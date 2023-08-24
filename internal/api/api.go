package api

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/nduyphuong/gorya/internal/api/config"
	"github.com/nduyphuong/gorya/internal/api/handler"
	"github.com/nduyphuong/gorya/internal/logging"
	"github.com/nduyphuong/gorya/internal/os"
	queueOptions "github.com/nduyphuong/gorya/internal/queue/options"
	"github.com/nduyphuong/gorya/internal/store"
	"github.com/nduyphuong/gorya/internal/version"
	"github.com/nduyphuong/gorya/internal/worker"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"github.com/nduyphuong/gorya/pkg/aws"
	awsOptions "github.com/nduyphuong/gorya/pkg/aws/options"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	PROVIDER_AWS = "aws"
	PROVIDER_GCP = "gcp"
)

type Server interface {
	Serve(ctx context.Context, l net.Listener) error
}

type server struct {
	cfg           config.ServerConfig
	sc            store.Interface
	aws           *aws.AwsPool
	taskProcessor worker.Interface
}

func NewServer(cfg config.ServerConfig) (Server, error) {
	return &server{
		cfg: cfg,
	}, nil
}

type CredentialRef struct {
	credentialRef map[string]bool
	lock          sync.Mutex
}

func (s *server) Serve(ctx context.Context, l net.Listener) error {
	var err error
	errCh := make(chan error)
	log := logging.LoggerFromContext(ctx)
	log.Infof("Server is listening on %q", l.Addr().String())
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})
	s.sc, err = store.GetOnce()
	if err != nil {
		return err
	}
	//TODO: change to sync map

	c := CredentialRef{}
	c.credentialRef = map[string]bool{}
	c.credentialRef[aws.Default] = true
	awsPool := aws.AwsPool{}
	ticker := time.NewTicker(30 * time.Second)
	updateAWSClientPool := func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		policies, err := s.sc.ListPolicyByProvider(PROVIDER_AWS)
		for _, pol := range *policies {
			for _, project := range pol.Projects {
				if project.CredentialRef != "" {
					c.credentialRef[project.CredentialRef] = true
				}
			}
		}
		//c.credentialRef["arn:aws:iam::043159268388:role/test"] = true
		s.aws, err = awsPool.New(
			ctx,
			c.credentialRef,
			//currently we only support multi account in 1 region
			awsOptions.WithRegion(os.GetEnv("AWS_REGION", "ap-southeast-1")),
			awsOptions.WithEndpoint(os.GetEnv("AWS_ENDPOINT", "")),
		)
		if err != nil {
			log.Errorf("update aws client pool %v", err)
			return
		}
	}
	//we update client pool for the first time then periodically update it every ticker.C seconds
	updateAWSClientPool()
	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				updateAWSClientPool()
			}
		}

	}(ctx.Done())

	s.taskProcessor = worker.NewClient(worker.Options{
		QueueOpts: queueOptions.Options{
			Addr:        os.GetEnv("GORYA_REDIS_ADDR", "localhost:6379"),
			Name:        os.GetEnv("GORYA_QUEUE_NAME", "gorya"),
			PopInterval: 2 * time.Second,
		},
	})
	path, svcHandler := svcv1alpha1.NewGoryaServiceHandler(ctx, s.sc, s)
	mux.Handle(path, svcHandler)
	srv := &http.Server{
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: time.Minute,
	}
	go func() { errCh <- srv.Serve(l) }()
	select {
	case <-ctx.Done():
		log.Info("Gracefully stopping server...")
		time.Sleep(s.cfg.GracefulShutdownTimeout)
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func (s *server) GetTimeZone() http.Handler {
	return handler.GetTimeZoneV1Alpha1()
}

func (s *server) GetVersionInfo() http.Handler {
	return handler.GetVersionInfoV1Alpha1(version.GetVersion())
}

func (s *server) AddSchedule(ctx context.Context) http.Handler {
	return handler.AddScheduleV1Alpha1(ctx, s.sc)
}

func (s *server) GetSchedule(ctx context.Context) http.Handler {
	return handler.GetScheduleV1alpha1(ctx, s.sc)
}

func (s *server) ListSchedule(ctx context.Context) http.Handler {
	return handler.ListScheduleV1alpha1(ctx, s.sc)
}

func (s *server) DeleteSchedule(ctx context.Context) http.Handler {
	return handler.DeleteScheduleV1alpha1(ctx, s.sc)
}

func (s *server) AddPolicy(ctx context.Context) http.Handler {
	return handler.AddPolicyV1Alpha1(ctx, s.sc)
}

func (s *server) GetPolicy(ctx context.Context) http.Handler {
	return handler.GetPolicyV1Alpha1(ctx, s.sc)
}

func (s *server) ListPolicy(ctx context.Context) http.Handler {
	return handler.ListPolicyV1alpha1(ctx, s.sc)
}

func (s *server) DeletePolicy(ctx context.Context) http.Handler {
	return handler.DeletePolicyV1alpha1(ctx, s.sc)
}

func (s *server) ChangeState(ctx context.Context) http.Handler {
	return handler.ChangeStateV1alpha1(ctx, s.aws)
}

func (s *server) ScheduleTask(ctx context.Context) http.Handler {
	return handler.ScheduleTaskV1alpha1(ctx, s.sc, s.taskProcessor)
}
