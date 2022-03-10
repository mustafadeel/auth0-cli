package cli

import (
	"context"
	"fmt"
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
		Use:   "signup",
		Args:  cobra.NoArgs,
		Short: "Get a new account and authenticate the Auth0 CLI",
		Long:  "Get your Auth0 account and authorize the CLI to access the Management API.",
		Example: `auth0 signup -email <email>
auth0 signup -e <email>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := signupEmail.Ask(cmd, &inputs.Email, nil); err != nil {
				return err
			}
			ctx := cmd.Context()
			_, err := RunSignup(ctx, cli, inputs.Email)
			if err == nil {
				cli.tracker.TrackCommandRun(cmd, cli.config.InstallID)
			}
			fmt.Print("error occurred\n\n")
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
		browserURL := fmt.Sprintf(state.VerificationURI + "&loginhint=" + email)
		err = browser.OpenURL(browserURL)

		if err != nil {
			cli.renderer.Warnf("Couldn't open the URL, please do it manually: %s.", state.VerificationURI)
		}
	}

	var res auth.Result
	err = ansi.Spinner("Waiting for login to complete in browser", func() error {
		res, err = cli.authenticator.Wait(ctx, state)
		return err
	})

	if err != nil {
		return tenant{}, fmt.Errorf("login error: %w", err)
	}

	fmt.Print("\n")
	cli.renderer.Infof("Successfully logged in.")
	cli.renderer.Infof("Tenant: %s\n", res.Domain)
	cli.renderer.Infof("Access Token: %s\n", res.AccessToken)
	cli.renderer.Infof("Refresh Token: %s\n", res.RefreshToken)

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
