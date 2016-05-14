package config

import "fmt"

type Port struct {
	AgentName  string
	PortNumber int
}
type PortList []Port

func (a *PortList) String() string {
	s := ""
	sep := ""
	for _, v := range *a {
		s = s + fmt.Sprintf("%s[%d]%s", sep, v.PortNumber, v.AgentName)
		sep = ", "
	}

	return s
}

func (a *PortList) StringReversePort() string {
	s := ""
	sep := ""
	for _, v := range *a {
		s = s + fmt.Sprintf("%s%s[%d]", sep, v.AgentName, v.PortNumber)
		sep = ", "
	}

	return s
}
