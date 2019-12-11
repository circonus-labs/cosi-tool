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

// Fetch retrieves a check using the Circonus API
func Fetch(client CircAPI, regDir, id, checkType, name, target string) (*circapi.CheckBundle, error) {
	// logger := log.With().Str("cmd", "cosi check fetch").Logger()

	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}

	if id != "" {
		b, err := FetchByID(client, id)
		if err != nil {
			return nil, errors.Wrap(err, "check by id")
		}
		return b, nil
	}

	if checkType != "" {
		b, err := FetchByType(client, regDir, checkType)
		if err != nil {
			return nil, errors.Wrap(err, "check by type")
		}
		return b, nil
	}

	if name != "" {
		b, err := FetchByName(client, name)
		if err != nil {
			return nil, errors.Wrap(err, "check by name")
		}
		return b, nil
	}

	if target != "" {
		b, err := FetchByTarget(client, target)
		if err != nil {
			return nil, errors.Wrap(err, "check by target")
		}
		return b, nil
	}

	return nil, errors.Errorf("missing required argument identifying which check to fetch")
}

// FetchByID retrieves a check bundle by CID using Circonus API
func FetchByID(client CircAPI, id string) (*circapi.CheckBundle, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	cid := id
	if !strings.HasPrefix(cid, "/check_bundle/") {
		cid = "/check_bundle/" + id
	}
	if ok, err := regexp.MatchString(`^/check_bundle/[0-9]+`, cid); err != nil {
		return nil, errors.Wrap(err, "compile check id regexp")
	} else if !ok {
		return nil, errors.Errorf("invalid check bundle id (%s)", id)
	}
	b, err := client.FetchCheckBundle(circapi.CIDType(&cid))
	if err != nil {
		return nil, errors.Wrap(err, "fetch api")
	}
	return b, nil
}

// FetchByType retrieves a check bundle by COSI check type (system|group) using Circonus API
func FetchByType(client CircAPI, regDir, checkType string) (*circapi.CheckBundle, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if regDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if checkType == "" {
		return nil, errors.New("invalid check type (empty)")
	}

	if ok, err := regexp.MatchString(`^(system|group)$`, checkType); err != nil {
		return nil, errors.Wrap(err, "compile check type regexp")
	} else if !ok {
		return nil, errors.Errorf("invalid check type (%s)", checkType)
	}

	regFile := filepath.Join(regDir, "registration-check-"+checkType+".json")
	data, err := ioutil.ReadFile(regFile)
	if err != nil {
		return nil, errors.Wrap(err, "loading check type")
	}

	var b *circapi.CheckBundle
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, errors.Wrap(err, "parsing json")
	}

	return FetchByID(client, b.CID)
}

// FetchByName retrieves a check bundle by Display Name using Circonus API
func FetchByName(client CircAPI, name string) (*circapi.CheckBundle, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if name == "" {
		return nil, errors.Errorf("invalid display name (empty)")
	}

	query := circapi.SearchQueryType("\"" + name + "\" (active:1)")

	bundles, err := client.SearchCheckBundles(&query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "search api")
	}

	if bundles == nil {
		return nil, errors.Errorf("no check bundles or error returned")
	}

	if len(*bundles) == 0 {
		return nil, errors.Errorf("no checks found matching (%s)", name)
	}

	if len(*bundles) > 1 {
		return nil, errors.Errorf("multiple checks matching (%s)", name)
	}

	b := (*bundles)[0]
	return &b, nil
}

// FetchByTarget retrieves a check bundle by Target using Circonus API
func FetchByTarget(client CircAPI, target string) (*circapi.CheckBundle, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if target == "" {
		return nil, errors.Errorf("invalid target (empty)")
	}

	query := circapi.SearchQueryType("(host:\"" + target + "\")(active:1)")

	bundles, err := client.SearchCheckBundles(&query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "search api")
	}

	if bundles == nil {
		return nil, errors.Errorf("no check bundles or error returned")
	}

	if len(*bundles) == 0 {
		return nil, errors.Errorf("no checks found matching (%s)", target)
	}

	if len(*bundles) > 1 {
		return nil, errors.Errorf("multiple checks matching (%s)", target)
	}

	b := (*bundles)[0]
	return &b, nil
}
