package diskqueue

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type tbLog interface {
	Log(...interface{})
}

func NewTestLogger(tbl tbLog) AppLogFunc {
	return func(lvl LogLevel, f string, args ...interface{}) {
		tbl.Log(fmt.Sprintf(lvl.String()+": "+f, args...))
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

func NotNil(t *testing.T, object interface{}) {
	if isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tExpected value not to be <nil>\033[39m\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   %#v (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func Nil(t *testing.T, object interface{}) {
	if !isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   <nil> (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, object)
		t.FailNow()
	}
}

func TestDiskQueue(t *testing.T) {
	l := NewTestLogger(t)

	dqName := "test_disk_queue" + strconv.Itoa(int(time.Now().Unix()))
	//创建文件夹
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	//删除文件夹
	defer os.RemoveAll(tmpDir)
	dq := New(dqName, tmpDir, 1024, 4, 1<<10, 2500, 2*time.Second, l)
	defer dq.Close()
	NotNil(t, dq)
	Equal(t, int64(0), dq.Depth())
	msg := []byte("test")
	err = dq.Put(msg)
	Nil(t, err)
	Equal(t, int64(1), dq.Depth())
	msgOut := <-dq.ReadChan()
	fmt.Println(string(msgOut))
	Equal(t, msg, msgOut)
}

func TestDiskQueueRoll(t *testing.T) {
	l := NewTestLogger(t)
	dqName := "test_disk_queue_roll" + strconv.Itoa(int(time.Now().Unix()))
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("nsq-test-%d", time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmpDir)
	//创建一个10大小的byte切片
	msg := bytes.Repeat([]byte{0}, 10)
	ml := int64(len(msg))
	dq := New(dqName, tmpDir, 10*(ml+4), int32(ml), 1<<10, 2500, 2*time.Second, l)
	defer dq.Close()
	NotNil(t, dq)
	Equal(t, int64(0), dq.Depth())

	for i := 0; i < 10; i++ {
		err := dq.Put(msg)
		Nil(t, err)
		Equal(t, int64(i+1), dq.Depth())
	}

	Equal(t, int64(1), dq.(*diskQueue).writeFileNum)
	Equal(t, int64(0), dq.(*diskQueue).writePos)

	for i := 10; i > 0; i-- {
		Equal(t, msg, <-dq.ReadChan())
		Equal(t, int64(i-1), dq.Depth())
	}

}
