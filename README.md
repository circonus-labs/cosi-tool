# PRE-RELEASE

# Circonus One Step Install

Circonus One Step Install (cosi) is comprised of two separate pieces.

1. [`cosi-tool`](https://github.com/circonus-labs/cosi-tool), this repository, contains the command line tool used to register a system with Circonus and manage the local registration.
1. [`cosi-server`](https://github.com/circonus-labs/cosi-server) contains the server used during the installation and registration process. It serves the installation script, whether a specific OS is supported, what [`circonus-agent`](https://github.com/circonus-labs/circonus-agent) package to use, and templates for creating assets in the Circonus system (checks, dashboards, graphs, rulesets, and worksheets).

---

# Circonus One Step Install Tool

> NOTE: this repository is in **active** development! As such, it may not be entirely feature complete, may contain bugs, and should **not** be used in production at this time.

## Installation (automated)

The new cosi ecosystem is currently in pre-alpha. When it is available for public testing this section will be updated.

## Installation (manual)

1. Download from [latest release](https://github.com/circonus-labs/cosi-tool/releases/latest)
1. Create an installation directory (e.g. `mkdir -p /opt/circonus/cosi`)
1. Unpack release archive into installation directory
1. See `bin/cosi --help`
    1. Configure `etc/example-cosi.json` (edit, rename `cosi.json` - see `cosi config -h` to get started)
    1. Optionally, configure `etc/example-reg-conf.toml` for customizing the registration portion - if applicable
