package lock

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"log"
	"sync"
	"testing"
	"time"
)

var cli *clientv3.Client

func init() {
	endpoints := []string{"10.10.114.123:2379"}

	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	}
	cli1, err := clientv3.New(cfg)
	if err != nil {
		log.Println("new cli error:", err)
		panic(err)
	}
	cli = cli1
}

func TestL(t *testing.T) {
	addrs := []string{"10.10.114.123:2379"}
	lock := NewEtcd(addrs)

	for i := 0; i < 10; i++ {
		//go func(i int) {
		//	lock.TryLockWithFunc("lockname",i, func() error {
		//		fmt.Println(i)
		//		time.Sleep(3* time.Second)
		//		return nil
		//	})
		//}(i)

		go func(i int) {
			id, err := lock.TryLock("lockname")
			if err != nil {
				return
			}
			fmt.Println(i)
			time.Sleep(3 * time.Second)
			lock.Unlock(id)
		}(i)

	}

	forever := make(chan bool)
	<-forever
}

var n = 0

// 使用worker模拟锁的抢占
func worker(key string) error {

	s, err := concurrency.NewSession(cli)
	if err != nil {
		return err
	}

	m := concurrency.NewMutex(s, "/"+key)

	err = m.Lock(context.TODO())
	if err != nil {
		log.Println("lock error:", err)
		return err
	}

	defer func() {
		err = m.Unlock(context.TODO())
		if err != nil {
			log.Println("unlock error:", err)
		}
	}()

	log.Println("get lock: ", n)
	n++
	time.Sleep(3 * time.Second) // 模拟执行代码

	return nil
}

func TestL3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		err := worker("lockname")
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := worker("lockname")
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := worker("lockname")
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
}
