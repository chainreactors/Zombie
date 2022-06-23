package exec

import (
	utils2 "Zombie/v1/pkg/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type SshService struct {
	utils2.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	MysqlInf
	Cmd    string
	SshCon *ssh.Client
}

func (s *SshService) Connect() bool {
	err, _, conn := SSHConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SshCon = conn
		return true
	}
	return false
}

func (s *SshService) DisConnect() bool {
	s.SshCon.Close()
	return false
}

func (s *SshService) GetInfo() bool {

	if s.Cmd != "" {
		session, err := s.SshCon.NewSession()
		defer session.Close()
		defer s.SshCon.Close()
		cmd := "ping -c 5 " + s.Cmd
		buf, err := session.Output(cmd)

		if err != nil {
			return false
		}

		re, _ := regexp.Compile(`\d received`)

		FindRes := string(re.Find([]byte(buf)))

		reslist := strings.Split(FindRes, " ")
		if reslist[1] == "received" {
			if reslist[0] != "0" {
				fmt.Printf("%v can reach %v\n", s.Ip, s.Cmd)
			}
		}
	} else {
		panic("Please input ip")
	}

	return true
}

func (s *SshService) SetQuery(cmd string) {
	s.Cmd = cmd
}

func (s *SshService) Query() bool {

	session, err := s.SshCon.NewSession()
	defer session.Close()
	defer s.SshCon.Close()
	buf, err := session.Output(s.Cmd)

	if err != nil {
		return false
	}
	res := fmt.Sprintf(s.Ip + ":\n" + string(buf) + "\n")
	s.Output(res)
	return true
}

func (s *SshService) Output(res interface{}) {
	finres := res.(string)
	utils2.TDatach <- finres
}

func SSHConnect(User string, Password string, info utils2.IpInfo) (err error, result bool, connect *ssh.Client) {
	config := &ssh.ClientConfig{
		User: User,

		Timeout: time.Duration(utils2.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if strings.HasPrefix(Password, "pk:") {
		config.Auth = []ssh.AuthMethod{
			publicKeyAuthFunc(Password[3:]),
		}
	} else {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(Password),
		}
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", info.Ip, info.Port), config)
	if err == nil {
		session, err := client.NewSession()
		defer session.Close()
		errRet := session.Run("whoami")
		if err == nil && errRet == nil {
			result = true
		}
		connect = client
	}
	return err, result, connect
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {

	key, err := ioutil.ReadFile(kPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
