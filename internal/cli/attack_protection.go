package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/auth0/go-auth0/management"
	"github.com/spf13/cobra"

	"github.com/auth0/auth0-cli/internal/ansi"
)

const (
	breachedPasswordDetection = "breached-password-detection"
	bruteForceProtection      = "brute-force-protection"
	suspiciousIPThrottling    = "suspicious-ip-throttling"
)

var (
	attackProtectionType = Argument{
		Name: "Attack Protection Type",
		Help: "Specific attack protection type to manage.",
	}

	apFlags = attackProtectionFlags{
		bpd: breachedPasswordDetectionFlags{
			Enabled: Flag{
				Name:         "Enabled",
				LongForm:     "bpd-enabled",
				Help:         "Enable (or disable) breached password detection.",
				AlwaysPrompt: true,
			},
			Shields: Flag{
				Name:         "Shields",
				LongForm:     "bpd-shields",
				Help:         "Action to take when a breached password is detected. Possible values: block, user_notification, admin_notification. Comma-separated.",
				AlwaysPrompt: true,
			},
			AdminNotificationFrequency: Flag{
				Name:     "Admin Notification Frequency",
				LongForm: "bpd-admin-notification-frequency",
				Help: "When \"admin_notification\" is enabled, determines how often email notifications " +
					"are sent. Possible values: immediately, daily, weekly, monthly. Comma-separated.",
				AlwaysPrompt: true,
			},
			Method: Flag{
				Name:         "Method",
				LongForm:     "bpd-method",
				Help:         "The subscription level for breached password detection methods. Use \"enhanced\" to enable Credential Guard. Possible values: standard, enhanced.",
				AlwaysPrompt: true,
			},
		},
		bfp: bruteForceProtectionFlags{
			Enabled: Flag{
				Name:         "Enabled",
				LongForm:     "bfp-enabled",
				Help:         "Enable (or disable) brute force protection.",
				AlwaysPrompt: true,
			},
			Shields: Flag{
				Name:     "Shields",
				LongForm: "bfp-shields",
				Help: "Action to take when a brute force protection threshold is violated." +
					" Possible values: block, user_notification. Comma-separated.",
				AlwaysPrompt: true,
			},
			AllowList: Flag{
				Name:     "Allow List",
				LongForm: "bfp-allowlist",
				Help: "List of trusted IP addresses that will not have " +
					"attack protection enforced against them. Comma-separated.",
				AlwaysPrompt: true,
			},
			Mode: Flag{
				Name:     "Mode",
				LongForm: "bfp-mode",
				Help: "Account Lockout: Determines whether or not IP address is used when counting " +
					"failed attempts. Possible values: count_per_identifier_and_ip, count_per_identifier.",
				AlwaysPrompt: true,
			},
			MaxAttempts: Flag{
				Name:         "MaxAttempts",
				LongForm:     "bfp-max-attempts",
				Help:         "Maximum number of unsuccessful attempts.",
				AlwaysPrompt: true,
			},
		},
		sit: suspiciousIPThrottlingFlags{
			Enabled: Flag{
				Name:         "Enabled",
				LongForm:     "sit-enabled",
				Help:         "Enable (or disable) suspicious ip throttling.",
				AlwaysPrompt: true,
			},
			Shields: Flag{
				Name:     "Shields",
				LongForm: "sit-shields",
				Help: "Action to take when a suspicious IP throttling threshold is violated. " +
					"Possible values: block, admin_notification. Comma-separated.",
				AlwaysPrompt: true,
			},
			AllowList: Flag{
				Name:     "Allow List",
				LongForm: "sit-allowlist",
				Help: "List of trusted IP addresses that will not have attack protection enforced against " +
					"them. Comma-separated.",
				AlwaysPrompt: true,
			},
			StagePreLoginMaxAttempts: Flag{
				Name:     "StagePreLoginMaxAttempts",
				LongForm: "sit-pre-login-max",
				Help: "Configuration options that apply before every login attempt. " +
					"Total number of attempts allowed per day.",
				AlwaysPrompt: true,
			},
			StagePreLoginRate: Flag{
				Name:     "StagePreLoginRate",
				LongForm: "sit-pre-login-rate",
				Help: "Configuration options that apply before every login attempt. " +
					"Interval of time, given in milliseconds, at which new attempts are granted.",
				AlwaysPrompt: true,
			},
			StagePreUserRegistrationMaxAttempts: Flag{
				Name:     "StagePreUserRegistrationMaxAttempts",
				LongForm: "sit-pre-registration-max",
				Help: "Configuration options that apply before every user registration attempt. " +
					"Total number of attempts allowed.",
				AlwaysPrompt: true,
			},
			StagePreUserRegistrationRate: Flag{
				Name:     "StagePreUserRegistrationRate",
				LongForm: "sit-pre-registration-rate",
				Help: "Configuration options that apply before every user registration attempt. " +
					"Interval of time, given in milliseconds, at which new attempts are granted.",
				AlwaysPrompt: true,
			},
		},
	}
)

