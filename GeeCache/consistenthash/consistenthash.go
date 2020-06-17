package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//一致性ashh算法
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //自定义的hash函数，默认为crc32.ChecksumIEEE算法
	replicas int            //虚拟节点倍数：
	keys     []int          //hash环
	hashMap  map[int]string //虚拟节点和真实节点的映射表，键是虚拟节点的hash值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

//添加真实节点，keys：真实节点
//对每一个真实节点key：对应创建m.replicas个虚拟节点，虚拟节点的名称是:strconv.Itoa(i)+key，即通过添加编号的方式区分不同虚拟节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//将虚拟节点添加到环上
			m.keys = append(m.keys, hash)
			//增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	//以递增顺序排序
	sort.Ints(m.keys)
}

//选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算key的hash值
	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配的虚拟节点的下标idx，从m.keys中获取到对应的hash值，如何idx == len(m.keys)
	//说明应选择m.keys[0]
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	//通过 hashMap 映射得到真实的节点。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
