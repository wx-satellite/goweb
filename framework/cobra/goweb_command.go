package cobra

import (
	"github.com/robfig/cron/v3"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"github.com/wxsatellite/goweb/framework/provider/distributed"
	"log"
	"time"
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

/*** 分布式定时器 ***/

// AddDistributedCronCommand 实现一个分布式定时器
// serviceName 这个服务的唯一名字，不允许带有空格
// spec 具体的执行时间
// cmd 具体的执行命令
// holdTime 表示如果我选择上了，这次选择持续的时间，也就是锁释放的时间
func (c *Command) AddDistributedCronCommand(serviceName string, spec string, cmd *Command, holdTime time.Duration) {
	root := c.Root()

	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}

	// cron命令的注释，这里注意Type为distributed-cron，ServiceName需要填写
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type:        "distributed-cron", // 注意这里是 distributed-cron
		Cmd:         cmd,
		Spec:        spec,
		ServiceName: serviceName, // 用于生成锁文件名
	})

	appService := root.Container().MustMake(app.Key).(app.App)
	distributedService := root.Container().MustMake(distributed.Key).(distributed.Distributed)

	// appId 应用的唯一ID
	appId := appService.AppId()

	// 复制传入的cmd
	ctx := cmd.Context()
	cronCmd := *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	_, _ = root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				return
			}
		}()

		// 节点选择器
		selectedAppId, err := distributedService.Select(serviceName, appId, holdTime)
		if err != nil {
			return
		}

		// 如果自己没有被选择到则退出
		if selectedAppId != appId {
			return
		}

		// 如果自己选择到了，则执行任务
		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
		return
	})

}

/*** 分布式定时器 ***/
