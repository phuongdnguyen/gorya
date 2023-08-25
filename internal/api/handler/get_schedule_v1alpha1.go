package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nduyphuong/gorya/internal/store"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"gorm.io/gorm"
	"net/http"
)

func isEmpty(s string) bool {
	return s == ""
}
func GetScheduleV1alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("schedule")
		if isEmpty(name) {
			http.Error(w, errors.New("empty schedule name").Error(), http.StatusBadRequest)
			return
		}
		schedule, err := store.GetSchedule(name)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp := svcv1alpha1.GetScheduleResponse{
			Name:        schedule.Name,
			DisplayName: notEmpty(schedule.Name, schedule.DisplayName),
			TimeZone:    schedule.TimeZone,
			Dtype:       schedule.Schedule.Data().Dtype,
			Corder:      schedule.Schedule.Data().Corder,
			NdArray:     schedule.Schedule.Data().NdArray,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
