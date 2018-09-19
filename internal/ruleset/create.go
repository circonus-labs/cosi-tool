// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pkg/errors"
)

// CreateFromFile uses Circonus API to create a check from supplied configuration file
func CreateFromFile(client API, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi ruleset create").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}

	if in == "" {
		return errors.New("invalid input file (empty)")
	}

	data, err := ioutil.ReadFile(in)
	if err != nil {
		return errors.Wrap(err, "reading configuration file")
	}

	var cfg api.RuleSet
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	c, err := Create(client, &cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error creating ruleset")
	}

	if err = regfiles.Save(out, c, force); err != nil {
		return errors.Wrap(err, "saving created ruleset")
	}

	return nil
}

// Create rulset from supplied configuration, returns created ruleset or error
func Create(client API, cfg *api.RuleSet) (*api.RuleSet, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}

	return client.CreateRuleSet(cfg)
}
