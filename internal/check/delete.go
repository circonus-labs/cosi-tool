// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// Delete uses Circonus API to delete a check from supplied configuration file, type or id
func Delete(client CircAPI, regDir, id, checkType, in string) error {
	// logger := log.With().Str("cmd", "cosi check delete").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}

	cid := ""

	if id != "" {
		cid = id
		if !strings.HasPrefix(cid, "/check_bundle/") {
			cid = "/check_bundle/" + id
		}
		if ok, err := regexp.MatchString(`^/check_bundle/[0-9]+`, cid); err != nil {
			return errors.Wrap(err, "compile check id regexp")
		} else if !ok {
			return errors.Errorf("invalid check bundle id (%s)", id)
		}
	} else if in != "" {
		data, err := ioutil.ReadFile(in)
		if err != nil {
			return errors.Wrap(err, "reading configuration file")
		}

		var cfg circapi.CheckBundle
		if err = json.Unmarshal(data, &cfg); err != nil {
			return errors.Wrap(err, "loading configuration")
		}

		cid = cfg.CID
	} else if checkType != "" {
		if regDir == "" {
			return errors.Errorf("invalid registration directory (empty)")
		}
		if ok, err := regexp.MatchString(`^(system|group)$`, checkType); err != nil {
			return errors.Wrap(err, "compile check type regexp")
		} else if !ok {
			return errors.Errorf("invalid check type (%s)", checkType)
		}

		regFile := filepath.Join(regDir, "registration-check-"+checkType+".json")
		data, err := ioutil.ReadFile(regFile)
		if err != nil {
			return errors.Wrap(err, "loading check type")
		}

		var c *circapi.CheckBundle
		if err := json.Unmarshal(data, &c); err != nil {
			return errors.Wrap(err, "parsing json")
		}

		cid = c.CID
	}

	if cid == "" {
		return errors.New("missing required argument identifying check bundle")
	}

	ok, err := client.DeleteCheckBundleByCID(&cid)
	if err != nil {
		return errors.Wrap(err, "Circonus API error deleting check bundle")
	}
	if !ok {
		return errors.New("unable to delete check bundle")
	}

	return nil
}
