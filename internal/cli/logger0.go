package cli

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/auth0-cli/internal/ansi"
	"github.com/auth0/logger0"
	"github.com/auth0/logger0/rpc/client"
	controlv1 "github.com/auth0/logger0/rpc/control/v1"
	"github.com/joeshaw/envdecode"
)

func logger0Stream(ctx context.Context, tenant, trigger, actionID string, actionMappings map[string]string) error {
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

	req := &controlv1.CreateSessionRequest{
		Type:   logger0.LogRecord_TYPE_ACTIONS,
		Tenant: translateTenant(tenant),

		// TODO(cyx): we should use constants exposed by logger0 here
		// instead, or an enum in the proto layer.
		Filters: map[string]string{
			"action_id":    actionID,
			"trigger_type": translateTrigger(trigger),
		},
	}

	stream, err := client.CreateSession(ctx, req)
	if err != nil {
		return err
	}

	for {
		rec, err := stream.Recv()
		if err != nil {
			return err
		}

		formatLogRecord(rec, actionMappings)
	}
}

func logger0DrainAdd(ctx context.Context, tenant string, target string, filters map[string]string) error {
	// NOTE(cyx): this is just a ship, the eventual state is we'll talk to
	// API2 anyway so it's fine to hack around this.
	var cfg struct {
		Token string   `env:"LOGGER0_CONTROL_TOKEN,required"`
		URL   *url.URL `env:"LOGGER0_CONTROL_URL,required"`
	}

	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}

	client, err := client.NewDrainEndpointClient(ctx, cfg.URL, cfg.Token)
	if err != nil {
		return err
	}

	req := &controlv1.CreateDrainRequest{
		Tenant: translateTenant(tenant),
		Sink: &logger0.Sink{
			Target: &logger0.Sink_Url{Url: target},
		},

		// TODO(cyx): we should use constants exposed by logger0 here
		// instead, or an enum in the proto layer.
		Filters: translateFilters(filters),
	}

	_, err = client.CreateDrain(ctx, req)
	return err
}

func logger0DrainDelete(ctx context.Context, tenant string, target string) error {
	// NOTE(cyx): this is just a ship, the eventual state is we'll talk to
	// API2 anyway so it's fine to hack around this.
	var cfg struct {
		Token string   `env:"LOGGER0_CONTROL_TOKEN,required"`
		URL   *url.URL `env:"LOGGER0_CONTROL_URL,required"`
	}

	if err := envdecode.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}

	client, err := client.NewDrainEndpointClient(ctx, cfg.URL, cfg.Token)
	if err != nil {
		return err
	}

	list, err := client.ListDrains(ctx, &controlv1.ListDrainsRequest{Tenant: translateTenant(tenant)})
	if err != nil {
		return err
	}

	var found *logger0.Drain
	for _, d := range list.Drains {
		switch v := d.GetSink().Target.(type) {
		case *logger0.Sink_Url:
			if v.Url == target {
				found = d
			}
		default:
			return fmt.Errorf("Unknown sink type: %T %v", d.GetSink(), d.GetSink())
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to find drain with target %s", target)
	}

	_, err = client.DeleteDrain(ctx, &controlv1.DeleteDrainRequest{
		Id: found.Id,
	})

	return err
}

func formatLogRecord(rec *logger0.LogRecord, actionMappings map[string]string) {
	msgs := bytes.Join(rec.Messages, []byte("\n"))

	kvs := []string{
		"trigger=" + unwrapTrigger(rec.Metadata["trigger_type"]),
		"action=" + actionMappings[rec.Metadata["action_id"]],
	}
	fmt.Printf("%s - %s - %s\n",
		ansi.Yellow(time.Now().Format(time.RFC3339)),
		ansi.Green(strings.Join(kvs, " ")),
		msgs,
	)
}

// NOTE(cyx): this is only necessary because we're punching directly to actions
// protocol. API2 in theory should do this for us (it already does this
// translation).
func translateTrigger(t string) string {
	return strings.Replace(strings.ToUpper(t), "-", "_", -1)
}

// NOTE(cyx): this is only necessary because we're punching directly to actions
// protocol. API2 in theory should do this for us (it already does this
// translation).
func unwrapTrigger(t string) string {
	return strings.Replace(strings.ToLower(t), "_", "-", -1)
}

// NOTE(cyx): should use a domain util helper
func translateTenant(t string) string {
	chunks := strings.Split(t, ".")
	return chunks[0]
}

func translateFilters(f map[string]string) []*logger0.Filter {
	var result []*logger0.Filter

	for k, v := range f {
		result = append(result, &logger0.Filter{
			Key: k,
			Val: v,
		})
	}

	return result
}
