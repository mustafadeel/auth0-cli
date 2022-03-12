package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/auth0/auth0-cli/internal/prompt"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/spf13/cobra"
)

var (
	tenantDomain = Argument{
		Name: "Tenant",
		Help: "Tenant to select",
	}

	tenantClientID = Flag{
		Name:       "Client ID",
		LongForm:   "client-id",
		ShortForm:  "i",
		Help:       "Client ID of the application.",
		IsRequired: true,
	}

	tenantClientSecret = Flag{
		Name:       "Client Secret",
		LongForm:   "client-secret",
		ShortForm:  "s",
		Help:       "Client Secret of the application.",
		IsRequired: true,
	}

	tenantName = Flag{
		Name:       "Tenant Name",
		LongForm:   "name",
		ShortForm:  "n",
		Help:       "Name of the application.",
		IsRequired: true,
	}

	tenantEnvironment = Flag{
		Name:       "Environment",
		LongForm:   "environment",
		ShortForm:  "e",
		Help:       "Environment in which to create tenant. <development|production>",
		IsRequired: true,
	}
)

func tenantsCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenants",
		Short: "Manage configured tenants",
		Long:  "Manage configured tenants.",
	}

	cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(useTenantCmd(cli))
	cmd.AddCommand(listTenantCmd(cli))
	cmd.AddCommand(openTenantCmd(cli))
	cmd.AddCommand(addTenantCmd(cli))
	cmd.AddCommand(createTenantCmd(cli))
	return cmd
}

func listTenantCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List your tenants",
		Long:    "List your tenants.",
		Example: "auth0 tenants list",
		RunE: func(cmd *cobra.Command, args []string) error {
			tens, err := cli.listTenants()
			if err != nil {
				return fmt.Errorf("Unable to load tenants due to an unexpected error: %w", err)
			}

			tenNames := make([]string, len(tens))
			for i, t := range tens {
				tenNames[i] = t.Domain
			}

			cli.renderer.TenantList(tenNames)
			return nil
		},
	}
	return cmd
}

func useTenantCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Set the active tenant",
		Long:    "Set the active tenant.",
		Example: "auth0 tenants use <tenant>",
		RunE: func(cmd *cobra.Command, args []string) error {
			var selectedTenant string
			if len(args) == 0 {
				tens, err := cli.listTenants()
				if err != nil {
					return fmt.Errorf("Unable to load tenants due to an unexpected error: %w", err)
				}

				tenNames := make([]string, len(tens))
				for i, t := range tens {
					tenNames[i] = t.Domain
				}

				input := prompt.SelectInput("tenant", "Tenant:", "Tenant to activate", tenNames, tenNames[0], true)
				if err := prompt.AskOne(input, &selectedTenant); err != nil {
					return handleInputError(err)
				}
			} else {
				requestedTenant := args[0]
				t, ok := cli.config.Tenants[requestedTenant]
				if !ok {
					return fmt.Errorf("Unable to find tenant %s; run 'auth0 tenants use' to see your configured tenants or run 'auth0 login' to configure a new tenant", requestedTenant)
				}
				selectedTenant = t.Domain
			}

			cli.config.DefaultTenant = selectedTenant
			if err := cli.persistConfig(); err != nil {
				return fmt.Errorf("An error occurred while setting the default tenant: %w", err)
			}
			cli.renderer.Infof("Default tenant switched to: %s", selectedTenant)
			return nil
		},
	}

	return cmd
}

func openTenantCmd(cli *cli) *cobra.Command {
	var inputs struct {
		Domain string
	}

	cmd := &cobra.Command{
		Use:     "open",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Open tenant settings page in the Auth0 Dashboard",
		Long:    "Open tenant settings page in the Auth0 Dashboard.",
		Example: "auth0 tenants open <tenant>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := tenantDomain.Pick(cmd, &inputs.Domain, cli.tenantPickerOptions)
				if err != nil {
					return err
				}
			} else {
				inputs.Domain = args[0]

				if _, ok := cli.config.Tenants[inputs.Domain]; !ok {
					return fmt.Errorf("Unable to find tenant %s; run 'auth0 login' to configure a new tenant", inputs.Domain)
				}
			}

			openManageURL(cli, inputs.Domain, "tenant/general")
			return nil
		},
	}

	return cmd
}

