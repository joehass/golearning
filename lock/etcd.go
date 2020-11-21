package lock

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"golearning/lock/concurrency"
	"log"
	"sync"
	"time"
)

type Lock interface {
	TryLock(key string) (clientv3.LeaseID, error)
	Unlock(leaseId clientv3.LeaseID) error
	TryLockWithFunc(key string, i int, f func() error) error
}

type Etcd struct {
	cli  *clientv3.Client
	sess *concurrency.Session
	m    map[clientv3.LeaseID]*concurrency.Mutex
	mu   sync.Mutex
}

func NewEtcd(etcdAddrs []string) Lock {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: time.Duration(3) * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return &Etcd{
		cli: cli,
		m:   make(map[clientv3.LeaseID]*concurrency.Mutex),
	}
}

func (e *Etcd) TryLockWithFunc(key string, i int, f func() error) error {
	s1, err := concurrency.NewSession(e.cli)
	if err != nil {
		log.Fatal(err)
	}

	m1 := concurrency.NewMutex(s1, "/my-lock/"+key)

	err = m1.Lock(context.TODO(), i)
	defer func() {
		err = m1.Unlock(context.TODO())
		if err != nil {
			log.Println("unlock error:", err)
		}
	}()
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
		return err
	}

	return nil
}

func (e *Etcd) TryLock(key string) (clientv3.LeaseID, error) {
	s1, err := concurrency.NewSession(e.cli)
	if err != nil {
		log.Fatal(err)
	}
	e.sess = s1
	m1 := concurrency.NewMutex(s1, "/my-lock/"+key)
	// 会话s1获取锁
	err = m1.Lock(context.TODO(), 1)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	e.mu.Lock()
	e.m[e.sess.Lease()] = m1
	e.mu.Unlock()
	return e.sess.Lease(), nil
}

func (e *Etcd) Unlock(leaseId clientv3.LeaseID) error {
	m := e.m[leaseId]
	err := m.Unlock(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
