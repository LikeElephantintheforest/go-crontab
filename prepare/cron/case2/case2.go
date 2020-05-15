package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

// 定义一个协程，定时检查所有的cron任务，谁过期了执行谁。
func main() {

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // key -> jobName
	)

	scheduleTable = make(map[string]*CronJob)

	now = time.Now()

	// 声明2个定时job
	expr = cronexpr.MustParse("*/5 * * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job2"] = cronJob

	go func() {

		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)

		//定时检查调度表
		for {

			now = time.Now()

			for jobName, cronJob = range scheduleTable {
				//check is past
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//在启动子协程做任务
					go func(jobName string) {
						fmt.Println("执行", jobName)
					}(jobName)

					//计算下次调度时间，并进行赋值
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间为", cronJob.nextTime)

				} else {
					//fmt.Println("无定时任务需要执行")
				}
			}

			select {
			case <-time.NewTimer(100 * time.Microsecond).C: //阻塞去读，timer会100毫米投递一次

			}

		}

	}()

	time.Sleep(50 * time.Second)

}
