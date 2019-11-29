package controller

import (
	"github.com/cjheppell/bouncer/pkg/controller/tcpservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, tcpservice.Add)
}
