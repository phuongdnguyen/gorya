package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	constants "github.com/nduyphuong/gorya/internal/constants"
	svcv1alpha1 "github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
	"github.com/nduyphuong/gorya/pkg/aws"
	"github.com/nduyphuong/gorya/pkg/gcp"
	pkgerrors "github.com/pkg/errors"
)

/*
*
For every request, naively we will init a new:
  - EndpointResolverWithOptions to be able to emulate with localstack
    and return aws config in order to pass to its super.

Problem: if there are N in queue -> N request to aws although there will be k amount of request will use the same
AssumeRoleProvider (k<=N)

Optimization:
- We can init a client pool identified by AssumeRoleARN
*/
func ChangeStateV1alpha1(ctx context.Context, awsClientPool *aws.ClientPool, gcpClientPool *gcp.ClientPool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		m := svcv1alpha1.ChangeStateRequest{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			http.Error(w, pkgerrors.Wrap(err, "decode state change request body").Error(), http.StatusBadRequest)
			return
		}
		if awsClientPool != nil && m.Provider == constants.PROVIDER_AWS {
			awsClient, ok := awsClientPool.GetForCredential(m.CredentialRef)
			if !ok {
				http.Error(w, fmt.Errorf("client not found for credential %v", m.CredentialRef).Error(),
					http.StatusBadRequest)
				return
			}
			compute := awsClient.EC2()
			if err := compute.ChangeStatus(ctx, m.Action, m.TagKey, m.TagValue); err != nil {
				http.Error(w, pkgerrors.Wrap(err, "change compute status").Error(), http.StatusInternalServerError)
				return
			}
		}
		if gcpClientPool != nil && m.Provider == constants.PROVIDER_GCP {
			gcpClient, ok := gcpClientPool.GetForCredential(m.CredentialRef)
			if !ok {
				http.Error(w, fmt.Errorf("client not found for credential %v", m.CredentialRef).Error(),
					http.StatusBadRequest)
				return
			}
			compute := gcpClient.GCE()
			if err := compute.ChangeStatus(ctx, m.Action, m.TagKey, m.TagValue); err != nil {
				http.Error(w, pkgerrors.Wrap(err, "change compute status").Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
