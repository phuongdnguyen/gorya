package handler

import (
	"context"
	"errors"
	"github.com/nduyphuong/gorya/internal/store"
	"net/http"
)

func DeletePolicyV1alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("policy")
		if isEmpty(name) {
			http.Error(w, errors.New("empty policy name").Error(), http.StatusBadRequest)
			return
		}
		if err := store.DeletePolicy(name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
