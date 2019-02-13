# Circonus One Step Install

Circonus One Step Install (cosi) is comprised of two separate pieces.

1. [`cosi-tool`](https://github.com/circonus-labs/cosi-tool), this repository, contains the command line tool used to register a system with Circonus and manage the local registration.
1. [`cosi-server`](https://github.com/circonus-labs/cosi-server) contains the server used during the installation and registration process. It serves the installation script, whether a specific OS is supported, what [`circonus-agent`](https://github.com/circonus-labs/circonus-agent) package to use, and templates for creating assets in the Circonus system (checks, dashboards, graphs, rulesets, and worksheets).

The [circonus-agent](https://github.com/circonus-labs/circonus-agent) is comprised of:

  * replacement for NAD, written in go, with builtin plugins for the common metrics needed for cosi visuals (graphs, worksheets, & dashboards)
  * includes (if OS supports) [protocol_observer](https://github.com/circonus-labs/wirelatency), no longer needs to be built/installed manually
  * includes (if OS supports) [circonus-logwatch](https://github.com/circonus-labs/circonus-logwatch), no longer needs to be installed manually
  * includes OS/version/architecture-specific NAD plugins (non-javascript only) -- **Note:** the circonus-agent is **not** capable of using NAD _native plugins_ since they require NodeJS

The cosi-tool does **not** currently include a functional `cosi plugin` command. This capability will be included in a future release, as the individual `cosi plugin ...` sub-commands (postgres and cassandra) are completed.

Supported Operating Systems (x86_64 and/or amd64):

  * RHEL7 (CentOS, RedHat, Oracle)
  * RHEL6 (CentOS, RedHat, amzn)
  * Ubuntu18
  * Ubuntu16
  * Ubuntu14
  * Debian9
  * Debian8
  * FreeBSD 12
  * FreeBSD 11

Please continue to use the original cosi(w/NAD) for OmniOS and Raspian - cosi v2 support for these is TBD. Note: after installing NAD a binary circonus-agent can be used as a drop-in replacement (configure circonus-agent _plugins directory_ to be NAD plugins directory -- javascript plugins will not function). Binaries for OmniOS (`solaris_x86_64`) and Raspian (`linux_arm`) are available in the [circonus-agent repository](https://github.com/circonus-labs/circonus-agent/releases/latest).

---

# Circonus One Step Install Tool

## Installation (automated)

```
curl -sSL https://setup.circonus.com/install | bash \
    -s -- \
    --cosiurl https://setup.circonus.com/ \
    --key <insert api key> \
    --app <insert api app>
```

## Installation (manual)

1. Download from [latest release](https://github.com/circonus-labs/cosi-tool/releases/latest)
1. Create an installation directory (e.g. `mkdir -p /opt/circonus/cosi`)
1. Unpack release archive into installation directory
1. See `bin/cosi --help`
    1. Configure `etc/example-cosi.json` (edit, rename `cosi.json` - see `cosi config -h` to get started)
    1. Optionally, configure `etc/example-reg-conf.toml` for customizing the registration portion - if applicable

## Options (general)

```
$ /opt/circonus/cosi/bin/cosi -h
A command line tool for registering a system with Circonus
and managing the local registration.

Usage:
  cosi [flags]
  cosi [command]

Available Commands:
  broker      Information about Circonus Brokers
  check       Manage COSI registered check(s)
  config      COSI configuration file
  dashboard   Manage COSI registered dashboard(s)
  graph       Manage COSI registered graph(s)
  help        Help about any command
  plugin      Manage specific NAD plugins
  register    COSI registration of this system
  reset       Reset system - remove COSI created artifacts
  ruleset     Manage rulesets for the system check
  template    Manage COSI templates
  version     Display version and exit
  worksheet   Manage COSI registered worksheet(s)

  Flags:
        --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
        --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
        --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name
        --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
        --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
        --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
        --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
        --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
        --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "<hostname>")
    -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
        --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://onestep.circonus.com/")
    -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
        --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
    -h, --help                  help for cosi
        --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
        --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
        --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
        --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
        --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
        --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
        --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
        --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)

Use "cosi [command] --help" for more information about a command.
```

### Configuration file options

```

agent:
  mode: ""
  url: ""
api:
  app: ""
  ca_file: ""
  key: ""
  url: ""
base_ui_url: ""
cosi_url: ""
debug: false
host:
  broker:
    id: ""
    type: ""
  group_id: ""
  target: ""
log:
  level: ""
  pretty: false
reg_conf: ""
system:
  arch: ""
  dmi: ""
  os_dist: ""
  os_type: ""
  os_vers: ""

```

### Registration configuration file

```

brokers:
  group:
    list: []
    default: 0
  system:
    list: []
    default: 0
checks:
  group:
    broker_id: ""
    create: false
    display_name: ""
    id: ""
    tags: []
  system:
    broker_id: ""
    display_name: ""
    tags: []
    target: ""
dashboards:
  system:
    create: false
    title: ""
graphs:
  configs: {}
  exclude: []
  include: []
host:
  ip: ""
  name: ""
worksheets:
  system:
    create: false
    title: ""
    description: ""
    tags: []

```

## Commands

### Broker

```
$ /opt/circonus/cosi/bin/cosi broker -h
Obtain information on Circonus Brokers available to the token being used.

Usage:
  cosi broker [command]

Available Commands:
  default     Show default broker
  list        List available brokers
  show        Show information for a specific broker

Flags:
  -h, --help   help for broker

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Check

```
$ /opt/circonus/cosi/bin/cosi check -h
Intended for managing local COSI checks.

Usage:
  cosi check [command]

Available Commands:
  create      Create a check from a configuration file
  delete      Delete a check from Circonus
  fetch       Fetch an existing check bundle from API
  list        List checks
  update      Update a check using configuration file

Flags:
  -h, --help   help for check

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Config

```
$ /opt/circonus/cosi/bin/cosi config -h
Initialize and show the COSI configuration file.

Usage:
  cosi config [command]

Available Commands:
  init        Create an initial default configuration file
  show        Display current configuration

Flags:
  -h, --help   help for config

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Dashboard

```
$ /opt/circonus/cosi/bin/cosi dashboard -h
Intended for managing local COSI dashboards.

Usage:
  cosi dashboard [command]

Available Commands:
  create      Create a dashboard from a configuration file
  delete      Delete a dashboard from Circonus
  fetch       Fetch an existing dashboard from API
  list        List dashboards
  update      Update a dashboard using configuration file

Flags:
  -h, --help   help for dashboard

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Graph

```
$ /opt/circonus/cosi/bin/cosi graph -h
Intended for managing local COSI graphs.

Usage:
  cosi graph [command]

Available Commands:
  create      Create a graph from a configuration file
  delete      Delete a graph from Circonus
  fetch       Fetch an existing graph from API
  list        List graphs
  update      Update a graph using configuration file

Flags:
  -h, --help   help for graph

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Register

> Note: the system check will always be created. All of the other items (group check, graphs, worksheets, dashboards, rulesets) are optional.

```
$ /opt/circonus/cosi/bin/cosi register -h
Register this system using COSI method.

Create a system check and optional group check.
Create graphs for known plugins.
Create system worksheet.
Create system dashboard.
Create rulesets for system check.

Usage:
  cosi register [flags]

Flags:
  -h, --help                 help for register
      --show-config string   Show registration options configuration using format yaml|json|toml
      --templates strings    Template ID list (type-name[,type-name,...] e.g. check-system,graph-cpu)

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Reset

> NOTE: when `--force` is _not_ used, reset will prompt for confirmation.

```
$ /opt/circonus/cosi/bin/cosi reset -h
Reset will delete all COSI registration created artifacts.
Checks, graphs, worksheets, rulesets, dashboards and associated registration files.

Usage:
  cosi reset [flags]

Flags:
      --force   Do not prompt for confirmation
  -h, --help    help for reset

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Ruleset

```
$ /opt/circonus/cosi/bin/cosi ruleset -h
Intended for managing rulesets for the system check.

Usage:
  cosi ruleset [command]

Available Commands:
  create      Create a ruleset from a configuration file
  delete      Delete a ruleset from Circonus
  fetch       Fetch an existing ruleset from API
  list        List rulesets

Flags:
  -h, --help   help for ruleset

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Template

```
$ /opt/circonus/cosi/bin/cosi template -h
Intended for managing local templates for COSI.

Usage:
  cosi template [command]

Available Commands:
  fetch       Fetch an existing template from COSI API
  list        A brief description of your command

Flags:
  -h, --help   help for template

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

### Worksheet

```
$ /opt/circonus/cosi/bin/cosi worksheet -h
Intended for managing local COSI worksheets.

Usage:
  cosi worksheet [flags]
  cosi worksheet [command]

Available Commands:
  create      Create a worksheet from a configuration file
  delete      Delete a worksheet from Circonus
  fetch       Fetch an existing worksheet from API
  list        List worksheets
  update      Update a worksheet using configuration file

Flags:
  -h, --help   help for worksheet

Global Flags:
      --agent-mode string     [ENV: COSI_AGENT_MODE] Agent mode for check (reverse|pull) (default "reverse")
      --agent-url string      [ENV: COSI_AGENT_URL] URL the Circonus Agent is listening on (default "http://localhost:2609/")
      --api-app string        [ENV: COSI_API_APP] Circonus API Token App Name (default "cosi")
      --api-ca-file string    [ENV: COSI_API_CA_FILE] Circonus API Certificate CA file
      --api-key string        [ENV: COSI_API_KEY] Circonus API Token Key
      --api-url string        [ENV: COSI_API_URL] Circonus API URL (default "https://api.circonus.com/v2/")
      --broker-id uint        [ENV: COSI_BROKER_ID] Broker ID to use when creating check [0=auto select] (default 0)
      --broker-type string    [ENV: COSI_BROKER_TYPE] Limit automatic broker selection to a specific type of broker (default "any")
      --check-target string   [ENV: COSI_CHECK_TARGET] Check target(host) to use when creating system check (default "cosi-tool-c7")
  -c, --config string         config file (default: /opt/circonus/cosi/etc/cosi.yaml|.json|.toml)
      --cosi-url string       [ENV: COSI_URL] Circonus One Step Install (cosi server) URL (default "https://setup.circonus.com/")
  -d, --debug                 [ENV: COSI_DEBUG] Enable debug messages
      --group-id string       [ENV: COSI_GROUP_ID] Group ID for multi-system check
      --log-level string      [ENV: COSI_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty            [ENV: COSI_LOG_PRETTY] Output formatted/colored log lines [ignored on windows] (default true)
      --os-distro string      [ENV: COSI_OS_DISTRO] OS distribution (generated by cosi-install)
      --os-type string        [ENV: COSI_OS_TYPE] OS type (generated by cosi-install)
      --os-version string     [ENV: COSI_OS_VERSION] OS distribution version (generated by cosi-install)
      --regconf string        [ENV: COSI_REG_CONF] Registration options configuration file
      --sys-arch string       [ENV: COSI_SYS_ARCH] System architecture (generated by cosi-install)
      --sys-dmi string        [ENV: COSI_SYS_DMI] System dmi bios version (generated by cosi-install, only used in AWS)
```

---

Unless otherwise noted, the source files are distributed under the BSD-style license found in the [LICENSE](LICENSE) file.
