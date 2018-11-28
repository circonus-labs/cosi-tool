// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// CreateFromFile uses Circonus API to create a graph from supplied configuration file
func CreateFromFile(client CircAPI, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi graph create").Logger()

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

	var cfg circapi.Graph
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	g, err := Create(client, &cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error creating graph")
	}

	if err = regfiles.Save(out, g, force); err != nil {
		return errors.Wrap(err, "saving created graph")
	}

	return nil
}

// Create graph from supplied configuration, returns created graph or error
func Create(client CircAPI, cfg *circapi.Graph) (*circapi.Graph, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}

	return client.CreateGraph(cfg)
}
