package command

import "github.com/wxsatellite/goweb/framework/cobra"

func AddKernelCommands(rootCommand *cobra.Command) {

	// app
	rootCommand.AddCommand(initAppCommand())

	// cron
	rootCommand.AddCommand(initCronCommand())

	// env
	rootCommand.AddCommand(initEnvCommand())
	return
}
