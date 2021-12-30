package console

import (
	"github.com/wxsatellite/goweb/app/console/command/demo"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/cobra"
	"github.com/wxsatellite/goweb/framework/command"
	"time"
)

// RunCommand 初始化根command并运行
func RunCommand(container framework.Container) (err error) {

	// 根 command
	rootCommand := &cobra.Command{
		Use:   "goweb",
		Short: "goweb 命令",
		Long:  "goweb 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return cmd.Help()
		},
		//Args: cobra.ExactArgs(1),
		// 不需要出现 cobra 默认的 completion 子命令
		// 不设置的话，当存在子命令的时候，会多出一个cobra自动创建的completion子命令：
		//  Available Commands:
		//  	completion  generate the autocompletion script for the specified shell
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 为根 command 设置容器
	rootCommand.SetContainer(container)
	// 绑定框架的命令
	command.AddKernelCommands(rootCommand)

	// 绑定业务的命令
	AddAppCommand(rootCommand)

	// 执行
	return rootCommand.Execute()
}

func AddAppCommand(rootCommand *cobra.Command) {
	rootCommand.AddCronCommand("* * * * * *", demo.PrintCommand)

	// 启动一个分布式任务调度，调度的服务名称为init_func_for_test，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s才释放
	rootCommand.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.PrintCommand, 2*time.Second)
	return
}
