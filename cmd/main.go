package main

import (
	"flag"
	"os"

	"image-clone-controller/pkg/config"
	"image-clone-controller/pkg/controller"

	"github.com/spf13/pflag"
	ctrconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

const (
	leaderLockConfigMap = "image-clone-controller"
)

var (
	log = logf.Log.WithName("main")
)

func main() {
	// add flags registered by all imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := config.GlobalConfig.Validate()
	if err != nil {
		panic(err)
	}

	logf.SetLogger(zap.Logger(false))
	log.Info("Starting image clone controller")

	cfg, err := ctrconfig.GetConfig()
	if err != nil {
		log.Error(err, "Failed to get KubeConfig")
		os.Exit(1)
	}

	mgr, err := manager.New(cfg, manager.Options{
		// leader election: not strictly necessary for this test task
		// idea was to prevent the concurrent event handling during RollingUpgrade
		LeaderElection:     false,
		LeaderElectionID:   leaderLockConfigMap,
		MetricsBindAddress: "0",
	})
	if err != nil {
		log.Error(err, "Failed to create the manager")
		os.Exit(1)
	}

	log.Info("Registering controllers")
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "Failed to register all the controllers")
		os.Exit(1)
	}

	log.Info("Starting controllers")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}
