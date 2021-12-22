package command

import (
	"context"
	"github.com/wxsatellite/goweb/framework/cobra"
	"github.com/wxsatellite/goweb/framework/provider/kernel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	return appCommand
}

// appCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "应用服务控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 打印帮助信息
		return cmd.Help()
	},
}

// appStartCommand 是 app 命令的子命令，用于启动应用服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动应用服务",
	Long:  "启动应用服务，它是一个web服务",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		/* 启动一个web服务 */

		// 获取根command中存放的容器
		container := cmd.Container()

		// 从容器中获取web服务引擎
		service := container.MustMake(kernel.Key).(kernel.Kernel)
		engine := service.Engine()

		// 创建一个 server 服务
		server := &http.Server{Addr: ":8888", Handler: engine}

		// 启动服务
		go func() {
			_ = server.ListenAndServe()
		}()

		// 创建信号等待，用于安全退出服务
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 只有监听到上面三个信号的时候才会继续走后面的逻辑
		<-quit

		// 调用 Shutdown graceful 退出
		timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		// 优雅退出
		if err = server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown", err)
		}

		return
	},
}
