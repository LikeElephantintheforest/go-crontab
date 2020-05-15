package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err            error
	combinedOutput []byte
}

//创建一个子协程，做两秒休眠后打印信息。但在第一秒休眠后将其中断。
func main() {

	var (
		command       *exec.Cmd
		ctx           context.Context
		cancelFunc    context.CancelFunc
		resultChannel chan *result
		res           *result
	)

	resultChannel = make(chan *result, 1000)

	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {

		var (
			combinedOutput []byte
			err            error
		)

		command = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2 ; echo 1 ; echo 2;")

		combinedOutput, err = command.CombinedOutput()

		resultChannel <- &result{
			err:            err,
			combinedOutput: combinedOutput,
		}

	}()

	time.Sleep(1 * time.Second)

	cancelFunc()

	res = <-resultChannel

	fmt.Println(res.err, string(res.combinedOutput))

}
