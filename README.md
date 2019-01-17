# Windows Syslog BOSH Release
* Slack: #syslog on <https://slack.cloudfoundry.org>
* Tracker: [CF Platform Logging Improvements][tracker]
* CI: [Syslog CI][CI]
* Syslog Release for Linux: [Release][syslogLinux]

This is a BOSH release for forwarding logs
from BOSH jobs on Windows VMs.
It forwards logs in
c:/var/vcap/sys/log/ (and any subdirectories, recursively)
to a configured syslog server.

## Differences From `syslog-release`
This release is designed to share configuration with its linux sibling.
However, it uses a very different architecture,
and parity would often be prohibitively expensive.
Here are a list of some of the major differences with `syslog-release`,
with special attention to standard configuration that will be ignored or fail.

- RELP is not a supported transport protocol. While `syslog-forwarder-windows` will start, it will not forward logs.
- Does not forward os event logs: use [Event Log Release][event-log-release]
- Does not support fallback servers (and will not use any if configured)
- Does not support custom rules (and will not respect them if configured)
- Does not support configuration of `permitted_peer` (what addresses or dns resolutions are allowed to recieve syslog messages) in tls mode
- Does not support custom configuration of various queuing parameters

If any of these constraints are a problem for you,
please open an issue explaining your use case.

## Configuring Log Forwarding
Add the `syslog_forwarder`
to forward all local syslog messages
from an instance
to a syslog endpoint.
You can use `addons` to add syslog forwarder to all instances;
if you are using `cf-deployment,`
there is an ops file to accomplish this [here](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/experimental/windows-enable-component-syslog.yml)
Configure `address` and,
optionally,
`port` and `transport`:

```yml
instance_groups:
- name: some-instance-group
  jobs:
  - name: syslog_forwarder_windows
    release: windows-syslog
  properties:
    syslog:
      address: <IP or hostname>
...
releases:
- name: "windows-syslog"
  version: "X.X.X"
  url: "https://bosh.io/d/github.com/cloudfoundry/windows-syslog-release?v=X.X.X"
  sha1: ""
```

You can get the windows-syslog-release from [bosh.io/releases](https://bosh.io/releases/github.com/cloudfoundry/windows-syslog-release?all=1).

If the syslog endpoint is unavailable,
messages will be briefly queued.

TLS over TCP is supported with additional properties.
`tls-enabled` should be set to true if you wish to use it.
In a future version, this will likely be true by default.
By default, the windows certificate API is used to validate certs.
If the cert you wish to respect isn't validated by the Windows API,
you will need to set the full cert chain with `ca_cert`.
Note that this fully replaces use of the Windows API.

## Testing and Debugging
Here are some tips for debugging issues with log forwarding.
We're assuming some familiarity with BOSH,
but not necessarily with Windows.

There is a useful manifest at `tests/manifests/tcp.yml`.
To use it, you will need to provide a deployment name, like so:
`bosh deploy -d windows-syslog -v deployment=windows-syslog tests/manifests/tcp.yml`.

To ssh on to a windows VM, you will need the windows utility release.
Our test manifest includes it.
You can then use `bosh -d windows-syslog ssh --opts=-T forwarder "echo example command"`
to run a single command non interactively on the forwarder job.
Interactive SSH works more or less like you'd expect,
except that you arrive in `cmd.exe`.

Our test manifest includes a storer from the linux syslog release,
and the windows forwarder job is configured to send logs to it.

## Maintainer's Note
The blobstore for this release
is on Google Cloud Storage.
Access is controlled by membership
in the cf-syslog@pivotal.io mailing list.

[tracker]: https://www.pivotaltracker.com/n/projects/2126318
[CI]: https://syslog.ci.cf-app.com
[syslogLinux]: https://github.com/cloudfoundry/syslog-release
[event-log-release]: https://github.com/cloudfoundry-incubator/event-log-release
