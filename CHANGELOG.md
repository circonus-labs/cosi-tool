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
