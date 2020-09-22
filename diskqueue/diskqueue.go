package diskqueue

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG = LogLevel(1)
	INFO  = LogLevel(2)
	WARN  = LogLevel(3)
	ERROR = LogLevel(4)
	FATAL = LogLevel(5)
)

type AppLogFunc func(lvl LogLevel, f string, args ...interface{})

func (l LogLevel) String() string {
	switch l {
	case 1:
		return "DEBUG"
	case 2:
		return "INFO"
	case 3:
		return "WARNING"
	case 4:
		return "ERROR"
	case 5:
		return "FATAL"
	}
	panic("invalid LogLevel")
}

type Interface interface {
	Close() error
	Depth() int64
	Put([]byte) error
	ReadChan() <-chan []byte
}

func (d *diskQueue) Delete() error {
	panic("implement me")
}

func (d *diskQueue) Empty() error {
	panic("implement me")
}

type diskQueue struct {
	name              string
	dataPath          string        //数据文件地址
	depth             int64         //消息数
	syncTimeout       time.Duration //同步时间
	depthChan         chan int64
	exitFlag          int32       //退出标记
	needSync          bool        //定时同步消息，即将消息内容持久化到磁盘
	writeChan         chan []byte //消息存储队列
	readChan          chan []byte //消息读取队列
	writeResponseChan chan error
	writeFile         *os.File //写文件句柄
	writeFileNum      int64    //写文件数量
	readFile          *os.File
	readFileNum       int64 //读文件数量
	writeBuf          bytes.Buffer
	reader            *bufio.Reader

	writePos        int64 //写位置
	readPos         int64 //记录当前readFileNum指向文件已经读取并发送出去的文件位置
	nextReadPos     int64 //记录当前readFileNum指向的文件已经读取但是还没发送出去的文件位置
	nextReadFileNum int64

	exitChan     chan int //退出信号
	exitSyncChan chan int

	sync.RWMutex
	logf AppLogFunc
}

func New(name string, dataPath string, maxBytesPerFile int64,
	minMsgSize int32, maxMsgSize int32,
	syncEvery int64, syncTimeout time.Duration, logf AppLogFunc) Interface {

	d := diskQueue{
		name:              name,
		dataPath:          dataPath,
		syncTimeout:       syncTimeout,
		readChan:          make(chan []byte),
		depthChan:         make(chan int64),
		writeChan:         make(chan []byte),
		writeResponseChan: make(chan error),
		exitChan:          make(chan int),
		exitSyncChan:      make(chan int),
		logf:              logf,
	}
	err := d.retrieveMetaData()
	//文件不存在
	if err != nil && !os.IsNotExist(err) {
		d.logf(ERROR, "DISKQUEUE(%s) failed to retrieveMetaData - %s", d.name, err)
	}

	go d.ioLoop()
	return &d
}

func (d *diskQueue) ioLoop() {
	var dataRead []byte
	var err error
	var count int64
	var r chan []byte

	syncTicker := time.NewTicker(d.syncTimeout)

	for {

		if d.needSync {
			err = d.sync()
			if err != nil {
				d.logf(ERROR, "DISKQUEUE(%s) failed to sync - %s", d.name, err)
			}
			count = 0
		}

		if d.readPos < d.writePos {
			if d.nextReadPos == d.readPos {
				dataRead, err = d.readOne()
				if err != nil {
					d.logf(ERROR, "DISKQUEUE(%s) reading at %d of %s - %s",
						d.name, d.readPos, d.fileName(d.readFileNum), err)
				}
			}
			r = d.readChan
		} else {
			r = nil
		}

		select {
		case r <- dataRead:
			count++

		case d.depthChan <- d.depth:
		case dataWrite := <-d.writeChan:
			count++
			d.writeResponseChan <- d.writeOne(dataWrite)
		case <-syncTicker.C:
			if count == 0 {
				continue
			}
			d.needSync = true
		}
	}
}

func (d *diskQueue) moveForward() {
	oldReadFileNum := d.readFileNum
	d.readFileNum = d.nextReadFileNum
	d.readPos = d.nextReadFileNum
	d.depth -= 1

	if oldReadFileNum != d.nextReadFileNum {

	}

}

//把读队列长度置空
func (d *diskQueue) checkTailCorruption(depth int64) {

}

