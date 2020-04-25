package controller

import (
	"github.com/kube-sailmaker/request-operator/pkg/controller/releaserequest"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, releaserequest.Add)
}
