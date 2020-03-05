package controller

import (
	"image-clone-controller/pkg/controller/daemonset"
	"image-clone-controller/pkg/controller/deployment"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, daemonset.Add)
	AddToManagerFuncs = append(AddToManagerFuncs, deployment.Add)
}

// AddToManagerFuncs is a list of functions to add all controllers to the manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all controllers to the manager
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
