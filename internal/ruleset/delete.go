// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// Delete uses Circonus API to delete a ruleset from supplied configuration file or id
func Delete(client CircAPI, id, in string) error {
	// logger := log.With().Str("cmd", "cosi ruleset delete").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}

	cid := ""

	if id != "" {
		cid = id
		if !strings.HasPrefix(cid, "/rule_set/") {
			cid = "/rule_set/" + id
		}
		if ok, err := regexp.MatchString(`^/rule_set/[0-9]+`, cid); err != nil {
			return errors.Wrap(err, "compile ruleset id regexp")
		} else if !ok {
			return errors.Errorf("invalid ruleset id (%s)", id)
		}
	} else if in != "" {
		data, err := ioutil.ReadFile(in)
		if err != nil {
			return errors.Wrap(err, "reading configuration file")
		}

		var cfg circapi.RuleSet
		if err = json.Unmarshal(data, &cfg); err != nil {
			return errors.Wrap(err, "loading configuration")
		}

		cid = cfg.CID
	}

	if cid == "" {
		return errors.New("missing required argument identifying ruleset")
	}

	ok, err := client.DeleteRuleSetByCID(&cid)
	if err != nil {
		return errors.Wrap(err, "Circonus API error deleting ruleset")
	}
	if !ok {
		return errors.New("unable to delete ruleset")
	}

	return nil
}
