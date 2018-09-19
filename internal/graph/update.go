// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pkg/errors"
)

// Update uses Circonus API to update a graph from supplied configuration file
func Update(client API, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi graph update").Logger()

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

	var cfg api.Graph
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	c, err := client.UpdateGraph(&cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error updating graph")
	}

	if err = regfiles.Save(out, c, force); err != nil {
		return errors.Wrap(err, "saving updated graph")
	}

	return nil
}
