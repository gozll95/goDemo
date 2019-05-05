package main

import (
	"fmt"
	"os"
	"time"

	"context"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	checkErr(err)

	kv := clientv3.NewKV(cli)

	// !+-PUT---------------------------------------
	putResp, err := kv.Put(context.TODO(), "/test/a", "something")
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:2 raft_term:2  <nil>}
	checkErr(err)
	fmt.Println(putResp)

	// 再写一个孩子
	kv.Put(context.TODO(), "/test/b", "another")

	// 再写一个同前缀的干扰项
	kv.Put(context.TODO(), "/testxxx", "干扰")
	// !--------------------------------------------

	// !+-GET---------------------------------------
	getResp, err := kv.Get(context.TODO(), "/test/a")
	checkErr(err)

	fmt.Println(getResp)
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:8 raft_term:2  [key:"/test/a" create_revision:2 mod_revision:6 version:3 value:"something" ] false 1}

	rangeResp, err := kv.Get(context.TODO(), "/test/", clientv3.WithPrefix())
	checkErr(err)
	fmt.Println(rangeResp)
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 revision:11 raft_term:2  [key:"/test/a" create_revision:2 mod_revision:9 version:4 value:"something"  key:"/test/b" create_revision:4 mod_revision:10 version:3 value:"another" ] false 2}
	// !--------------------------------------------

	// lease
	lease := clientv3.NewLease(cli)
	grantResp, err := lease.Grant(context.TODO(), 10)
	checkErr(err)

	kv.Put(context.TODO(), "/test/expireme", "gone...", clientv3.WithLease(grantResp.ID))

	/*
		keepResp, err := lease.KeepAliveOnce(context.TODO(), grantResp.ID)
		checkErr(err)
	*/
	// sleep一会..

	// op
	/*
		op1 := clientv3.OpPut("/hi", "hello", clientv3.WithPrevKV())
		opResp, err := kv.Do(context.TODO(), op1)
	*/

	// txn 事务
	txn := kv.Txn(context.TODO())

	txnResp, err := txn.If(clientv3.Compare(clientv3.Value("/hi"), "=", "hello")).
		Then(clientv3.OpGet("/hi")).
		Else(clientv3.OpGet("/test/", clientv3.WithPrefix())).
		Commit()
	checkErr(err)

	if txnResp.Succeeded { // If = true
		fmt.Println("~~~", txnResp.Responses[0].GetResponseRange().Kvs)
	} else { // If =false
		fmt.Println("!!!", txnResp.Responses[0].GetResponseRange().Kvs)
	}

}

func checkErr(err error) {
	if err != nil {
		os.Exit(-1)
	}
}

// 参考:https://yuerblog.cc/2017/12/12/etcd-v3-sdk-usage/
