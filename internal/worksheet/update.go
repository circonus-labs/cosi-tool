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

// Update uses Circonus CircAPI to update a worksheet from supplied configuration file
func Update(client CircAPI, in, out string, force bool) error {
	// logger := log.With().Str("cmd", "cosi worksheet update").Logger()

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

	c, err := client.UpdateWorksheet(&cfg)
	if err != nil {
		return errors.Wrap(err, "Circonus API error updating worksheet")
	}

	if err = regfiles.Save(out, c, force); err != nil {
		return errors.Wrap(err, "saving updated worksheet")
	}

	return nil
}
