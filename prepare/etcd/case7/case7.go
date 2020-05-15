package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {

	var (
		config             clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
		watcher            clientv3.Watcher
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 1 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd中的kv对
	kv = clientv3.NewKV(client)

	//创建一个协程，用于变化指定key
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")

			kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值 kvs:", getResp.Kvs)              // [key:"/cron/jobs/job7" create_revision:24 mod_revision:24 version:1 value:"i am job7" ]
		fmt.Println("当前值:", string(getResp.Kvs[0].Value)) //i am job7
	}

	// watch当前 revision之后的变化记录
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)
	// 监听5秒后进行删除移除
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	//处理变化
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
