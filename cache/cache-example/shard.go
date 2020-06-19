package main

import (
	"encoding/binary"
	"errors"
	"sync"
)

const (
	headerEntrySize = 4
	defaultValue    = 1024
)

type cacheShard struct {
	items        map[uint64]uint32 //key是hash后的key，value是字节数组中的index
	lock         sync.RWMutex
	array        []byte //全局字节数组，用来保存value
	tail         int    //存储字节数组的尾部下标
	headerBuffer []byte
}

func initNewShard() *cacheShard {
	return &cacheShard{
		items:        make(map[uint64]uint32, defaultValue),
		array:        make([]byte, defaultValue),
		tail:         1,
		headerBuffer: make([]byte, headerEntrySize),
	}
}

func (s *cacheShard) set(hashedKey uint64, entry []byte) {
	w := wrapEntry(entry)
	s.lock.Lock()
	index := s.push(w)
	s.items[hashedKey] = uint32(index)
	s.lock.Unlock()
}

func (s *cacheShard) push(data []byte) int {
	dataLen := len(data)
	index := s.tail
	s.save(data, dataLen)
	return index
}

func (s *cacheShard) save(data []byte, len int) {
	//使用小端序存储
	binary.LittleEndian.PutUint32(s.headerBuffer, uint32(len))
	s.copy(s.headerBuffer, headerEntrySize)
	s.copy(data, len)
}

func (s *cacheShard) copy(data []byte, len int) {
	s.tail += copy(s.array[s.tail:], data[:len])
}

func wrapEntry(entry []byte) []byte {
	blobLength := len(entry)
	blob := make([]byte, blobLength)
	copy(blob, entry)
	return blob
}

func (s *cacheShard) get(key string, hashedKey uint64) ([]byte, error) {
	s.lock.RLock()
	itemIndex := int(s.items[hashedKey])

	if itemIndex == 0 {
		s.lock.RUnlock()
		return nil, errors.New("key not fount")
	}

	blockSize := int(binary.LittleEndian.Uint32(s.array[itemIndex : itemIndex+headerEntrySize]))
	entry := s.array[itemIndex+headerEntrySize : itemIndex+headerEntrySize+blockSize]
	s.lock.RUnlock()
	return readEntry(entry), nil
}

func readEntry(data []byte) []byte {
	dst := make([]byte, len(data))
	copy(dst, data)
	return dst
}
