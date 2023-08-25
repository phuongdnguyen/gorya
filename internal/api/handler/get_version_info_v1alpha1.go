package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nduyphuong/gorya/internal/version"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
)

func GetVersionInfoV1Alpha1(v version.Version) http.HandlerFunc {
	resp := &svcv1alpha1.GetVersionInfoResponse{
		VersionInfo: &svcv1alpha1.VersionInfo{
			Version:      v.Version,
			GitCommit:    v.GitCommit,
			GitTreeDirty: v.GitTreeDirty,
			BuildTime:    v.BuildDate,
			GoVersion:    v.GoVersion,
			Compiler:     v.Compiler,
			Platform:     v.Platform,
		},
	}
	return func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
