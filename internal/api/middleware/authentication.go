package middleware

import (
	"context"
	"crypto/tls"
	"github.com/coreos/go-oidc"
	"github.com/nduyphuong/gorya/internal/os"
	"net/http"
	"time"
)

// Claims claims component of jwt contains many fields , we need only roles of GoryaServiceClient
// "GoryaServiceClient":{"GoryaServiceClient":{"roles":["get-timezone","list-policy","delete-policy"]}},
type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	GoryaServiceClient clientRoles `json:"gorya,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

var (
	realmConfigUrl        = os.GetEnv("GORYA_KEYCLOAK_REALM_URL", "http://localhost:8080/auth/realms/demorealm")
	clientID       string = os.GetEnv("GORYA_KEYCLOAK_CLIENT_ID",
		"GoryaServiceClient")
)

func JWTAuthorization(h http.HandlerFunc, role string, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rawAccessToken := r.Header.Get("Authorization")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(ctx, client)
		provider, err := oidc.NewProvider(ctx, realmConfigUrl)
		if err != nil {
			authorisationFailed("authorisation failed while getting the provider: "+err.Error(), w, r)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: clientID,
			//skip check aud in jwt payload since in keycloak it's always `account`
			SkipClientIDCheck: true,
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)
		if err != nil {
			authorisationFailed("authorisation failed while verifying the token: "+err.Error(), w, r)
			return
		}

		var IDTokenClaims Claims
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorisationFailed("claims : "+err.Error(), w, r)
			return
		}
		//checking the roles
		user_access_roles := IDTokenClaims.ResourceAccess.GoryaServiceClient.Roles
		for _, b := range user_access_roles {
			if b == role {
				h(w, r)
				return
			}
		}

		authorisationFailed("user not allowed to access this api", w, r)
	}
}

func authorisationFailed(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}
