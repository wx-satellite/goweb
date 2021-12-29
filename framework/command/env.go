package command

import (
	"fmt"
	"github.com/wxsatellite/goweb/framework/cobra"
	"github.com/wxsatellite/goweb/framework/provider/env"
	"github.com/wxsatellite/goweb/framework/utils"
)

/**
环境变量服务按照先读取本地默认.env 文件，再读取运行环境变量的方式来实现，并且为其设置了最关键的环境变量 APP_ENV 来表示这个应用当前运行的环境。
后续再根据这个 APP_ENV 来获取具体环境的本地配置文件。
*/
func initEnvCommand() *cobra.Command {
	envCommand.AddCommand(envListCommand)
	return envCommand
}

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前app环境",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.Container()
		envService := container.MustMake(env.Key).(env.Env)
		fmt.Println("environment", envService.AppEnv())
	},
}

var envListCommand = &cobra.Command{
	Use:   "list",
	Short: "获取所有的环境变量",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.Container()
		envService := container.MustMake(env.Key).(env.Env)
		envs := envService.All()
		var outs [][]string
		for k, v := range envs {
			outs = append(outs, []string{k, v})
		}
		utils.PrettyPrint(outs)
	},
}
