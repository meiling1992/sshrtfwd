package xdaemon

import (
	"fmt"
	"os"
	"os/exec"
	"sshrts/api/logger"
	"sshrts/api/single"
	"strconv"
	"syscall"
	"time"
)

const ENV_NAME = "XDAEMON_INDEX"

var RunIndex int = 0

type Deamon struct {
	LogFile     string
	MaxCount    int // if 0 无限循环
	MaxError    int
	MinExitTime int64 // 单位秒(s)
	s           *single.Single
}

func NewDaemon(LogFileVar string, sl *single.Single) *Deamon {

	return &Deamon{
		LogFile:     LogFileVar,
		MaxCount:    0,
		MaxError:    3,
		MinExitTime: 10,
		s:           sl,
	}
}

func BackGround(logfile string, isExit bool) (*exec.Cmd, error) {

	RunIndex++
	logger.InfoField("XDaemon-BackGround", "ENV_NAME", os.Getenv(ENV_NAME))
	envIndex, err := strconv.Atoi(os.Getenv(ENV_NAME))
	if err != nil {
		envIndex = 0
	}
	if RunIndex <= envIndex { // 子进程 退出
		return nil, nil
	}
	logger.InfoField("XDaemon-BackGround", "envIndex", strconv.Itoa(envIndex))
	logger.InfoField("XDaemon-BackGround", "RunIndex", strconv.Itoa(RunIndex))
	// 设置子进程环境变量
	env := os.Environ()
	logger.InfoField("XDaemon-BackGround", "RunIndex", fmt.Sprintf("ENV_NAME(%s)=RunIndex(%s)", ENV_NAME, RunIndex))

	env = append(env, fmt.Sprintf("%s=%s", ENV_NAME, RunIndex))

	NewSysProAttr := func() *syscall.SysProcAttr {
		return &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	// 启动子进程
	startProc := func(args, env []string, logFile string) (*exec.Cmd, error) {
		cmd := &exec.Cmd{
			Path:        args[0],
			Args:        args,
			Env:         env,
			SysProcAttr: NewSysProAttr(),
		}
		if logFile != "" {
			stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				logger.ErrorField("BackGround.ERROR", "OpenFile", fmt.Sprintf("error:%v", err))
				return nil, err
			}
			cmd.Stderr = stdout
			cmd.Stdout = stdout
		}
		err := cmd.Start()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
	cmd, err := startProc(os.Args, env, logfile)
	if err != nil {
		logger.ErrorField("BackGround.ERROR", "child process ", fmt.Sprintf("error:%s", err))
		return nil, err
	} else {
		logger.InfoField("BackGround", "Child Process (Success)", fmt.Sprintf("Pid(%v)-->Process(%v)", os.Getegid(), cmd.Process.Pid))
	}
	if isExit {
		os.Exit(0)
	}

	return nil, nil
}

func (d *Deamon) Run() {
	c, err := BackGround(d.LogFile, true)
	if err != nil {
		logger.ErrorField("ERROR-MSG", "Run.err", err)
		return
	}
	if c != nil {
		logger.ErrorField("ERROR-MSG", "Run.c", c)
		//daemon process!!
		//daemon process should be a single instance app.
		err = d.s.Lock()
		if err != nil {
			logger.ErrorField("ERROR-Run", "Already", "Cannot lock file.May another instance running??")
			os.Exit(1)
		}
	}
	//守护进程启动一个子进程 并循环监控
	logger.InfoField("Run-Frist-SubProcess", "MSG", fmt.Sprintf("Start-Running(PID:%d)", os.Getpid()))
	var t int64
	count := 1

	errNum := 0

	for {
		Message := fmt.Sprintf("Start-Running(PID:%d,count:%d/%d,errNum:%d/%d)", os.Getpid(), count, d.MaxCount, errNum, d.MaxError)
		logger.InfoField("For Run-Frist-SubProcess", "MSG", Message)
		if errNum > d.MaxError {
			logger.InfoField("For Run-Frist-SubProcess", "ERROR", "启动子进程太多，退出")
			os.Exit(1)
		}
		if d.MaxCount > 0 && count > d.MaxCount {
			logger.InfoField("For Run-Frist-SubProcess", "ERROR", "重启太多次，退出")
			os.Exit(0)
		}

		count++
		t = time.Now().Unix() //start time 启动时间戳
		cmd, err := BackGround(d.LogFile, false)
		if err != nil {
			logger.ErrorField(Message, "ERROR", err)
			errNum++
			continue
		}
		//
		if cmd == nil {
			logger.InfoField("For Run-Frist-SubProcess", "StartRunning", fmt.Sprintf("(子进程PID:%d):Running .... ", os.Getpid()))
			break
		}
		//
		err = cmd.Wait()
		dat := time.Now().Unix() - t
		if dat < d.MinExitTime {
			errNum++
		} else {
			errNum = 0
		}
		EndMsg := fmt.Sprintf("%s监测到子进程(%d)退出吗,共运行了:%d秒:%v\n", Message, cmd.ProcessState.Pid(), dat, err)
		logger.InfoField("For Run-End", "EndMsg", EndMsg)
	}

}
