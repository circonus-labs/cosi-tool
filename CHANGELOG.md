# v0.6.0

* upd: release config
* upd: lint enforcement
* upd: dependencies
* fix: handle fs graphs creating invalid file names
* fix: handle missing graphs in dashbaord (remove graph)
* doc: package comments

# v0.5.3

* add: more tests covering non-active check bundles
* fix: fail when attempting to use a bad check bundle
* fix: gofrs/uuid changed to not follow modules std
* upd: dependencies

# v0.5.2

* doc: add more information to readme for subcommands
* upd: emit error message when a broker with no IP address is rejected

# v0.5.1

* fix: option names in config file

# v0.5.0

* upd: short->long(w/stream tags) metric name mapping right before submitting graph configuration via api. if short metric names are mapped in situ, template/text will html entity encode then unicode escape the stream tags (no bueno).
* add: `registration complete` info message at end of `register`
* fix: default cosi url `https://setup.circonus.com/`, no longer `onestep`
* fix: derive format from --config extension for `cosi config init --config=<filespec>`
* upd: clean up temporary files in testdata before running graph creation tests
* upd: dependencies
* add: matching static and variable datapoint metric specs in graph templates to active metrics from agent (support metrics from agent with dynamic stream tags)
* upd: remove 'REV:' prefix from check target (when cosi creates a check for reverse)

# v0.4.2

* fix: remove default register show-config format, otherwise it forces exit

# v0.4.1

* upd: dependency (gofrs/uuid/v3, yaml v2.2.2)

# v0.4.0

* upd: move register show-config to configure method (before broker selection)
* upd: add default API token app name
* upd: create etc directory if it does not exist (should already exist, if installed by cosi installer)
* upd: use api config consts for config keys
* upd: finish config of group check
* upd: switch to modules
* upd: normalize mock api names for clarification in test code
* upd: only update check bundle if NOT using metric filters
* upd: default to metric_filters for checks
* upd: build tag go1.11
* upd: USE github.com/circonus-labs/go-apiclient for circonus api
* upd: DEPRECATE github.com/circonus-labs/circonus-gometrics/api for circonus api
* upd: normalize configuration options
* add: `--show-config=fmt` to `cosi register` for skeleton configuration
* doc: wip - flesh out documentation
* upd: dependencies (circonus-agent, go-apiclient, zerolog)
* upd: `SearchCheckBundles` method signature
* upd: add zerolog shim for apiclient Logger interface

# v0.3.2

* upd: dependencies

# v0.3.1

* upd: switch to circonus-labs release target

# v0.3.0

* doc: pre-release

# v0.2.8

* upd: release file names use x86_64, facilitate automated builds and testing
* upd: fix repo refs to circonus-labs in Vagrantfile
* upd: to use go1.11 in Vagrantfile

# v0.2.7

* upd: ensure `cosi_placeholder` metric is disabled

# v0.2.6

* upd: switch to gofrs/uuid from satori/go.uuid see [issue #84](https://github.com/satori/go.uuid/issues/84)

# v0.2.5

* upd: finalize dependencies w/promoted cosi-server and updates
* fix: incorrect toml import
* upd: refactor/condense, prep for pre-alpha release

# v0.2.4

* add: example custom registration options config file

# v0.2.3

* fix: release info for `cosi version` command

# v0.2.2

* add: reset command
* fix: cosi id file name `.cosi_id` not `.cosi.id`

# v0.2.1

* upd: turn off draft
* add: rulesets to registration (last step)
* upd: change 'Config' to 'Options' (same as other reg components)

# v0.2.0

* add: implement `cosi ruleset list` command
* upd: load a common `cosi_id` - generate and save, if needed
* upd: use common notes, append any notes provided in template
* upd: integrate base ui url (check, dashboard, graph, worksheet)
* doc: update registration long description
* doc: remove ui url todo
* upd: remove templates key from main cosi config (it's a registration-specific config option)
* add: base ui url if not set, fetch and save to main cosi config (if it exists)
* add: save config
* add: implement 'cosi ruleset fetch' command
* upd: ensure ruleset registration files go in regular registration directory, so reset will function correctly
* add: implementation of 'cosi ruleset delete' command
* upd: rename default settings vars (ruleset create)
* fix: typo in graph command line option
* upd: base ruleset command description
* add: tests for ruleset package
* add: implementation of 'cosi ruleset create' command
* upd: refactor and condense graph

# v0.1.0

* Initial _foundation_ feature complete (e.g. `cosi register` creates check and visuals)
