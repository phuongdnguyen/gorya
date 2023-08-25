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

func GetPolicyV1Alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("policy")
		if isEmpty(name) {
			http.Error(w, errors.New("empty policy name").Error(), http.StatusBadRequest)
			return
		}
		policy, err := store.GetPolicyByName(name)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp := svcv1alpha1.GetPolicyResponse{
			Policy: *policy,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
