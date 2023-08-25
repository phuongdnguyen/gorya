package handler

import (
	"context"
	"encoding/json"
	"github.com/nduyphuong/gorya/internal/store"
	"github.com/nduyphuong/gorya/internal/types"
	"github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"net/http"
)

func ListScheduleV1alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		verboseString := req.URL.Query().Get("verbose")
		var verbose bool
		if verboseString != "" {
			verbose = types.MustParseBool(verboseString)
		}
		schedules, err := store.ListSchedule()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if verbose {
			resp := v1alpha1.ListResponsesVerbose{}

			for _, v := range *schedules {
				resp = append(resp, v1alpha1.ListResponseVerbose{
					Name:        v.Name,
					DisplayName: v.DisplayName,
				})
			}
			b, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		} else {
			resp := v1alpha1.ListResponse{}
			for _, v := range *schedules {
				resp = append(resp, v.Name)
			}
			b, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		}
	}
}
