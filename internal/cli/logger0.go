package cli

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/auth0/auth0-cli/internal/ansi"
	"github.com/auth0/logger0"
	"github.com/auth0/logger0/rpc/client"
	controlv1 "github.com/auth0/logger0/rpc/control/v1"
	"github.com/joeshaw/envdecode"
)

func logger0Stream(ctx context.Context, tenant, trigger, actionID string) error {
	// NOTE(cyx): this is just a ship, the eventual state is we'll talk to
	// API2 anyway so it's fine to hack around this.
	var cfg struct {
		Token string   `env:"LOGGER0_CONTROL_TOKEN,required"`
		URL   *url.URL `env:"LOGGER0_CONTROL_URL,required"`
	}

	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}

	client, err := client.NewSessionEndpointClient(ctx, cfg.URL, cfg.Token)
	if err != nil {
		return err
	}

	stream, err := client.CreateSession(ctx, &controlv1.CreateSessionRequest{
		Type:   logger0.LogRecord_TYPE_ACTIONS,
		Tenant: tenant,

		// TODO(cyx): we should use constants exposed by logger0 here
		// instead, or an enum in the proto layer.
		Filters: map[string]string{
			"action_id":    actionID,
			"trigger_type": trigger,
		},
	})
	if err != nil {
		return err
	}

	for {
		rec, err := stream.Recv()
		if err != nil {
			return err
		}

		formatLogRecord(rec)
	}
}

func formatLogRecord(rec *logger0.LogRecord) {
	msgs := bytes.Join(rec.Messages, []byte("\n"))
	fmt.Printf("%s - %s\n", ansi.Yellow(time.Now().Format(time.RFC3339)), msgs)
}
