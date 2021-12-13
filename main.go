package main

import (
	"context"
	"github.com/wxsatellite/goweb/framework"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/**
SIGINT   ctrl+c  该信息可以捕获和处理
SIGQUIT  ctrl+\  该信号可以捕获和处理
SIGTERM  kill    该信号可以捕获和处理
SIGKILL  kill -9 不可捕获和处理，进程会被直接杀死


每隔多少时间，执行一次操作，应该会使用 time.Sleep 来做间隔时长，而在 Shutdown 里面演示了如何使用 time.Ticker 来进行轮询设置。
该方法消耗的 CPU 会远低于 time.Sleep。因为其内部使用了 runtimeTimer 数据结构，Go做了很多复杂的优化，而 time.Sleep 就是 gopark
和 goready，它需要调度先唤醒当前 goroutine，才能执行后续的逻辑

*/

// 主goroutine 和 子goroutine 都是可以接收信号的，不过主goroutine中接收信号做收尾工作比较简单
func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}
	// 该 goroutine 负责启动服务
	go func() {
		_ = server.ListenAndServe()
	}()

	// main 所在的当前goroutine负责监听信号
	quit := make(chan os.Signal)
	// 监听指定信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞当前 goroutine，只能等到捕获到了信号才继续后续的流程
	<-quit

	// server.Shutdown 方法是个阻塞方法，一旦执行之后，它会阻塞当前 Goroutine，并且在所有连接请求都结束之后，才继续往后执行
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown：", err)
	}

}

//func main() {
//	http.HandleFunc("/haha", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Println(request.RequestURI)
//		fmt.Println(request.URL.Path)
//		panic("123")
//	})
//	http.HandleFunc("/name", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Println(request.RequestURI)
//		fmt.Println(request.URL.Path)
//	})
//	http.ListenAndServe(":8080", nil)
//}
