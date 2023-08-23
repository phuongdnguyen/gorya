package handler

import (
	"context"
	"encoding/json"
	"net/http"

	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"github.com/nduyphuong/gorya/pkg/aws"
	pkgerrors "github.com/pkg/errors"
)

func ChangeStateV1alpha1(ctx context.Context, awsClient aws.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		compute := awsClient.EC2()
		m := svcv1alpha1.ChangeStateRequest{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, pkgerrors.Wrap(err, "decode state change request body").Error(), http.StatusBadRequest)
			return
		}
		if err := compute.ChangeStatus(ctx, m.Action, m.TagKey, m.TagValue); err != nil {
			http.Error(w, pkgerrors.Wrap(err, "change compute status").Error(), http.StatusInternalServerError)
			return
		}
	}
}
