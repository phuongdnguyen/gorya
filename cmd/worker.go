package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nduyphuong/gorya/internal/os"
	queueOptions "github.com/nduyphuong/gorya/internal/queue/options"
	"github.com/nduyphuong/gorya/internal/types"
	"github.com/nduyphuong/gorya/internal/worker"
	"github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newWorkCommand() *cobra.Command {
	return &cobra.Command{
		Use:               "worker",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			errCh := make(chan error, 2)
			defer close(errCh)
			taskProcessResultChan := make(chan string)
			taskProcessor := worker.NewClient(worker.Options{
				QueueOpts: queueOptions.Options{
					Addr:        os.GetEnv("GORYA_REDIS_ADDR", "localhost:6379"),
					Name:        os.GetEnv("GORYA_QUEUE_NAME", "gorya"),
					PopInterval: 2 * time.Second,
				},
			})

			numWorkers := types.MustParseInt(os.GetEnv("GORYA_NUM_WORKER", "1"))
			for i := 0; i < numWorkers; i++ {
				go taskProcessor.Process(ctx, ctx.Done(), errCh)
			}
			for task := range taskProcessResultChan {
				var elem worker.QueueElem
				err := json.Unmarshal([]byte(task), &elem)
				if err != nil {
					return pkgerrors.Wrap(err, "unmarshal elem from queue")
				}
				changeStateRequest := v1alpha1.ChangeStateRequest{
					Action:        elem.Action,
					Project:       elem.Project,
					TagKey:        elem.TagKey,
					TagValue:      elem.TagValue,
					CredentialRef: elem.CredentialRef,
					Provider:      elem.Project,
				}
				requestURL := fmt.Sprintf("http://localhost:%d%s", types.MustParseInt(os.GetEnv("PORT",
					"8080")), v1alpha1.GoryaTaskChangeStageProcedure)
				b, err := json.Marshal(changeStateRequest)
				if err != nil {
					return pkgerrors.Wrap(err, "unmarshal changeStateRequest")
				}
				req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(b))
				if err != nil {
					return pkgerrors.Wrap(err, "creating request")
				}
				req.Header.Set("Content-Type", "application/json")
				_, err = http.DefaultClient.Do(req)
				if err != nil {
					return pkgerrors.Wrap(err, "making request")
				}
			}
			close(errCh)
			var resErr error
			for err := range errCh {
				resErr = errors.Join(resErr, err)
			}
			return resErr
		},
	}
}
