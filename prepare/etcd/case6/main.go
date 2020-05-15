package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {

	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		kv             clientv3.KV
		lease          clientv3.Lease
		leaseId        clientv3.LeaseID
		leaseGrantResp *clientv3.LeaseGrantResponse
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 1 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//通过客户端申请租约
	lease = clientv3.NewLease(client)
	//默认时间是秒
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	leaseId = leaseGrantResp.ID

	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("因为长时间失联等异常情况，租约已过期，无法续约")
					goto END
				} else {
					fmt.Println("续约成功", keepResp.ID)
				}
			}
		}
	END:
	}()

	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "default", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功", putResp.Header.Revision)

	for {

		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv 被移除了")
		} else {
			fmt.Println(getResp.Kvs)
		}

		time.Sleep(2 * time.Second)

	}

}
