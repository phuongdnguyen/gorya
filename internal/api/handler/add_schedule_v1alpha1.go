package handler

import (
	"context"
	"encoding/json"
	"gorm.io/datatypes"
	"net/http"

	"github.com/nduyphuong/gorya/internal/models"
	"github.com/nduyphuong/gorya/internal/store"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
)

func notEmpty(a, b string) string {
	if b == "" {
		return a
	}
	return b
}
func AddScheduleV1Alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	resp := &svcv1alpha1.OkResponse{
		Message: "ok",
	}
	return func(w http.ResponseWriter, req *http.Request) {

		m := svcv1alpha1.AddScheduleRequest{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: converter func
		s := models.ScheduleModel{
			Name:        m.Name,
			DisplayName: notEmpty(m.Name, m.DisplayName),
			TimeZone:    m.TimeZone,
			Schedule: datatypes.NewJSONType(models.Schedule{
				Dtype:   m.Dtype,
				Corder:  m.Corder,
				Shape:   m.Shape,
				NdArray: m.NdArray,
			}),
		}
		if err := store.SaveSchedule(s); err != nil {
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
