package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nduyphuong/gorya/internal/store"
	"github.com/nduyphuong/gorya/internal/types"
	"github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
)

func ListPolicyV1alpha1(ctx context.Context, store store.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		verboseString := req.URL.Query().Get("verbose")
		var verbose bool
		if verboseString != "" {
			verbose = types.MustParseBool(verboseString)
		}
		policy, err := store.ListPolicy()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if verbose {
			resp := v1alpha1.ListPolicyResponsesVerbose{}
			for _, v := range *policy {
				resp = append(resp, v1alpha1.ListPolicyResponseVerbose{
					ListResponseVerbose: v1alpha1.ListResponseVerbose{
						Name:        v.Name,
						DisplayName: v.DisplayName,
					},
					Provider: v.Provider,
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
			for _, v := range *policy {
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
