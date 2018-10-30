// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// CreateFromFile uses supplied configuration file to create a worksheet
func CreateFromFile(client CircAPI, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi worksheet create").Logger()

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

	var cfg circapi.Worksheet
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	w, err := Create(client, &cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error creating worksheet")
	}

	if err = regfiles.Save(out, w, force); err != nil {
		return errors.Wrap(err, "saving created worksheet")
	}

	return nil
}

// Create uses Circonus CircAPI to create a worksheet
func Create(client CircAPI, cfg *circapi.Worksheet) (*circapi.Worksheet, error) {
	if client == nil {
		return nil, errors.New("invalid client (nil)")
	}
	if cfg == nil {
		return nil, errors.New("invalid config (nil)")
	}
	return client.CreateWorksheet(cfg)
}
