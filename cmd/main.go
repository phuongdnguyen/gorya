package main

import (
	"os"

	"github.com/nduyphuong/gorya/internal/signals"
	"github.com/nduyphuong/gorya/pkg/gce"

	"github.com/nduyphuong/gorya/internal/logging"
)

func main() {
	ctx := signals.SetupSignalHandler()
	gceClient, err := gce.NewGCEClient(ctx)
	if err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
	if err := Execute(ctx, gceClient); err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
}
