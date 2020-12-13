# cw-agent

![Build Status](https://github.com/asmarques/cw-agent/workflows/CI/badge.svg)

An agent for reporting memory and disk metrics to AWS CloudWatch. Supports Linux and FreeBSD hosts.

## Usage

Report all available memory and disk metrics for the current EC2 instance every 5 minutes (the default):
```bash
cw-agent -all-metrics
```

Metrics can be selected individually:
```bash
cw-agent -mem-util -mem-avail -mem-used
```

You can also choose to report metrics once and exit (useful if running from cron):
```bash
cw-agent -all-metrics -once
```

## Installation

- [Download](https://github.com/asmarques/cw-agent/releases) a binary for the latest version
- Copy `cw-agent` to the filesystem of the host you wish to report metrics for
- Ensure the binary is executable: `chmod +x cw-agent`
- Run `cw-agent` as a process managed by your init system (e.g. systemd) or periodically from cron
 
## Credentials

If running from an EC2 instance with an associated IAM role, no credentials need to be supplied, otherwise
you can set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables to supply the necessary
credentials.

In either case, you must ensure the user/role is associated with a policy containing the appropriate
permissions: 

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:PutMetricData"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
} 
```

## License

[MIT](LICENSE)
