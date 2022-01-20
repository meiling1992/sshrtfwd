package main

import (
	"fmt"
	"os"
	"sshrts/api/logger"
	"sshrts/api/service"
	"sshrts/api/single"
	"sshrts/api/xdaemon"
	// log "github.com/sirupsen/logrus"
	// logg "sshrts/api/log"
)

// var log = logrus.New()

func main() {
	file, _ := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// // logg.DefaultLogger = file
	// logger := logg.NewStdLogger(file)
	// // loggertest := logg.NewStdLogger(os.Stdout)
	// // loggertesterr := logg.NewStdLogger(os.Stderr)
	// logger = logg.With(logg.MultiLogger(logger, loggertest), "test", logg.DefaultTinestamp, "Caller", logg.DefaultCaller)
	defer file.Close()
	sobject, err := single.New("sshrt")
	if err != nil {
		logger.ErrorField("Start-ERROR-MSG", "msg", "single new failed !")
		os.Exit(1)
	}
	logger.InfoField("Start-Test-MSG", "msg", fmt.Sprintf("Frist Isntance For New(PID:%d)", os.Getpid()))
	logFile := "daemon.log"
	daemon := xdaemon.NewDaemon(logFile, sobject)
	logger.InfoJsonField("Start-Test-MSG", "msg", daemon)
	daemon.Run()
	rtcfg, err := service.SSHRseTunnelCfg()
	if err != nil {
		logger.ErrorField("SSHRseTunnelCfg.Retrun.Error", "ERROR.MSG", fmt.Sprintf("SSH reverse tunnel get config failed :%v", err))
		os.Exit(10)
	}
	service.SSHRTService(rtcfg.SshServerEndPoint, rtcfg.LocalEndPoint, rtcfg.SshForwardEndPoint)

	// logger.Log(logg.LevelInfo, "msg", "values-main")
	// logger.Log(logg.LevelInfo, "msg", "values-test")
	// logger.Log(logg.LevelDebug, "mgsbug", "vaule-bug")
	// logger.Log(logg.ParseLevel("debug"), "mgsbugp", "vaule-bug-test")
	// // file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)

	// // if err == nil {
	// // 	log.Out = file
	// // } else {
	// // 	log.Info("Failed to log to file, using default stderr")
	// // }

	// // log.WithFields(logrus.Fields{
	// // 	"animal": "walrus",
	// // 	"size":   10,
	// // }).Info("A group of walrus emerges from the ocean")
	// logg.Info("sss")
	// logg.InfoW
	// 存放文章信息的 Post 结构体
	// type Book struct {
	// 	Id      int    `json:"id"`
	// 	Title   string `json:"title"`
	// 	Summary string `json:"summary"`
	// 	Author  string `json:"author"`
	// }
	// var books map[int]*Book = make(map[int]*Book)
	// book1 := Book{Id: 1, Title: "Go Web 编程", Summary: "Go Web 编程入门指南", Author: "学院君"}
	// books[book1.Id] = &book1
	// logger.InfoFields("ssss", logrus.Fields{"test": "sss"})
	// logger.InfoField("kkk", "test", "100")
	// logger.InfoJsonFields("mm", logrus.Fields{
	// 	"qq": books,
	// })
	// logger.InfoField("kkk", "test", "100")
}
