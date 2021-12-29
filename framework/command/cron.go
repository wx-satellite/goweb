package command

import (
	"errors"
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
	"github.com/wxsatellite/goweb/framework/cobra"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"github.com/wxsatellite/goweb/framework/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

/****  定时命令行 ****/

var cronDaemon = false

// 初始化定时命令行
func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start cron daemon")
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronRestartCommand)
	cronCommand.AddCommand(cronStopCommand)
	cronCommand.AddCommand(cronStateCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// 列出所有的定时任务
var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronSpecs := cmd.Root().CronSpecs
		var ps [][]string
		for _, spec := range cronSpecs {
			line := []string{spec.Type, spec.Spec, spec.Cmd.Use, spec.Cmd.Short, spec.ServiceName}
			ps = append(ps, line)
		}
		return nil
	},
}

// 启动 cron 进程
var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.Container()

		// 获取 app 服务
		appService := container.MustMake(app.Key).(app.App)

		// 获取 cron 的日志地址和pid地址
		pidFolder := appService.RuntimeFolder()
		serverPidFolder := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		serverLogFolder := filepath.Join(logFolder, "cron.log")
		currentFolder := appService.BaseFolder()

		// 守护进程的方式启动定时脚本
		// TODO：当开启多个子进程cron时会存在问题，pid文件、日志文件会相互覆盖（路径都一样），简单的处理方式就是每一个子进程的文件都不一样，引入app_id，并且将文件路径设置到子进程的环境变量中
		if cronDaemon {
			ctx := &daemon.Context{
				// 设置pid文件及权限
				PidFileName: serverPidFolder,
				PidFilePerm: 0664,

				// 设置日志文件及权限
				LogFileName: serverLogFolder,
				LogFilePerm: 0664,

				// 设置工作路径
				WorkDir: currentFolder,

				Umask: 027,

				// 子进程的参数，按照这个参数设置，子进程的命令为 ./goweb cron start --daemon=true
				Args: []string{"", "cron", "start", "--daemon=true"},

				Env: []string{}, // 设置环境变量
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := ctx.Reborn()
			if err != nil {
				return err
			}

			// d 不为空的时候，表示当前进程是父进程，可以从d中获取到子进程的信息
			// d 为空的时候，表示当前进程是子进程
			if d != nil {
				// 打印子进程的信息
				fmt.Println("cron serve started, pid:", d.Pid)
				fmt.Println("log file:", serverLogFolder)
				return nil
			}

			/* d == nil 即为子进程，那么启动定时脚本 */

			// 退出时，释放资源
			defer d.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("goweb cron")
			// 会阻塞
			cmd.Root().Cron.Start()
			return nil
		}

		// 非守护进程的方式，直接挂起
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFolder, []byte(content), 0664)
		if err != nil {
			return err
		}
		// 设置之后会影响ps的最后一列：
		//	501  6555   828   0  5:18PM ttys002    0:00.02 goweb cron
		gspt.SetProcTitle("goweb cron")
		cmd.Root().Cron.Run()
		return nil
	},
}

// 重新启动
var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.Container()

		appService := container.MustMake(app.Key).(app.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) > 0 {
			pid, _ := strconv.Atoi(string(content))
			if pid <= 0 {
				return errors.New("pid 不存在")
			}
			if !utils.CheckProcessExist(pid) {
				return errors.New("pid 不存在")
			}
			if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}

			// 检测是否真的退出
			for i := 0; i < 10; i++ {
				if utils.CheckProcessExist(pid) == false {
					break
				}
				time.Sleep(1 * time.Second)
			}
			fmt.Println("kill process:" + strconv.Itoa(pid))
		}

		cronDaemon = true
		return cronStartCommand.RunE(cmd, args)
	},
}

// 停止进程
var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {

		container := cmd.Container()

		appService := container.MustMake(app.Key).(app.App)

		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) > 0 {
			pid, _ := strconv.Atoi(string(content))
			if pid <= 0 {
				return errors.New("pid 不存在")
			}
			if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err = ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}

// 进程状态
var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "cron 常驻进程状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.Container()
		appService := container.MustMake(app.Key).(app.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if len(content) <= 0 {
			fmt.Println("no cron server start")
			return nil
		}
		pid, _ := strconv.Atoi(string(content))
		if pid <= 0 {
			return errors.New("pid 不存在")
		}
		if utils.CheckProcessExist(pid) {
			fmt.Println("cron server started, pid:", pid)
		}
		return nil
	},
}
