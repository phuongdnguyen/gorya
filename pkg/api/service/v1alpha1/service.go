package v1alpha1

import (
	"context"
	"net/http"

	"github.com/nduyphuong/gorya/internal/store"
)

//go:generate mockery --name GoryaServiceHandler
type GoryaServiceHandler interface {
	GetTimeZone() http.Handler
	GetVersionInfo() http.Handler
	AddSchedule(ctx context.Context) http.Handler
	GetSchedule(ctx context.Context) http.Handler
	ListSchedule(ctx context.Context) http.Handler
	DeleteSchedule(ctx context.Context) http.Handler
	AddPolicy(ctx context.Context) http.Handler
	GetPolicy(ctx context.Context) http.Handler
	ListPolicy(ctx context.Context) http.Handler
	DeletePolicy(ctx context.Context) http.Handler
	ChangeState(ctx context.Context) http.Handler
	ScheduleTask(ctx context.Context) http.Handler
}

const (
	GoryaTaskChangeStageProcedure = "/tasks/change_state"
	GoryaTaskScheduleProcedure    = "/tasks/schedule"
	GoryaGetTimeZoneProcedure     = "/api/v1alpha1/time_zones"
	GoryaAddScheduleProcedure     = "/api/v1alpha1/add_schedule"
	GoryaGetScheduleProcedure     = "/api/v1alpha1/get_schedule"
	GoryaListScheduleProcedure    = "/api/v1alpha1/list_schedules"
	GoryaDeleteScheduleProcedure  = "/api/v1alpha1/del_schedule"
	GoryaAddPolicyProcedure       = "/api/v1alpha1/add_policy"
	GoryaGetPolicyProcedure       = "/api/v1alpha1/get_policy"
	GoryaListPolicyProcedure      = "/api/v1alpha1/list_policies"
	GoryaDeletePolicyProcedure    = "/api/v1alpha1/del_policy"
	GoryaGetVersionInfo           = "/api/v1alpha1/version_info"
)

// NewGoryaServiceHandler builds an HTTP handler from the service implementation. It returns the
//
//	path on which to mount the handler and the handler itself.
//
//	 https://stackoverflow.com/questions/33646948/go-using-mux-router-how-to-pass-my-db-to-my-handlers
func NewGoryaServiceHandler(ctx context.Context, store store.Interface, svc GoryaServiceHandler) (string,
	http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(GoryaGetTimeZoneProcedure, svc.GetTimeZone())
	mux.Handle(GoryaGetVersionInfo, svc.GetVersionInfo())
	mux.Handle(GoryaAddScheduleProcedure, svc.AddSchedule(ctx))
	mux.Handle(GoryaGetScheduleProcedure, svc.GetSchedule(ctx))
	mux.Handle(GoryaListScheduleProcedure, svc.ListSchedule(ctx))
	mux.Handle(GoryaDeleteScheduleProcedure, svc.DeleteSchedule(ctx))
	mux.Handle(GoryaAddPolicyProcedure, svc.AddPolicy(ctx))
	mux.Handle(GoryaGetPolicyProcedure, svc.GetPolicy(ctx))
	mux.Handle(GoryaListPolicyProcedure, svc.ListPolicy(ctx))
	mux.Handle(GoryaDeletePolicyProcedure, svc.DeletePolicy(ctx))
	mux.Handle(GoryaTaskChangeStageProcedure, svc.ChangeState(ctx))
	mux.Handle(GoryaTaskScheduleProcedure, svc.ScheduleTask(ctx))
	return "/", mux
}
