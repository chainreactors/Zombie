package ldap

import (
	"github.com/chainreactors/zombie/pkg"
	ldap "github.com/go-ldap/ldap/v3"
)

type LdapPlugin struct {
	*pkg.Task
	Input string
	conn  *ldap.Conn
}

func (s *LdapPlugin) Unauth() (bool, error) {
	//TODO implement me
	return false, nil
}

//func (s *LdapPlugin) Query() bool {
//	panic("implement me")
//}

func (s *LdapPlugin) Login() error {
	var conn *ldap.Conn
	ldap.DefaultTimeout = s.Duration()
	conn, err := ldap.Dial("tcp", s.Address())

	if err != nil {
		return err
	}

	err = conn.Bind(s.Username, s.Password)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *LdapPlugin) Close() error {
	if s.conn != nil {
		s.conn.Close()
		return nil
	}
	return pkg.NilConnError{s.Service}
}

func (s *LdapPlugin) Name() string {
	return s.Service.String()
}

func (s *LdapPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

//func (s *LdapPlugin) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *LdapPlugin) Output(res interface{}) {
//
//}
//
//func (s *LdapPlugin) GetInfo() bool {
//	s.conn.Close()
//	return true
//}