type (
	attackProtectionFlags struct {
		bpd breachedPasswordDetectionFlags
		bfp bruteForceProtectionFlags
		sit suspiciousIPThrottlingFlags
	}

	breachedPasswordDetectionFlags struct {
		Enabled                    Flag
		Shields                    Flag
		AdminNotificationFrequency Flag
		Method                     Flag
	}

	bruteForceProtectionFlags struct {
		Enabled     Flag
		Shields     Flag
		AllowList   Flag
		Mode        Flag
		MaxAttempts Flag
	}

	suspiciousIPThrottlingFlags struct {
		Enabled                             Flag
		Shields                             Flag
		AllowList                           Flag
		StagePreLoginMaxAttempts            Flag
		StagePreLoginRate                   Flag
		StagePreUserRegistrationMaxAttempts Flag
		StagePreUserRegistrationRate        Flag
	}

	attackProtectionInputs struct {
		bpd breachedPasswordDetectionInputs
		bfp bruteForceProtectionInputs
		sit suspiciousIPThrottlingInputs
	}

	breachedPasswordDetectionInputs struct {
		Enabled                    bool
		Shields                    []string
		AdminNotificationFrequency []string
		Method                     string
	}

	bruteForceProtectionInputs struct {
		Enabled     bool
		Shields     []string
		AllowList   []string
		Mode        string
		MaxAttempts int
	}

	suspiciousIPThrottlingInputs struct {
		Enabled                             bool
		Shields                             []string
		AllowList                           []string
		StagePreLoginMaxAttempts            int
		StagePreLoginRate                   int
		StagePreUserRegistrationMaxAttempts int
		StagePreUserRegistrationRate        int
	}
)

func attackProtectionCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "protection",
		Aliases: []string{"attack-protection", "ap"},
		Short:   "Manage resources for attack protection",
		Long:    "Manage resources for attack protection.",
	}

	cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(showAttackProtectionCmd(cli))
	cmd.AddCommand(updateAttackProtectionCmd(cli))

	return cmd
}

func showAttackProtectionCmd(cli *cli) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Show attack protection settings",
		Long:    "Show attack protection settings.",
		Example: `auth0 protection show`,
		RunE:    showAttackProtectionCmdRun(cli),
	}
}

func updateAttackProtectionCmd(cli *cli) *cobra.Command {
	var inputs attackProtectionInputs

	cmd := &cobra.Command{
		Use:     "update",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Update attack protection settings",
		Long:    "Update attack protection settings.",
		Example: `auth0 protection update`,
		RunE:    updateAttackProtectionCmdRun(cli, &inputs),
	}

	apFlags.bpd.Enabled.RegisterBoolU(cmd, &inputs.bpd.Enabled, false)
	apFlags.bpd.Shields.RegisterStringSliceU(cmd, &inputs.bpd.Shields, []string{})
	apFlags.bpd.AdminNotificationFrequency.RegisterStringSliceU(cmd, &inputs.bpd.AdminNotificationFrequency, []string{})
	apFlags.bpd.Method.RegisterString(cmd, &inputs.bpd.Method, "")

	apFlags.bfp.Enabled.RegisterBoolU(cmd, &inputs.bfp.Enabled, false)
	apFlags.bfp.Shields.RegisterStringSliceU(cmd, &inputs.bfp.Shields, []string{})
	apFlags.bfp.AllowList.RegisterStringSliceU(cmd, &inputs.bfp.AllowList, []string{})
	apFlags.bfp.Mode.RegisterString(cmd, &inputs.bfp.Mode, "")
	apFlags.bfp.MaxAttempts.RegisterIntU(cmd, &inputs.bfp.MaxAttempts, 1)

	apFlags.sit.Enabled.RegisterBoolU(cmd, &inputs.sit.Enabled, false)
	apFlags.sit.Shields.RegisterStringSliceU(cmd, &inputs.sit.Shields, []string{})
	apFlags.sit.AllowList.RegisterStringSliceU(cmd, &inputs.sit.AllowList, []string{})
	apFlags.sit.StagePreLoginMaxAttempts.RegisterIntU(cmd, &inputs.sit.StagePreLoginMaxAttempts, 1)
	apFlags.sit.StagePreLoginRate.RegisterIntU(cmd, &inputs.sit.StagePreLoginRate, 34560)
	apFlags.sit.StagePreUserRegistrationMaxAttempts.RegisterIntU(cmd, &inputs.sit.StagePreUserRegistrationMaxAttempts, 1)
	apFlags.sit.StagePreUserRegistrationRate.RegisterIntU(cmd, &inputs.sit.StagePreUserRegistrationRate, 1200)

	return cmd
}

