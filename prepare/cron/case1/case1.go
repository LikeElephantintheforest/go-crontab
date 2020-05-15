package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {

	var (
		express  *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	if express, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
	}

	now = time.Now()

	nextTime = express.Next(now)

	fmt.Println(now, nextTime)

	//等待X执行
	time.AfterFunc(nextTime.Sub(now) /*下一个时间减去当前时间*/, func() {
		fmt.Println("被调度了 : ", nextTime)
	})

	time.Sleep(50 * time.Second)

}

//2020-05-10 23:12:37.093362 +0800 CST m=+0.000631941 2020-05-10 23:12:40 +0800 CST
//被调度了 :  2020-05-10 23:12:40 +0800 CST
