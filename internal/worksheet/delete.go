// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// Delete uses Circonus API to delete a worksheet from supplied configuration file or id
func Delete(client API, id, in string) error {
	// logger := log.With().Str("cmd", "cosi worksheet delete").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}

	cid := ""

	if id != "" {
		cid = id
		if !strings.HasPrefix(cid, "/worksheet/") {
			cid = "/worksheet/" + id
		}
		if ok, err := regexp.MatchString(`^/worksheet/[0-9]+`, cid); err != nil {
			return errors.Wrap(err, "compile worksheet id regexp")
		} else if !ok {
			return errors.Errorf("invalid worksheet id (%s)", id)
		}
	} else if in != "" {
		data, err := ioutil.ReadFile(in)
		if err != nil {
			return errors.Wrap(err, "reading configuration file")
		}

		var cfg api.Worksheet
		if err = json.Unmarshal(data, &cfg); err != nil {
			return errors.Wrap(err, "loading configuration")
		}

		cid = cfg.CID
	}

	if cid == "" {
		return errors.New("missing required argument identifying worksheet")
	}

	ok, err := client.DeleteWorksheetByCID(&cid)
	if err != nil {
		return errors.Wrap(err, "Circonus API error deleting worksheet")
	}
	if !ok {
		return errors.New("unable to delete worksheet")
	}

	return nil
}