func (d *diskQueue) readOne() ([]byte, error) {
	var err error
	var msgSize int32

	if d.readFile == nil {
		curFileName := d.fileName(d.readFileNum)
		d.readFile, err = os.OpenFile(curFileName, os.O_RDONLY, 0600)
		if err != nil {
			return nil, err
		}

		d.logf(INFO, "DISKQUEUE(%s): readOne() opened %s", d.name, curFileName)

		if d.readPos > 0 {

		}
		//读取文件
		d.reader = bufio.NewReader(d.readFile)
	}

	//从reader中读取数据到msgSize中
	err = binary.Read(d.reader, binary.BigEndian, &msgSize)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}
	readBuf := make([]byte, msgSize)
	//读取指定长度的数据到readBuf中
	_, err = io.ReadFull(d.reader, readBuf)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}

	totalBytes := int64(4 + msgSize)

	d.nextReadPos = d.readPos + totalBytes
	d.nextReadFileNum = d.readFileNum

	return readBuf, nil
}

//持久化数据
func (d *diskQueue) sync() error {
	if d.writeFile != nil {
		//当前内容持久化，马上写到磁盘
		err := d.writeFile.Sync()
		if err != nil {
			d.writeFile.Close()
			d.writeFile = nil
			return err
		}
	}
	err := d.persistMetaData()
	if err != nil {
		return err
	}
	d.needSync = false
	return nil
}

func (d *diskQueue) ReadChan() <-chan []byte {
	return d.readChan
}

//持久化数据到磁盘
func (d *diskQueue) persistMetaData() error {
	var f *os.File
	var err error

	fileName := d.metaDataFileName()
	tmpFileName := fmt.Sprintf("%s.%d.tmp", fileName, rand.Int())

	//以读写模式打开文件，文件不存在则创建
	f, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	//文件句柄写入内容
	_, err = fmt.Fprintf(f, "%d\n%d,%d\n%d,%d\n",
		d.depth,
		d.readFileNum, d.readPos,
		d.writeFileNum, d.writePos)
	if err != nil {
		f.Close()
		return err
	}
	f.Sync()
	f.Close()

	return os.Rename(tmpFileName, fileName)
}

func (d *diskQueue) writeOne(data []byte) error {
	var err error

	if d.writeFile == nil {
		curFileName := d.fileName(d.writeFileNum)
		//以读写模式打开，不存在则创建
		d.writeFile, err = os.OpenFile(curFileName, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		d.logf(INFO, "DISKQUEUE(%s): writeOne() opened %s", d.name, curFileName)
	}
	dataLen := int32(len(data))

	//充值buffer
	d.writeBuf.Reset()
	err = binary.Write(&d.writeBuf, binary.BigEndian, dataLen)
	if err != nil {
		return err
	}
	//数据写入buffer
	_, err = d.writeBuf.Write(data)
	if err != nil {
		return err
	}

	//buffer写入文件
	_, err = d.writeFile.Write(d.writeBuf.Bytes())
	if err != nil {
		d.writeFile.Close()
		d.writeFile = nil
		return err
	}

	totalBytes := int64(4 + dataLen)
	d.writePos += totalBytes
	d.depth += 1 //消息数

	return err
}

func (d *diskQueue) fileName(fileNum int64) string {
	return fmt.Sprintf(path.Join(d.dataPath, "%s.diskqueue.%06d.dat"), d.name, fileNum)
}

func (d *diskQueue) Close() error {
	err := d.exit(false)
	if err != nil {
		return err
	}

	return d.sync()
}

func (d *diskQueue) exit(deleted bool) error {
	d.Lock()
	defer d.Unlock()

	d.exitFlag = 1

	if deleted {
		d.logf(INFO, "DISKQUEU(%s): deleting", d.name)
	} else {
		d.logf(INFO, "DISKQUEUE(%S): closing", d.name)
	}
	close(d.exitChan)

	<-d.exitSyncChan

	close(d.depthChan)

	if d.readFile != nil {
		d.readFile.Close()
		d.readFile = nil
	}

	if d.writeFile != nil {
		d.writeFile.Close()
		d.writeFile = nil
	}

	return nil
}

func (d *diskQueue) Put(data []byte) error {
	d.RLock()
	defer d.RUnlock()

	if d.exitFlag == 1 {
		return errors.New("exiting")
	}
	d.writeChan <- data
	return <-d.writeResponseChan
}

func (d *diskQueue) Depth() int64 {
	depth, ok := <-d.depthChan
	if !ok {
		depth = d.depth
	}
	return depth
}

func (d *diskQueue) retrieveMetaData() error {
	var f *os.File
	var err error

	//数据文件地址
	fileName := d.metaDataFileName()
	//以只读模式和读写权限打开文件
	f, err = os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func (d *diskQueue) metaDataFileName() string {
	return fmt.Sprintf(path.Join(d.dataPath, "%s.diskQueue.meta.dat"), d.name)
}
