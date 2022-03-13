package cli

import (
//	"encoding/json"
	"fmt"
//	"strings"

	"github.com/auth0/auth0-cli/internal/ansi"
//	"github.com/auth0/auth0-cli/internal/auth0"
//	"github.com/auth0/auth0-cli/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/auth0/go-auth0/management"
)


func mfaCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mfa",
		Short: "manage mfa policies",
		Long:  "manage mfa policies.",
	}

	cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(showMultifactorCmd(cli))
	cmd.AddCommand(enableMultifactorCmd(cli))

	return cmd
}

func showMultifactorCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Args:    cobra.NoArgs,
		Short:   "show your mfa settings",
		Long: `Shows whether MFA is enabled for your applications`,
		Example: `auth0 mfa show`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var policy *management.MultiFactorPolicies

			if err := ansi.Waiting(func() error {
				var err error
				policy, err = cli.api.MultiFactor.Policy()
				return err
			}); err != nil {
				return fmt.Errorf("An unexpected error occurred: %w", err)
			}

			
			// TODO: write cli render for MFA policy read
			///fmt.Println("mfa enabled")
			//cli.renderer.MultiFactorPolicies(policy)
			return nil
		},
	}

	return cmd
}

func enableMultifactorCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable",
		Args:    cobra.NoArgs,
		Short:   "enable mfa for all-apps",
		Long: `Shows whether MFA is enabled for your applications`,
		Example: `auth0 mfa enable`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var policy *management.MultiFactorPolicies = &management.MultiFactorPolicies{"all-applications"}

			if err := ansi.Waiting(func() error {
				var err error
				err = cli.api.MultiFactor.UpdatePolicy(policy)
				return err
			}); err != nil {
				return fmt.Errorf("An unexpected error occurred: %w", err)
			}

			// TODO: write cli render for MFA policy enable
			//cli.renderer.MultiFactorPolicies(policy)
			return nil
		},
	}

	return cmd
}
