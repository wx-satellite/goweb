package console

import (
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/cobra"
	"github.com/wxsatellite/goweb/framework/command"
)

// RunCommand 初始化根command并运行
func RunCommand(container framework.Container) (err error) {

	// 根 command
	rootCommand := &cobra.Command{
		Use:   "goweb",
		Short: "goweb 命令",
		Long:  "goweb 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现 cobra 默认的 completion 子命令
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
	return
}
