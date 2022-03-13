package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/auth0/auth0-cli/internal/ansi"
	"github.com/auth0/auth0-cli/internal/auth"
	"github.com/auth0/auth0-cli/internal/prompt"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	signupEmail = Flag{
		Name:       "Email",
		LongForm:   "email",
		ShortForm:  "e",
		Help:       "Email to use for signup.",
		IsRequired: true,
	}
)

func signupCmd(cli *cli) *cobra.Command {
	var inputs struct {
		Email string
	}

	cmd := &cobra.Command{
		Use:     "signup",
		Args:    cobra.NoArgs,
		Short:   "Get a new account and authenticate the Auth0 CLI",
		Long:    "Get your Auth0 account and authorize the CLI to access the Management API.",
		Example: `auth0 signup`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			_, err := RunSignup(ctx, cli, inputs.Email)
			if err == nil {
				cli.tracker.TrackCommandRun(cmd, cli.config.InstallID)
			}
			return err
		},
	}

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Flags().MarkHidden("tenant")
		cmd.Parent().HelpFunc()(cmd, args)
	})

	signupEmail.RegisterString(cmd, &inputs.Email, "")
	return cmd
}

// RunLogin runs the login flow guiding the user through the process
// by showing the login instructions, opening the browser.
// Use `expired` to run the login from other commands setup:
// this will only affect the messages.
func RunSignup(ctx context.Context, cli *cli, email string) (tenant, error) {
	fmt.Print("✪ Welcome to the Auth0 CLI 🎊\n\n")

	state, err := cli.authenticator.Start(ctx, true)
	if err != nil {
		return tenant{}, fmt.Errorf("Could not start the signup process: %w", err)
	}

	fmt.Printf("Your Device Confirmation code is: %s\n\n", ansi.Bold(state.UserCode))

	if cli.noInput {
		cli.renderer.Infof("Open the following URL in a browser: %s\n", ansi.Green(state.VerificationURI))
	} else {
		cli.renderer.Infof("%s to open the browser to sign up or %s to quit...", ansi.Green("Press Enter"), ansi.Red("^C"))
		fmt.Scanln()
		err = browser.OpenURL(state.VerificationURI)

		if err != nil {
			cli.renderer.Warnf("Couldn't open the URL, please do it manually: %s.", state.VerificationURI)
		}
	}

	var res auth.Result
	err = ansi.Spinner("Waiting for signup to complete in browser", func() error {
		res, err = cli.authenticator.Wait(ctx, state)
		return err
	})

	if err != nil {
		return tenant{}, fmt.Errorf("signup error: %w", err)
	}

	fmt.Print("\n")
	cli.renderer.Infof("Successfully signed up!")
	if strings.Contains(authCfg.DeviceCodeEndpoint, res.Domain) {
		cli.renderer.Infof("Time to create your first tenant!\n\nUse `auth0 tenants create` to create your first tenant.")
	}

	// store the refresh token
	secretsStore := &auth.Keyring{}
	err = secretsStore.Set(auth.SecretsNamespace, res.Domain, res.RefreshToken)
	if err != nil {
		// log the error but move on
		cli.renderer.Warnf("Could not store the refresh token locally, please expect to login again once your access token expired. See https://github.com/auth0/auth0-cli/blob/main/KNOWN-ISSUES.md.")
	}

	t := tenant{
		Name:        res.Tenant,
		Domain:      res.Domain,
		AccessToken: res.AccessToken,
		ExpiresAt: time.Now().Add(
			time.Duration(res.ExpiresIn) * time.Second,
		),
		Scopes: auth.RequiredScopes(),
	}
	err = cli.addTenant(t)
	if err != nil {
		return tenant{}, fmt.Errorf("Could not add tenant to config: %w", err)
	}

	if err := checkInstallID(cli); err != nil {
		return tenant{}, fmt.Errorf("Could not update config: %w", err)
	}

	if cli.config.DefaultTenant != res.Domain {
		promptText := fmt.Sprintf("Your default tenant is %s. Do you want to change it to %s?", cli.config.DefaultTenant, res.Domain)
		if confirmed := prompt.Confirm(promptText); !confirmed {
			return tenant{}, nil
		}
		cli.config.DefaultTenant = res.Domain
		if err := cli.persistConfig(); err != nil {
			cli.renderer.Warnf("Could not set the default tenant, please try 'auth0 tenants use %s': %w", res.Domain, err)
		}
	}

	return t, nil
}
