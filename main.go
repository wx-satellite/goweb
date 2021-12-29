package main

import (
	"github.com/wxsatellite/goweb/app/console"
	"github.com/wxsatellite/goweb/app/http"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"github.com/wxsatellite/goweb/framework/provider/distributed"
	"github.com/wxsatellite/goweb/framework/provider/kernel"
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
//func main() {
//	core := framework.NewCore()
//	registerRouter(core)
//	server := &http.Server{
//		Addr:    ":8080",
//		Handler: core,
//	}
//	// 该 goroutine 负责启动服务
//	go func() {
//		_ = server.ListenAndServe()
//	}()
//
//	// main 所在的当前goroutine负责监听信号
//	quit := make(chan os.Signal)
//	// 监听指定信号
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//	// 阻塞当前 goroutine，只能等到捕获到了信号才继续后续的流程
//	<-quit
//
//	// server.Shutdown 方法是个阻塞方法，一旦执行之后，它会阻塞当前 Goroutine，并且在所有连接请求都结束之后，才继续往后执行
//	if err := server.Shutdown(context.Background()); err != nil {
//		log.Fatal("Server Shutdown：", err)
//	}
//
//}

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

//func m() {
//	core := gin.New()
//	// 注册服务
//	core.Use(gin.Recovery())
//
//	registerRouter(core)
//
//	server := &http.Server{
//		Addr:    ":8082",
//		Handler: core,
//	}
//
//	// 这个 goroutine 用于提供服务
//	go func() {
//		_ = server.ListenAndServe()
//
//	}()
//
//	// 当前 goroutine 等待信号
//	quit := make(chan os.Signal)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//	// 阻塞当前 goroutine 等待信号
//	<-quit
//
//	// 最长等待5秒用来做安全退出
//	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	if err := server.Shutdown(timeoutCtx); err != nil {
//		log.Fatal("Server Shutdown:", err)
//	}
//}

func main() {
	container := framework.NewGoWebContainer()

	// 绑定应用目录服务
	_ = container.Bind(&app.Provider{})
	_ = container.Bind(&distributed.LocalProvider{})

	// 这个 Web 引擎不仅仅是调用了 Gin 创建 Web 引擎的方法，更重要的是需要注册业务的路由
	// 所以 http.NewHttpEngine 这个创建 Web 引擎的方法必须放在业务层，不能放在框架中
	if engine, err := http.NewHttpEngine(); err == nil {
		_ = container.Bind(&kernel.Provider{Engine: engine})
	}
	// 运行根command
	_ = console.RunCommand(container)
}

//func main() {
//	gspt.SetProcTitle("hade cron")
//	time.Sleep(100 * time.Minute)
//}