func addTenantCmd(cli *cli) *cobra.Command {
	var inputs struct {
		Domain       string
		ClientID     string
		ClientSecret string
	}

	cmd := &cobra.Command{
		Use:     "add",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Add a tenant with client credentials",
		Long:    "Add a tenant with client credentials.",
		Example: "auth0 tenants add <tenant> --client-id <id> --client-secret <secret>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := tenantDomain.Pick(cmd, &inputs.Domain, cli.tenantPickerOptions)
				if err != nil {
					if !errors.Is(err, errUnauthenticated) {
						return err
					}

					if err := tenantDomain.Ask(cmd, &inputs.Domain); err != nil {
						return err
					}
				}
			} else {
				inputs.Domain = args[0]
			}

			if err := tenantClientID.Ask(cmd, &inputs.ClientID, nil); err != nil {
				return err
			}

			if err := tenantClientSecret.Ask(cmd, &inputs.ClientSecret, nil); err != nil {
				return err
			}

			t := tenant{
				Domain:       inputs.Domain,
				ClientID:     inputs.ClientID,
				ClientSecret: inputs.ClientSecret,
			}

			if err := cli.addTenant(t); err != nil {
				return err
			}

			cli.renderer.Infof("Tenant added successfully: %s", t.Domain)
			return nil
		},
	}

	return cmd
}

func createTenantCmd(cli *cli) *cobra.Command {
	var inputs struct {
		Name        string
		Environment string
	}

	cmd := &cobra.Command{
		Use:     "create",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Create a new tenant",
		Long:    "Create a new tenant.",
		Example: "auth0 tenants create --name <name> --environment <development|production>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := tenantName.Ask(cmd, &inputs.Name, nil); err != nil {
				return err
			}

			if err := tenantEnvironment.Ask(cmd, &inputs.Environment, nil); err != nil {
				return err
			}

			// Parse the access token for the tenant.
			t, err := jwt.ParseString(cli.config.Tenants[cli.tenant].AccessToken)
			if err != nil {
				return err
			}

			emailClaim := ""
			// get custom email claim
			if v, ok := t.Get(`https://auth0cli.api/email`); ok {
				emailClaim = fmt.Sprint(v)
			}

			tenantData := Payload{
				Name:        inputs.Name,
				Environment: inputs.Environment,
				Owner: Owner{
					ID:    fmt.Sprint(t.Subject()),
					Name:  "ADMIN",
					Email: emailClaim,
				},
			}

			tenantResult := createTenant(tenantData)

			if tenantResult.Error != "" {
				cli.renderer.Infof("Unable to add tenant: %s", tenantResult.Message)
			} else {
				//add tenant to cli
				ten := tenant{
					Domain:       tenantResult.Name,
					ClientID:     tenantResult.ClientId,
					ClientSecret: tenantResult.ClientSecret,
				}
				cli.addTenant(ten)

				cli.renderer.Infof("Tenant added successfully 🎉")
				cli.renderer.Infof("  Name         : %s", tenantResult.Name)
				cli.renderer.Infof("  Client Id    : %s", tenantResult.ClientId)
				cli.renderer.Infof("  Client Secret: %s", tenantResult.ClientSecret)
			}
			return nil
		},
	}

	return cmd
}

func (c *cli) tenantPickerOptions() (pickerOptions, error) {
	tens, err := c.listTenants()
	if err != nil {
		return nil, fmt.Errorf("Unable to load tenants due to an unexpected error: %w", err)
	}

	var priorityOpts, opts pickerOptions

	for _, t := range tens {
		opt := pickerOption{value: t.Domain, label: t.Domain}

		// check if this is currently the default tenant.
		if t.Domain == c.config.DefaultTenant {
			priorityOpts = append(priorityOpts, opt)
		} else {
			opts = append(opts, opt)
		}
	}

	if len(opts)+len(priorityOpts) == 0 {
		return nil, errNoApps
	}

	return append(priorityOpts, opts...), nil
}

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl --request POST --url http://localhost:3000/tenant \
// --data '{"name":"test-create-7","environment":"development", "owner":{ "id":"auth0|6229066879f160002a42e5a9", "name": "ADMIN", "email": "root@auth0.com" }}' \
// --header 'content-type: application/json'

func createTenant(p Payload) Tenant {

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:3000/tenant", body)
	if err != nil {
		fmt.Println("NewRequest Error:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("http.DefaultClient.Do Error:", err)
	}
	defer resp.Body.Close()

	var t Tenant

	decodeTenantErr := json.NewDecoder(resp.Body).Decode(&t)
	if decodeTenantErr != nil {
		fmt.Println("Error getting tenant")
	}
	return t
}

type Payload struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Owner       Owner  `json:"owner"`
}
type Owner struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Tenant struct {
	Name         string `json:"name"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	StatusCode   int    `json:statusCode`
	Error        string `json:error`
	ErrorCode    string `json:errorCode`
	Message      string `json:message`
}
