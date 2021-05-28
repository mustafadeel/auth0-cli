package client

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/auth0/logger0/rpc/auth"
	controlv1 "github.com/auth0/logger0/rpc/control/v1"
	ingressv1 "github.com/auth0/logger0/rpc/ingress/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewLogEndpointClient(ctx context.Context, u *url.URL, token string) (ingressv1.LogEndpointClient, error) {
	conn, err := dial(ctx, u, token)
	if err != nil {
		return nil, err
	}

	return ingressv1.NewLogEndpointClient(conn), nil
}

func NewSessionEndpointClient(ctx context.Context, u *url.URL, token string) (controlv1.SessionEndpointClient, error) {
	conn, err := dial(ctx, u, token)
	if err != nil {
		return nil, err
	}

	return controlv1.NewSessionEndpointClient(conn), nil
}

func NewDrainEndpointClient(ctx context.Context, u *url.URL, token string) (controlv1.DrainEndpointClient, error) {
	conn, err := dial(ctx, u, token)
	if err != nil {
		return nil, err
	}

	return controlv1.NewDrainEndpointClient(conn), nil
}

func dial(ctx context.Context, u *url.URL, token string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&auth.GRPCCredentials{Token: token}),
		grpc.WithBlock(),
	}

	if u.Scheme == "http" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	addr := u.Host
	if u.Port() == "" {
		switch u.Scheme {
		case "https":
			addr += ":443"
		case "http":
			addr += ":80"
		}
	}

	return grpc.DialContext(ctx, addr, opts...)
}
