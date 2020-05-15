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
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		idx     int
		kvPair  *mvccpb.KeyValue
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

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		if len(delResp.PrevKvs) != 0 {
			for idx, kvPair = range delResp.PrevKvs {
				fmt.Printf(" index =%d ,del key =%s  del value =%s: ", idx, string(kvPair.Key), string(kvPair.Value))
			}
		}
	}

}
