// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// Delete uses Circonus API to delete a dashboard from supplied configuration file or id
func Delete(client API, id, in string) error {
	// logger := log.With().Str("cmd", "cosi dashboard delete").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}

	cid := ""

	if id != "" {
		cid = id
		if !strings.HasPrefix(cid, "/dashboard/") {
			cid = "/dashboard/" + id
		}
		if ok, err := regexp.MatchString(`^/dashboard/[0-9]+`, cid); err != nil {
			return errors.Wrap(err, "compile dashboard id regexp")
		} else if !ok {
			return errors.Errorf("invalid dashboard id (%s)", id)
		}
	} else if in != "" {
		data, err := ioutil.ReadFile(in)
		if err != nil {
			return errors.Wrap(err, "reading configuration file")
		}

		var cfg api.Dashboard
		if err = json.Unmarshal(data, &cfg); err != nil {
			return errors.Wrap(err, "loading configuration")
		}

		cid = cfg.CID
	}

	if cid == "" {
		return errors.New("missing required argument identifying dashboard")
	}

	ok, err := client.DeleteDashboardByCID(&cid)
	if err != nil {
		return errors.Wrap(err, "Circonus API error deleting dashboard")
	}
	if !ok {
		return errors.New("unable to delete dashboard")
	}

	return nil
}