func showAttackProtectionCmdRun(cli *cli) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		attackProtectionType, err := pickedAttackProtectionType(cmd, args)
		if err != nil {
			return err
		}

		switch attackProtectionType {
		case breachedPasswordDetection:
			if err := displayBreachedPasswordDetection(cli); err != nil {
				return err
			}
		case bruteForceProtection:
			if err := displayBruteForceProtection(cli); err != nil {
				return err
			}
		case suspiciousIPThrottling:
			if err := displaySuspiciousIPThrottling(cli); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid attack protection type: %s", attackProtectionType)
		}

		return nil
	}
}

func updateAttackProtectionCmdRun(
	cli *cli,
	inputs *attackProtectionInputs,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		attackProtectionType, err := pickedAttackProtectionType(cmd, args)
		if err != nil {
			return err
		}

		switch attackProtectionType {
		case breachedPasswordDetection:
			if err := updateAndDisplayBreachedPasswordDetection(cli, cmd, &inputs.bpd); err != nil {
				return err
			}
		case bruteForceProtection:
			if err := updateAndDisplayBruteForceProtection(cli, cmd, &inputs.bfp); err != nil {
				return err
			}
		case suspiciousIPThrottling:
			if err := updateAndDisplaySuspiciousIPThrottling(cli, cmd, &inputs.sit); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid attack protection type: %s", attackProtectionType)
		}

		return nil
	}
}

func pickedAttackProtectionType(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	var inputType string
	err := attackProtectionType.Pick(cmd, &inputType, attackProtectionTypePickerOptions)

	return inputType, err
}

func attackProtectionTypePickerOptions() (pickerOptions, error) {
	return pickerOptions{
		pickerOption{
			label: "Breached Password Detection",
			value: breachedPasswordDetection,
		},
		pickerOption{
			label: "Brute Force Protection",
			value: bruteForceProtection,
		},
		pickerOption{
			label: "Suspicious IP Throttling",
			value: suspiciousIPThrottling,
		},
	}, nil
}

func displayBreachedPasswordDetection(cli *cli) error {
	var bpd *management.BreachedPasswordDetection
	err := ansi.Waiting(func() (err error) {
		bpd, err = cli.api.AttackProtection.GetBreachedPasswordDetection()
		return err
	})
	if err != nil {
		return err
	}

	cli.renderer.BreachedPasswordDetectionShow(bpd)

	return nil
}

func displayBruteForceProtection(cli *cli) error {
	var bfp *management.BruteForceProtection
	err := ansi.Waiting(func() (err error) {
		bfp, err = cli.api.AttackProtection.GetBruteForceProtection()
		return err
	})
	if err != nil {
		return err
	}

	cli.renderer.BruteForceProtectionShow(bfp)

	return nil
}

func displaySuspiciousIPThrottling(cli *cli) error {
	var sit *management.SuspiciousIPThrottling
	err := ansi.Waiting(func() (err error) {
		sit, err = cli.api.AttackProtection.GetSuspiciousIPThrottling()
		return err
	})
	if err != nil {
		return err
	}

	cli.renderer.SuspiciousIPThrottlingShow(sit)

	return nil
}

