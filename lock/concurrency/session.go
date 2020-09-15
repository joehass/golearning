package concurrency

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

type Session struct {
	client *clientv3.Client
	opts   *sessionOptions
	id     clientv3.LeaseID

	cancel context.CancelFunc
	donec  <-chan struct{}
}

const defaultSessionTTL = 60

type sessionOptions struct {
	ttl     int              //过期时间
	leaseID clientv3.LeaseID // 租约id
	ctx     context.Context
}

type SessionOption func(*sessionOptions)

func NewSession(client *clientv3.Client, opts ...SessionOption) (*Session, error) {
	ops := &sessionOptions{ttl: defaultSessionTTL, ctx: client.Ctx()}
	for _, opt := range opts {
		opt(ops)
	}

	id := ops.leaseID
	if id == clientv3.NoLease { //没有租约
		resp, err := client.Grant(ops.ctx, int64(ops.ttl)) //创建租约
		if err != nil {
			return nil, err
		}
		id = resp.ID
	}

	ctx, cancel := context.WithCancel(ops.ctx)
	keepAlive, err := client.KeepAlive(ctx, id) //租约续租
	if err != nil || keepAlive == nil {
		cancel()
		return nil, err
	}

	donec := make(chan struct{})
	s := &Session{client: client, opts: ops, id: id, cancel: cancel, donec: donec}

	//如果租约失效则关闭channel
	go func() {
		defer close(donec)
		for range keepAlive {

		}
	}()

	return s, nil
}

func (s *Session) Client() *clientv3.Client {
	return s.client
}

func (s *Session) Lease() clientv3.LeaseID {
	return s.id
}
