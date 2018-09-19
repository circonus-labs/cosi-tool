// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import (
	"regexp"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// Fetch retrieves a ruleset using the Circonus API
func Fetch(client API, id string) (*api.RuleSet, error) {
	// logger := log.With().Str("cmd", "cosi ruleset fetch").Logger()
	return FetchByID(client, id)
}

// FetchByID retrieves a ruleset by CID using Circonus API
func FetchByID(client API, id string) (*api.RuleSet, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	cid := id
	if !strings.HasPrefix(cid, "/rule_set/") {
		cid = "/rule_set/" + id
	}
	if ok, err := regexp.MatchString(`^/rule_set/[0-9]+`, cid); err != nil {
		return nil, errors.Wrap(err, "compile ruleset id regexp")
	} else if !ok {
		return nil, errors.Errorf("invalid ruleset id (%s)", id)
	}
	rs, err := client.FetchRuleSet(api.CIDType(&cid))
	if err != nil {
		return nil, errors.Wrap(err, "fetch api")
	}
	return rs, nil
}
