package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sshrts/api/logger"

	"github.com/astaxie/beego/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type EndPoint struct {
	Host string
	Port int
}

type SSHReverseTunnelCfg struct {
	LocalEndPoint      EndPoint
	SshServerEndPoint  EndPoint
	SshForwardEndPoint EndPoint
}

func (endpoint *EndPoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
func bridgeTunnel(client net.Conn, remote net.Conn) {
	defer client.Close()
	DoneCh := make(chan bool)
	// Start remote -> local data transfer

	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			logger.ErrorField("bridgeTunnel-ERROR", "ERROR", fmt.Sprintf("error: while copy remtoe->local:%s", err))

		}
		DoneCh <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			logger.ErrorField("bridgeTunnel-ERROR", "ERROR", fmt.Sprintf("error: while copy local->remtoe:%s", err))

		}
		DoneCh <- true
	}()
	<-DoneCh
}

func publicKeyFile(file string) ssh.AuthMethod {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		logger.ErrorField("publicKeyFile-ReadFile", "Error", fmt.Sprintf("Cannot read SSH pubilc key file:%s", err))
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		logger.ErrorField("publicKeyFile-ParsePrivateKey", "Error", fmt.Sprintf("Cannot Parse SSH pubilc key file:%s", err))
		return nil
	}
	return ssh.PublicKeys(key)
}

func SSHRseTunnelCfg() (*SSHReverseTunnelCfg, error) {
	cfgFile, err := config.NewConfig("ini", "sshrt.config")
	if err != nil {
		logger.ErrorField("SSHReverseTunnelCfg-Config.NewConfig", "ERROR", fmt.Sprintf("config.file.error:%s", err))
		return nil, err
	}
	remoteip := cfgFile.String("sshrt::remoteip")
	remoteport, err := cfgFile.Int("sshrt::remoteport")
	if err != nil {
		remoteport = 80
	}
	// MSG_RemoteIp := fmt.Sprintf("remoteip:%s", remoteip)
	// logger.InfoJsonField("SSHRseTunnelCfg.INFO", "SSHRseTunnelCfg-MSG_RemoteIp", remoteip)
	logger.InfoJsonFields("SSHRseTunnelCfg.INFO.MSG_Remote.Date", logrus.Fields{
		"MSG_RemoteIp":   remoteip,
		"MSG_RemotePort": remoteport,
	})
	sship := cfgFile.String("sshrt::sship")
	sshport, err := cfgFile.Int("sshrt::sshport")
	if err != nil {
		sshport = 80
	}
	logger.InfoJsonFields("SSHRseTunnelCfg.INFO.MSG_SSH.Date", logrus.Fields{
		"MSG_SSHIp":   sship,
		"MSG_SSHPort": sshport,
	})
	sshfwdip := cfgFile.String("sshrt::sshfwdip")
	if sshfwdip == "" {
		sshfwdip = "127.0.0.1"
	}
	sshfwdport, err := cfgFile.Int("sshrt::sshfwdport")
	if err != nil {
		sshfwdport = 80
	}
	logger.InfoJsonFields("SSHRseTunnelCfg.INFO.MSG_SSHFwd.Date", logrus.Fields{
		"MSG_SSHFwdIp":   sshfwdip,
		"MSG_SSHFwdPort": sshfwdport,
	})

	if remoteip == "" || sship == "" {
		logger.ErrorField("SSHRseTunnelCfg.Conf.Error", "ERROR", "remote and ssh ip must be configd.")
		return nil, err
	}

	return &SSHReverseTunnelCfg{
		LocalEndPoint: EndPoint{
			Host: remoteip,
			Port: remoteport,
		},
		SshServerEndPoint: EndPoint{
			Host: sship,
			Port: sshport,
		},
		SshForwardEndPoint: EndPoint{
			Host: sshfwdip,
			Port: sshfwdport,
		},
	}, nil
}

func SSHRTService(sshsvrep, lep, fwdep EndPoint) error {

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			publicKeyFile("./key/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Listen on remote server port
	serverConn, err := ssh.Dial("tcp", sshsvrep.String(), sshConfig)
	if err != nil {
		logger.ErrorField("ERROR.SSH", "Dial.MSG", fmt.Sprintf("Dial INFO remote server.error:%s", err))
		return err
	}
	defer serverConn.Close()
	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", fwdep.String())
	if err != nil {
		logger.ErrorField("ERROR.SSH", "Listen.MSG", fmt.Sprintf("Listen open port on remote server.error:%s", err))
		return err
	}
	defer listener.Close()
	// handle incoming connections on reverse forwarded tunnel
	for {
		// accept a new tcp session on tunnel
		client, err := listener.Accept()
		if err != nil {
			logger.ErrorField("ERROR.SSH", "Listen.Accept.MSG", fmt.Sprintf("Listen.Accept server.error:%s", err))
			continue
		}
		//
		local, err := net.Dial("tcp", lep.String())
		if err != nil {
			logger.ErrorField("ERROR.SSH", "Dial(local).MSG", fmt.Sprintf("Dial local  server.error:%s", err))
			continue
		}
		go bridgeTunnel(client, local)
	}
	// return nil
}
