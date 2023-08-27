package v1alpha1

import (
	"context"
	"github.com/nduyphuong/gorya/internal/api/middleware"
	"net/http"

	"github.com/nduyphuong/gorya/internal/store"
)

//go:generate mockery --name GoryaServiceHandler
type GoryaServiceHandler interface {
	GetTimeZone() http.HandlerFunc
	GetVersionInfo() http.HandlerFunc
	AddSchedule(ctx context.Context) http.HandlerFunc
	GetSchedule(ctx context.Context) http.HandlerFunc
	ListSchedule(ctx context.Context) http.HandlerFunc
	DeleteSchedule(ctx context.Context) http.HandlerFunc
	AddPolicy(ctx context.Context) http.HandlerFunc
	GetPolicy(ctx context.Context) http.HandlerFunc
	ListPolicy(ctx context.Context) http.HandlerFunc
	DeletePolicy(ctx context.Context) http.HandlerFunc
	ChangeState(ctx context.Context) http.HandlerFunc
	ScheduleTask(ctx context.Context) http.HandlerFunc
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
	mux.Handle(GoryaGetTimeZoneProcedure, middleware.JWTAuthorization(svc.GetTimeZone(), "get-timezone", ctx))
	mux.Handle(GoryaGetVersionInfo, middleware.JWTAuthorization(svc.GetVersionInfo(), "get-version-info", ctx))
	mux.Handle(GoryaAddScheduleProcedure, middleware.JWTAuthorization(svc.AddSchedule(ctx), "add-schedule", ctx))
	mux.Handle(GoryaGetScheduleProcedure, middleware.JWTAuthorization(svc.GetSchedule(ctx), "get-schedule", ctx))
	mux.Handle(GoryaListScheduleProcedure, middleware.JWTAuthorization(svc.ListSchedule(ctx), "list-schedule", ctx))
	mux.Handle(GoryaDeleteScheduleProcedure, middleware.JWTAuthorization(svc.DeleteSchedule(ctx), "delete-schedule", ctx))
	mux.Handle(GoryaAddPolicyProcedure, middleware.JWTAuthorization(svc.AddPolicy(ctx), "add-policy", ctx))
	mux.Handle(GoryaGetPolicyProcedure, middleware.JWTAuthorization(svc.GetPolicy(ctx), "get-policy", ctx))
	mux.Handle(GoryaListPolicyProcedure, middleware.JWTAuthorization(svc.ListPolicy(ctx), "list-policy", ctx))
	mux.Handle(GoryaDeletePolicyProcedure, middleware.JWTAuthorization(svc.DeletePolicy(ctx), "delete-policy", ctx))
	mux.Handle(GoryaTaskChangeStageProcedure, middleware.JWTAuthorization(svc.ChangeState(ctx), "change-state", ctx))
	mux.Handle(GoryaTaskScheduleProcedure, middleware.JWTAuthorization(svc.ScheduleTask(ctx), "schedule-task", ctx))
	return "/", mux
}
