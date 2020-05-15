package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {

	var (
		config clientv3.Config
		client *clientv3.Client
		err    error

		kv clientv3.KV

		lease          clientv3.Lease
		leaseId        clientv3.LeaseID
		leaseGrantResp *clientv3.LeaseGrantResponse

		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse

		ctx        context.Context
		cancelFunc context.CancelFunc

		txn     clientv3.Txn
		txnResp *clientv3.TxnResponse
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 1: 申请租约

	//通过客户端申请租约
	lease = clientv3.NewLease(client)
	//默认时间是秒
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	leaseId = leaseGrantResp.ID

	//创建一个用于取消自动续约的context
	ctx, cancelFunc = context.WithCancel(context.TODO())
	//确保函数退出后，自动续租会停止
	defer cancelFunc()

	defer lease.Revoke(context.TODO(), leaseId)

	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
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

	txn = kv.Txn(context.TODO())

	//如key不存在，则说明为新建事务
	txnResp, err = txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")).Commit()

	if err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		fmt.Println("locked ", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}

	fmt.Println("run~")
	time.Sleep(5 * time.Second)

}
