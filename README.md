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

## Evaluation
If you wish to interact with this release to evaluate its behavior,
here are some tips.
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

[tracker]: https://www.pivotaltracker.com/n/projects/2126318
[CI]: https://syslog.ci.cf-app.com
[syslogLinux]: https://github.com/cloudfoundry/syslog-release
