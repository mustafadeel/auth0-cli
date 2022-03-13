//go:generate mockgen -source=client.go -destination=client_mock.go -package=auth0

package auth0

import "github.com/auth0/go-auth0/management"

type MultiFactorAPI interface {

	// Read a guardian policies.
	// See: https://auth0.com/docs/api/management/v2#!/Guardian/get_policies
	//Policy(id string, opts ...management.RequestOption) (*management.MultiFactorPolicies, error)

	Policy(...management.RequestOption) (*management.MultiFactorPolicies, error)

	UpdatePolicy(m *management.MultiFactorPolicies, opts ...management.RequestOption) (err error)

}
