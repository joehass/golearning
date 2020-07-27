package _map

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type SyncMap struct {
	mu     sync.Mutex             // 加锁作用，保护dirty字段
	read   atomic.Value           //只读的数据，read中的数据是dirty中同步过来的
	dirty  map[interface{}]*entry //最新写入的数据，dirty中存有全部的数据
	misses int                    //计数器，每次需要读dirty 则 +1
}

type readOnly struct {
	m       map[interface{}]*entry //内建map
	amended bool                   //表示dirty里存在read里没有的key，通过该字段决定是否加锁读dirty,true:存在，false:不存在
}

var expunged = unsafe.Pointer(new(interface{}))

type entry struct {
	//p == nil：键值已经被删除，且m.dirty == nil
	//p == expunged：键值已经被删除，但m.dirty!=nil 且 m.dirty不存在该键值
	//除以上情况，则键值对存在，存在于m.read.m中，如果m.dirty!=nil 则也存在于m.dirty中
	p unsafe.Pointer // 等同于*interface{}
}

func newEntry(i interface{}) *entry {
	return &entry{p: unsafe.Pointer(&i)}
}

func (e *entry) Load() (value interface{}, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged {
		return nil, false
	}
	return *(*interface{})(p), true
}

func (e *entry) tryStore(i *interface{}) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

//确保这条记录没有被标记为删除
//如果该条目先前已删除，则必须在解锁m.mu之前将其添加到脏映射中。
func (e *entry) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}

func (e *entry) storeLocked(i *interface{}) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *SyncMap) Load(key interface{}) (value interface{}, ok bool) {
	//首先尝试从read中读取readonly对象
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	//如果不存在则从dirty中获取
	if !ok && read.amended {
		m.mu.Lock()
		//用于上面read获取没有加锁，为了安全再检查一次，避免遗漏
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		//确实不存在，从dirty中获取
		if !ok && read.amended {
			e, ok = m.dirty[key]
			//调用miss的逻辑
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.Load()
}

func (m *SyncMap) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	//当miss积累过多，会将dirty存入read，然后将amended = false，且m.dirty = nil
	m.read.Store(readOnly{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *SyncMap) dirtyLocked() {
	if m.dirty != nil {
		return
	}
	read, _ := m.read.Load().(readOnly)
	m.dirty = make(map[interface{}]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		//如果p == nil（即键值对被delete），则会在这个时机被置为expunged
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}

	return p == expunged
}

func (m *SyncMap) Store(key, value interface{}) {
	read, _ := m.read.Load().(readOnly)

	//如果read里存在，则尝试存到entry里，因为dirty中也是存的entry，所以改变entry的同时
	//也改变了dirty
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	//如果上一步没执行成功，则要分情况处理
	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)

	//和load一样，重新取一次
	if e, ok := read.m[key]; ok {
		//情况1：read里存在
		if e.unexpungeLocked() {
			//如果p==expunged，则需要先将entry赋值给dirty（因为expunged数据不会留在dirty中）
			m.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
		//情况2：read中不存在，但dirty里存在，则用值更新entry
		e.storeLocked(&value)
	} else {
		//情况3：read和dirty都不存在
		if !read.amended {
			//如果amended == false，则调用dirtyLocked 将read拷贝到dirty（除了被标记删除的数据）
			m.dirtyLocked()
			//然后将amended改为true
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		//将新的键值存入dirty
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}

func (m *SyncMap) Delete(key interface{}) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(m.dirty, key)
		}
		m.mu.Unlock()
	}
	if ok {
		e.delete()
	}
}

func (e *entry) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *SyncMap) Range(f func(key, value interface{}) bool) {
	read, _ := m.read.Load().(readOnly)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.Load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
