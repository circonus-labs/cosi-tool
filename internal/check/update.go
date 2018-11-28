// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

import (
	"encoding/json"
	"io/ioutil"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// UpdateFromFile uses Circonus API to update a check from supplied configuration file
func UpdateFromFile(client CircAPI, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi check update").Logger()

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

	var cfg circapi.CheckBundle
	if err = json.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(err, "loading configuration")
	}

	c, err := Update(client, &cfg)
	if err != nil {
		return err
	}

	if err = regfiles.Save(out, c, force); err != nil {
		return errors.Wrap(err, "saving updated check bundle")
	}

	return nil
}

// Update uses Circonus API to update a check
func Update(client CircAPI, cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
	if client == nil {
		return nil, errors.New("invalid client (nil)")
	}

	if cfg == nil {
		return nil, errors.New("invalid cofnig (nil)")
	}

	c, err := client.UpdateCheckBundle(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Circonus API error updating check bundle")
	}

	return c, nil
}
