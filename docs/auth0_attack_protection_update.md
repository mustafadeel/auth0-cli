---
layout: default
---
## auth0 attack-protection update

Update attack protection settings

### Synopsis

Update attack protection settings.

```
auth0 attack-protection update [flags]
```

### Examples

```
auth0 attack-protection update
auth0 attack-protection update <attack-protection-type>
```

### Options

```
      --bfp-allowlist strings                      List of trusted IP addresses that will not have attack protection enforced against them. Comma-separated.
      --bfp-enabled                                Enable (or disable) brute force protection.
      --bfp-max-attempts int                       Maximum number of unsuccessful attempts.
      --bfp-mode string                            Account Lockout: Determines whether or not IP address is used when counting failed attempts. Possible values:
                                                   count_per_identifier_and_ip, count_per_identifier.
      --bfp-shields strings                        Action to take when a brute force protection threshold is violated. Possible values: block, user_notification.
                                                   Comma-separated.
      --bpd-admin-notification-frequency strings   When "admin_notification" is enabled, determines how often email notifications are sent. Possible values:
                                                   immediately, daily, weekly, monthly. Comma-separated.
      --bpd-enabled                                Enable (or disable) breached password detection.
      --bpd-method string                          The subscription level for breached password detection methods. Use "enhanced" to enable Credential Guard. Possible
                                                   values: standard, enhanced.
      --bpd-shields strings                        Action to take when a breached password is detected. Possible values: block, user_notification, admin_notification.
                                                   Comma-separated.
  -h, --help                                       help for update
      --sit-allowlist strings                      List of trusted IP addresses that will not have attack protection enforced against them. Comma-separated.
      --sit-enabled                                Enable (or disable) suspicious ip throttling.
      --sit-pre-login-max int                      Configuration options that apply before every login attempt. Total number of attempts allowed per day.
      --sit-pre-login-rate int                     Configuration options that apply before every login attempt. Interval of time, given in milliseconds, at which new
                                                   attempts are granted.
      --sit-pre-registration-max int               Configuration options that apply before every user registration attempt. Total number of attempts allowed.
      --sit-pre-registration-rate int              Configuration options that apply before every user registration attempt. Interval of time, given in milliseconds,
                                                   at which new attempts are granted.
      --sit-shields strings                        Action to take when a suspicious IP throttling threshold is violated. Possible values: block, admin_notification.
                                                   Comma-separated.
```

### Options inherited from parent commands

```
      --debug           Enable debug mode.
      --force           Skip confirmation.
      --format string   Command output format. Options: json.
      --no-color        Disable colors.
      --no-input        Disable interactivity.
      --tenant string   Specific tenant to use.
```

### SEE ALSO

* [auth0 attack-protection](auth0_attack_protection.md)	 - Manage resources for attack protection
