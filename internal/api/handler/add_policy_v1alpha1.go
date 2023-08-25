package handler

import (
	"context"
	"encoding/json"
	"github.com/nduyphuong/gorya/internal/models"
	"github.com/nduyphuong/gorya/internal/store"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"gorm.io/datatypes"
	"net/http"
)

func AddPolicyV1Alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	resp := &svcv1alpha1.OkResponse{
		Message: "ok",
	}
	return func(w http.ResponseWriter, req *http.Request) {
		m := svcv1alpha1.AddPolicyRequest{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		schedule, err := store.GetSchedule(m.ScheduleName)
		if schedule == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s := models.Policy{
			Name:         m.Name,
			DisplayName:  notEmpty(m.Name, m.DisplayName),
			Projects:     datatypes.NewJSONSlice(m.Projects),
			Tags:         datatypes.NewJSONSlice(m.Tags),
			ScheduleName: m.ScheduleName,
		}
		if err := store.SavePolicy(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
