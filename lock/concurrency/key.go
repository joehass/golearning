package concurrency

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func waitDelete(ctx context.Context, client *clientv3.Client, key string, rev int64) error {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wr clientv3.WatchResponse
	wch := client.Watch(cctx, key, clientv3.WithRev(rev))
	for wr = range wch {
		for _, ev := range wr.Events {
			if ev.Type == mvccpb.DELETE {
				fmt.Printf("delete key %s modRevision:%d createRevision:%d \n ", string(ev.Kv.Value), ev.Kv.ModRevision, ev.Kv.CreateRevision)
				return nil
			}
		}
	}

	if err := wr.Err(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return fmt.Errorf("lost watcher waiting for delete")
}

func waitDeletes(ctx context.Context, client *clientv3.Client, pfx string, maxCreateRev int64) (*etcdserverpb.ResponseHeader, error) {
	getOpts := append(clientv3.WithLastCreate(), clientv3.WithMaxCreateRev(maxCreateRev))
	for {
		resp, err := client.Get(ctx, pfx, getOpts...)
		if err != nil {
			return nil, err
		}

		if len(resp.Kvs) == 0 {
			return resp.Header, nil
		}
		lastKey := string(resp.Kvs[0].Key)
		if err = waitDelete(ctx, client, lastKey, resp.Header.Revision); err != nil {
			return nil, err
		}
	}
}
