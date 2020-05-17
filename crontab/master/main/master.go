package main

import (
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
}

func main() {

	//初始化线程环境
	initEnv()

}
