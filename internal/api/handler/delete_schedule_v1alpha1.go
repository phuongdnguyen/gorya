package handler

import (
	"context"
	"errors"
	"github.com/nduyphuong/gorya/internal/store"
	"net/http"
)

func DeleteScheduleV1alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("schedule")
		if isEmpty(name) {
			http.Error(w, errors.New("empty schedule name").Error(), http.StatusBadRequest)
			return
		}
		if isUnusedSchedule(name, store) {
			if err := store.DeleteSchedule(name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}

func isUnusedSchedule(name string, store store.Interface) bool {
	policies, err := store.GetPolicyBySchedule(name)
	if err != nil {
		return false
	}
	for _, policy := range *policies {
		if policy.ScheduleName == name {
			//schedule in use by at least 1 policy, must  not delete
			return false
		}
	}
	return true
}
