package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/nduyphuong/gorya/internal/api"
	"github.com/nduyphuong/gorya/internal/api/config"
	"github.com/nduyphuong/gorya/internal/os"
	queueOptions "github.com/nduyphuong/gorya/internal/queue/options"
	"github.com/nduyphuong/gorya/internal/types"
	versionpkg "github.com/nduyphuong/gorya/internal/version"
	"github.com/nduyphuong/gorya/internal/worker"
	"github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	pkgerrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "api",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			version := versionpkg.GetVersion()
			log.WithFields(log.Fields{
				"version": version.Version,
				"commit":  version.GitCommit,
			}).Info("Starting Gorya API Server")
			var wg sync.WaitGroup
			errCh := make(chan error, 2)
			taskProcessor := worker.NewClient(worker.Options{
				QueueOpts: queueOptions.Options{
					Name: os.GetEnv("GORYA_QUEUE_NAME", "gorya"),
					Addr: os.GetEnv("GORYA_REDIS_ADDR", "localhost:6379"),
					//check in queue every 5 seconds
					PopInterval: 5 * time.Second,
				},
			})
			ticker := time.NewTicker(2 * time.Second)
			taskProcessResultChan := make(chan string)
			numWorkers := types.MustParseInt(os.GetEnv("GORYA_NUM_WORKER", "1"))
			for i := 0; i <= numWorkers; i++ {
				go func(stop <-chan struct{}) {
					for {
						select {
						case <-stop:
							return
						case <-ticker.C:
							taskProcessor.Process(ctx, taskProcessResultChan, errCh)
						}
					}
				}(ctx.Done())
			}
			//process result from channel
			go func() {
				for task := range taskProcessResultChan {
					log.Infof("popped item %v", task)
					var elem worker.QueueElem
					err := json.Unmarshal([]byte(task), &elem)
					if err != nil {
						log.Errorf("unmarshal elem from queue %v", err)
						return
					}
					changeStateRequest := v1alpha1.ChangeStateRequest{
						Action:   elem.Action,
						Project:  elem.Project,
						TagKey:   elem.TagKey,
						TagValue: elem.TagValue,
					}
					requestURL := fmt.Sprintf("http://localhost:%d%s", types.MustParseInt(os.GetEnv("PORT",
						"8080")), v1alpha1.GoryaTaskChangeStageProcedure)
					b, err := json.Marshal(changeStateRequest)
					if err != nil {
						log.Errorf("unmarshal changeStateRequest %v", err)
						return
					}
					req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(b))
					if err != nil {
						log.Errorf("creating request %v", err)
						return
					}
					req.Header.Set("Content-Type", "application/json")
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Errorf("making request %v", err)
					}
					if resp.StatusCode != http.StatusOK {
						log.Errorf("change state request failed with status code: %v", resp.StatusCode)
					}
				}
			}()
			//periodically check if there is action to be done
			go func(stop <-chan struct{}) {
				for {
					select {
					case <-stop:
						return
					case <-ticker.C:
						requestURL := fmt.Sprintf("http://localhost:%d%s", types.MustParseInt(os.GetEnv("PORT",
							"8080")), v1alpha1.GoryaTaskScheduleProcedure)
						req, err := http.NewRequest(http.MethodGet, requestURL, nil)
						if err != nil {
							errCh <- pkgerrors.Wrap(err, "creating request")
							return
						}
						_, err = http.DefaultClient.Do(req)
						if err != nil {
							errCh <- pkgerrors.Wrap(err, "making request")
						}
					}
				}
			}(ctx.Done())
			cfg := config.ServerConfigFromEnv()
			srv, err := api.NewServer(cfg)
			if err != nil {
				return pkgerrors.Wrap(err, "error creating API server")
			}
			l, err := net.Listen(
				"tcp",
				fmt.Sprintf(
					"%s:%s",
					os.GetEnv("HOST", "0.0.0.0"),
					os.GetEnv("PORT", "8080"),
				),
			)
			if err != nil {
				return pkgerrors.Wrap(err, "error creating listener")
			}
			defer func() {
				_ = l.Close()
			}()
			wg.Add(1)
			go func() {
				srvErr := srv.Serve(ctx, l)
				errCh <- pkgerrors.Wrap(srvErr, "serve")
				wg.Done()
			}()
			wg.Wait()
			close(errCh)
			var resErr error
			for err := range errCh {
				resErr = errors.Join(resErr, err)
			}
			return resErr
		},
	}
}