func updateAndDisplayBreachedPasswordDetection(
	cli *cli,
	cmd *cobra.Command,
	inputs *breachedPasswordDetectionInputs,
) error {
	var bpd *management.BreachedPasswordDetection
	err := ansi.Waiting(func() (err error) {
		bpd, err = cli.api.AttackProtection.GetBreachedPasswordDetection()
		return err
	})
	if err != nil {
		return err
	}

	if err := apFlags.bpd.Enabled.AskBoolU(cmd, &inputs.Enabled, bpd.Enabled); err != nil {
		return err
	}
	bpd.Enabled = &inputs.Enabled

	shieldsString := strings.Join(bpd.GetShields(), ",")
	if err := apFlags.bpd.Shields.AskManyU(cmd, &inputs.Shields, &shieldsString); err != nil {
		return err
	}
	if len(inputs.Shields) == 0 {
		inputs.Shields = bpd.GetShields()
	}
	bpd.Shields = &inputs.Shields

	adminNotificationFrequencyString := strings.Join(bpd.GetAdminNotificationFrequency(), ",")
	if err := apFlags.bpd.AdminNotificationFrequency.AskManyU(
		cmd,
		&inputs.AdminNotificationFrequency,
		&adminNotificationFrequencyString,
	); err != nil {
		return err
	}
	if len(inputs.AdminNotificationFrequency) == 0 {
		inputs.AdminNotificationFrequency = bpd.GetAdminNotificationFrequency()
	}
	bpd.AdminNotificationFrequency = &inputs.AdminNotificationFrequency

	if err := apFlags.bpd.Method.AskU(cmd, &inputs.Method, bpd.Method); err != nil {
		return err
	}
	if inputs.Method == "" {
		inputs.Method = bpd.GetMethod()
	}
	bpd.Method = &inputs.Method

	if err := ansi.Waiting(func() error {
		return cli.api.AttackProtection.UpdateBreachedPasswordDetection(bpd)
	}); err != nil {
		return err
	}

	cli.renderer.BreachedPasswordDetectionUpdate(bpd)

	return nil
}

func updateAndDisplayBruteForceProtection(
	cli *cli,
	cmd *cobra.Command,
	inputs *bruteForceProtectionInputs,
) error {
	var bfp *management.BruteForceProtection
	err := ansi.Waiting(func() (err error) {
		bfp, err = cli.api.AttackProtection.GetBruteForceProtection()
		return err
	})
	if err != nil {
		return err
	}

	if err := apFlags.bfp.Enabled.AskBoolU(cmd, &inputs.Enabled, bfp.Enabled); err != nil {
		return err
	}
	bfp.Enabled = &inputs.Enabled

	shieldsString := strings.Join(bfp.GetShields(), ",")
	if err := apFlags.bfp.Shields.AskManyU(cmd, &inputs.Shields, &shieldsString); err != nil {
		return err
	}
	if len(inputs.Shields) == 0 {
		inputs.Shields = bfp.GetShields()
	}
	bfp.Shields = &inputs.Shields

	allowListString := strings.Join(bfp.GetAllowList(), ",")
	if err := apFlags.bfp.AllowList.AskManyU(
		cmd,
		&inputs.AllowList,
		&allowListString,
	); err != nil {
		return err
	}
	if len(inputs.AllowList) == 0 {
		inputs.AllowList = bfp.GetAllowList()
	}
	bfp.AllowList = &inputs.AllowList

	if err := apFlags.bfp.Mode.AskU(cmd, &inputs.Mode, bfp.Mode); err != nil {
		return err
	}
	if inputs.Mode == "" {
		inputs.Mode = bfp.GetMode()
	}
	bfp.Mode = &inputs.Mode

	defaultMaxAttempts := strconv.Itoa(bfp.GetMaxAttempts())
	if err := apFlags.bfp.MaxAttempts.AskIntU(cmd, &inputs.MaxAttempts, &defaultMaxAttempts); err != nil {
		return err
	}
	if inputs.MaxAttempts == 0 {
		inputs.MaxAttempts = bfp.GetMaxAttempts()
	}
	bfp.MaxAttempts = &inputs.MaxAttempts

	if err := ansi.Waiting(func() error {
		return cli.api.AttackProtection.UpdateBruteForceProtection(bfp)
	}); err != nil {
		return err
	}

	cli.renderer.BruteForceProtectionUpdate(bfp)

	return nil
}

