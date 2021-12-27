package cobra

import (
	"github.com/robfig/cron/v3"
	"github.com/wxsatellite/goweb/framework"
	"log"
)

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) Container() framework.Container {
	return c.Root().container
}

/*** 用于定时脚本 ***/

// CronSpec 保存 cron 命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParentNull() {
	c.parent = nil
}

// AddCronCommand 用来创建一个 cron 任务
func (c *Command) AddCronCommand(spec string, cmd *Command) {

	// 获取根命令
	root := c.Root()

	// Dom 表示 day of month，每个月的第几天
	// Dow 表示 day of week，每个星球的第几天
	if root.Cron == nil {
		// cron.SecondOptional 设置支持秒：* * * * * * （秒分时日月周）
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}

	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Spec: spec,
		Cmd:  cmd,
	})

	// 新创建一个 cmd
	cronCmd := *cmd
	ctx := root.Context()
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	cronCmd.SetContainer(root.Container())
	_, _ = root.Cron.AddFunc(spec, func() {
		defer func() {
			// cron 中这个匿名函数是在 goroutine 中执行的，需要放置panic
			if err := recover(); err != nil {
				log.Println(err)
			}

			err := cronCmd.ExecuteContext(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
	})
}

/*** 用于定时脚本 ***/
