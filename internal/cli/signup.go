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

func signupCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signup",
		Args:  cobra.NoArgs,
		Short: "Signup for an Auth0 account",
		Long:  "If you dont alredy have an Auth0 account, you can signup for one using the CLI!",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			_, err := RunLogin(ctx, cli, false)
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

	return cmd
}

// RunSignup runs the signup flow guiding the user through the process
// by showing the signup instructions.
func RunSignup(ctx context.Context, cli *cli, expired bool) (tenant, error) {
	if expired {
		cli.renderer.Warnf("Please sign in to re-authorize the CLI.")
	} else {
		fmt.Print("✪ Welcome to the Auth0 CLI 🎊\n\n")
		fmt.Print("If you don't have an account, please go to https://auth0.com/signup\n\n")
	}

	state, err := cli.authenticator.Start(ctx)
	if err != nil {
		return tenant{}, fmt.Errorf("Could not start the authentication process: %w.", err)
	}

	fmt.Printf("Your Device Confirmation code is: %s\n\n", ansi.Bold(state.UserCode))

	if cli.noInput {
		cli.renderer.Infof("Open the following URL in a browser: %s\n", ansi.Green(state.VerificationURI))
	} else {
		cli.renderer.Infof("%s to open the browser to log in or %s to quit...", ansi.Green("Press Enter"), ansi.Red("^C"))
		fmt.Scanln()
		err = browser.OpenURL(state.VerificationURI)

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
