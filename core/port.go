package core

import "fmt"

type Port struct {
	AgentID    int
	PortNumber int
}
type PortList []Port

func (a *PortList) String() string {
	s := ""
	sep := ""
	for _, v := range *a {
		s = s + fmt.Sprintf("%s[%d]%d", sep, v.PortNumber, v.AgentID)
		sep = ", "
	}

	return s
}

func (a *PortList) StringReversePort() string {
	s := ""
	sep := ""
	for _, v := range *a {
		s = s + fmt.Sprintf("%s%d[%d]", sep, v.AgentID, v.PortNumber)
		sep = ", "
	}

	return s
}
