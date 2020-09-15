package concurrency

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"strconv"
)

type Mutex struct {
	s *Session

	pfx   string
	myKey string
	myRev int64
	hdr   *etcdserverpb.ResponseHeader
}

func NewMutex(s *Session, pfx string) *Mutex {
	return &Mutex{s, pfx + "/", "", -1, nil}
}

func (m *Mutex) Lock(ctx context.Context, i int) error {
	s := m.s
	client := m.s.Client()

	m.myKey = fmt.Sprintf("%s%x", m.pfx, s.Lease())
	//判断版本号是否等于0，0表示锁不存在
	cmp := clientv3.Compare(clientv3.CreateRevision(m.myKey), "=", 0)

	//向加锁的key中存储一个空值，这个操作就是一个加锁的操作，
	//但是这把锁是有超时时间的，超时的时间是session的默认时长。超时是为了防止锁没有被正常释放导致死锁
	put := clientv3.OpPut(m.myKey, strconv.Itoa(i), clientv3.WithLease(s.Lease()))

	//get就是通过key来查询
	get := clientv3.OpGet(m.myKey)

	//注意这里是用m.pfx来查询的，并且带了查询参数WithFirstCreate()。
	//使用pfx来查询是因为其他的session也会用同样的pfx来尝试加锁，
	//并且因为每个LeaseID都不同，所以第一次肯定会put成功。
	//但是只有最早使用这个pfx的session才是持有锁的，所以这个getOwner的含义就是这样的
	getOwner := clientv3.OpGet(m.pfx, clientv3.WithFirstCreate()...)

	//事务抢锁
	resp, err := client.Txn(ctx).If(cmp).Then(put, getOwner).Else(get, getOwner).Commit()
	if err != nil {
		return err
	}

	//当前版本号
	m.myRev = resp.Header.Revision
	//resp.Succeeded是cmp为true时值为true，否则是false。
	//这里的判断表明当同一个session非第一次尝试加锁，当前的版本号应该取这个key的最新的版本号
	if !resp.Succeeded {
		m.myRev = resp.Responses[0].GetResponseRange().Kvs[0].CreateRevision
	}

	//下面是取得锁的持有者的key。如果当前没有人持有这把锁，那么默认当前会话获得了锁。
	//或者锁持有者的版本号和当前的版本号一致， 那么当前的会话就是锁的持有者。
	ownerKey := resp.Responses[1].GetResponseRange().Kvs
	fmt.Printf("i:%d rev:%d resp.Succeeded:%t ownerKey[0].CreateRevision: %d \n", i, m.myRev, resp.Succeeded, ownerKey[0].CreateRevision)
	if len(ownerKey) == 0 || ownerKey[0].CreateRevision == m.myRev {
		m.hdr = resp.Header
		return nil
	}

	hdr, werr := waitDeletes(ctx, client, m.pfx, m.myRev-1)
	if werr != nil {
		_ = m.Unlock(client.Ctx())
	} else {
		m.hdr = hdr
	}
	return werr
}

func (m *Mutex) Unlock(ctx context.Context) error {
	client := m.s.client
	if _, err := client.Delete(ctx, m.myKey); err != nil {
		return err
	}

	m.myKey = "\x00"
	m.myRev = -1
	return nil
}
