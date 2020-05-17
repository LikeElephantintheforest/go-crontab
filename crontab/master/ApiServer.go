package master

import (
	"net"
	"net/http"
	"time"
)

var (
	G_apiServer *ApiServer
)

// 任务Http接口
type ApiServer struct {
	httpServer *http.Server
}

func handlerJobSave(w http.ResponseWriter, r *http.Request) {

}

//  初始化服务
func InitApiServer() (err error) {

	var (
		mux        *http.ServeMux //定制路由
		httpServer *http.Server
		listener   net.Listener
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handlerJobSave)

	if listener, err = net.Listen("tcp", ":8080"); err != nil {
		return
	}

	httpServer = &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Handler:           mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	return nil

}
