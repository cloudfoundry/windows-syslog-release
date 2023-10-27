# Windows Syslog BOSH Release

* Syslog Release for Linux: [Release][syslogLinux]

This is a BOSH release for forwarding logs from BOSH jobs on Windows VMs. It forwards logs in c:/var/vcap/sys/log/ (and
any subdirectories, recursively) to a configured syslog server.

If you have any questions, or want to get attention for a PR or issue please reach out on the [#logging-and-metrics
channel in the cloudfoundry slack](https://cloudfoundry.slack.com/archives/CUW93AF3M)

## Differences From `syslog-release`

This release is designed to share configuration with its linux sibling. However, it uses a very different architecture,
and parity would often be prohibitively expensive. Here are a list of some of the major differences with
`syslog-release`, with special attention to standard configuration that will be ignored or fail.

- RELP is not a supported transport protocol. While `syslog-forwarder-windows` will start, it will not forward logs.
- Does not support fallback servers (and will not use any if configured)
- Does not support custom rules (and will not respect them if configured)
- Does not support configuration of `permitted_peer` (what addresses or dns resolutions are allowed to recieve syslog messages) in tls mode
- Does not support custom configuration of various queuing parameters

If any of these constraints are a problem for you, please open an issue explaining your use case.

## Configuring Log Forwarding

Add the `syslog_forwarder` to forward all local syslog messages from an instance to a syslog endpoint. You can use
`addons` to add syslog forwarder to all instances.

Configure `address`, `port` and `transport`:

```yml
instance_groups:
- name: some-instance-group
  jobs:
  - name: syslog_forwarder_windows
    release: windows-syslog
  properties:
    syslog:
      address: <IP or hostname>
      port: <port>
      transport: <transport>
...
releases:
- name: windows-syslog
  version: <version>
  url: https://bosh.io/d/github.com/cloudfoundry/windows-syslog-release?v=<version>
  sha1: <sha>
```

you can also add the `event_logger` job to receive events

```yml
- name: event_logger
  release: windows-syslog
```

You can get the windows-syslog-release from [bosh.io/releases](https://bosh.io/releases/github.com/cloudfoundry/windows-syslog-release?all=1).

If the syslog endpoint is unavailable, messages will be briefly queued.

TLS over TCP is supported with additional properties. `tls-enabled` should be set to true if you wish to use it. In a
future version, this will likely be true by default. By default, the windows certificate API is used to validate certs.
If the cert you wish to respect isn't validated by the Windows API, you will need to set the full cert chain with
`ca_cert`. Note that this fully replaces use of the Windows API.

## Testing and Debugging

Here are some tips for debugging issues with log forwarding. We're assuming some familiarity with BOSH, but not
necessarily with Windows.

There is a useful manifest at `tests/manifests/tcp.yml`. To use it, you will need to provide a deployment name, like so:
`bosh deploy -d windows-syslog -v deployment=windows-syslog tests/manifests/tcp.yml`.

To ssh on to a windows VM, you will need the windows utility release. Our test manifest includes it. You can then use
`bosh -d windows-syslog ssh --opts=-T forwarder "echo example command"` to run a single command non interactively on the
forwarder job. Interactive SSH works more or less like you'd expect, except that you arrive in `cmd.exe`.

Our test manifest includes a storer from the linux syslog release, and the windows forwarder job is configured to send
logs to it.

[syslogLinux]: https://github.com/cloudfoundry/syslog-release
