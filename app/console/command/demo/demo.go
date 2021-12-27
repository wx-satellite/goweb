package demo

import (
	"github.com/wxsatellite/goweb/framework/cobra"
	"log"
)

var PrintCommand = &cobra.Command{
	Use:     "print",
	Short:   "测试输出",
	Long:    "测试输出",
	Aliases: []string{"f", "fo"},
	Example: "例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("hello world")
		return nil
	},
}
