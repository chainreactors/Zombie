package exec

import (
	"Zombie/v1/pkg/utils"
	"github.com/mitchellh/go-vnc"
	"net"
	"strconv"
	"time"
)

type VNCService struct {
	utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
}

func (s *VNCService) Query() bool {
	return false
}

func (s *VNCService) GetInfo() bool {
	return false
}

func (s *VNCService) Connect() bool {
	err, res := VNCConnect(s.Username, s.Password, s.IpInfo)
	if err == nil && res {
		return true
	}
	return false

}

func (s *VNCService) DisConnect() bool {
	return false
}

func (s *VNCService) SetQuery(query string) {
	s.Input = query
}

func (s *VNCService) Output(res interface{}) {

}

func VNCConnect(User string, Password string, info utils.IpInfo) (err error, result bool) {

	targetPort := strconv.Itoa(info.Port)

	target := info.Ip + ":" + targetPort

	conn, err := net.DialTimeout("tcp", target, time.Duration(utils.Timeout)*time.Second)
	if err == nil {
		config := vnc.ClientConfig{
			Auth: []vnc.ClientAuth{
				&vnc.PasswordAuth{Password: Password},
			},
		}
		vncClient, err := vnc.Client(conn, &config)
		if err == nil {
			err = vncClient.Close()
			if err == nil {
				return err, true
			}
		}
	}
	return err, result
}
