package auth

import (
	"context"

	"google.golang.org/grpc/credentials"
)

const scheme = "bearer"

var _ credentials.PerRPCCredentials = &GRPCCredentials{}

// GRPCCredentials implements PerRPCCredentials.
type GRPCCredentials struct {
	Token string

	TransportSecurity bool
}

// GetRequestMetadata maps the given credentials to the appropriate request
// headers.
func (c GRPCCredentials) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": scheme + " " + c.Token,
	}, nil
}

// RequireTransportSecurity implements PerRPCCredentials.
func (c GRPCCredentials) RequireTransportSecurity() bool {
	return c.TransportSecurity
}