func updateAndDisplaySuspiciousIPThrottling(
	cli *cli,
	cmd *cobra.Command,
	inputs *suspiciousIPThrottlingInputs,
) error {
	var sit *management.SuspiciousIPThrottling
	err := ansi.Waiting(func() (err error) {
		sit, err = cli.api.AttackProtection.GetSuspiciousIPThrottling()
		return err
	})
	if err != nil {
		return err
	}

	if err := apFlags.sit.Enabled.AskBoolU(cmd, &inputs.Enabled, sit.Enabled); err != nil {
		return err
	}
	sit.Enabled = &inputs.Enabled

	shieldsString := strings.Join(sit.GetShields(), ",")
	if err := apFlags.sit.Shields.AskManyU(cmd, &inputs.Shields, &shieldsString); err != nil {
		return err
	}
	if len(inputs.Shields) == 0 {
		inputs.Shields = sit.GetShields()
	}
	sit.Shields = &inputs.Shields

	allowListString := strings.Join(sit.GetAllowList(), ",")
	if err := apFlags.bfp.AllowList.AskManyU(
		cmd,
		&inputs.AllowList,
		&allowListString,
	); err != nil {
		return err
	}
	if len(inputs.AllowList) == 0 {
		inputs.AllowList = sit.GetAllowList()
	}
	sit.AllowList = &inputs.AllowList

	defaultPreLoginMaxAttempts := strconv.Itoa(sit.Stage.PreLogin.GetMaxAttempts())
	if err := apFlags.sit.StagePreLoginMaxAttempts.AskIntU(
		cmd,
		&inputs.StagePreLoginMaxAttempts,
		&defaultPreLoginMaxAttempts,
	); err != nil {
		return err
	}
	if inputs.StagePreLoginMaxAttempts == 0 {
		inputs.StagePreLoginMaxAttempts = sit.Stage.PreLogin.GetMaxAttempts()
	}
	sit.Stage.PreLogin.MaxAttempts = &inputs.StagePreLoginMaxAttempts

	defaultPreLoginRate := strconv.Itoa(sit.Stage.PreLogin.GetRate())
	if err := apFlags.sit.StagePreLoginRate.AskIntU(cmd, &inputs.StagePreLoginRate, &defaultPreLoginRate); err != nil {
		return err
	}
	if inputs.StagePreLoginRate == 0 {
		inputs.StagePreLoginRate = sit.Stage.PreLogin.GetRate()
	}
	sit.Stage.PreLogin.Rate = &inputs.StagePreLoginRate

	defaultPreUserRegistrationMaxAttempts := strconv.Itoa(sit.Stage.PreUserRegistration.GetMaxAttempts())
	if err := apFlags.sit.StagePreUserRegistrationMaxAttempts.AskIntU(
		cmd,
		&inputs.StagePreUserRegistrationMaxAttempts,
		&defaultPreUserRegistrationMaxAttempts,
	); err != nil {
		return err
	}
	if inputs.StagePreUserRegistrationMaxAttempts == 0 {
		inputs.StagePreUserRegistrationMaxAttempts = sit.Stage.PreUserRegistration.GetMaxAttempts()
	}
	sit.Stage.PreUserRegistration.MaxAttempts = &inputs.StagePreUserRegistrationMaxAttempts

	defaultPreUserRegistrationRate := strconv.Itoa(sit.Stage.PreUserRegistration.GetRate())
	if err := apFlags.sit.StagePreUserRegistrationRate.AskIntU(
		cmd,
		&inputs.StagePreUserRegistrationRate,
		&defaultPreUserRegistrationRate,
	); err != nil {
		return err
	}
	if inputs.StagePreUserRegistrationRate == 0 {
		inputs.StagePreUserRegistrationRate = sit.Stage.PreUserRegistration.GetRate()
	}
	sit.Stage.PreUserRegistration.Rate = &inputs.StagePreUserRegistrationRate

	if err := ansi.Waiting(func() error {
		return cli.api.AttackProtection.UpdateSuspiciousIPThrottling(sit)
	}); err != nil {
		return err
	}

	cli.renderer.SuspiciousIPThrottlingUpdate(sit)

	return nil
}
