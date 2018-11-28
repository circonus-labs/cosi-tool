// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// CreateFromFile create dashboard from supplied configuration file
func CreateFromFile(client CircAPI, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi dashboard create").Logger()

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

	var cfg circapi.Dashboard
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	d, err := Create(client, &cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error creating dashboard")
	}

	if err = regfiles.Save(out, d, force); err != nil {
		return errors.Wrap(err, "saving created dashboard")
	}

	return nil
}

// Create dashboard using Circonus API
func Create(client CircAPI, cfg *circapi.Dashboard) (*circapi.Dashboard, error) {
	if client == nil {
		return nil, errors.New("invalid client (nil)")
	}
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}
	return client.CreateDashboard(cfg)
}
