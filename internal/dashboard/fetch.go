// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"regexp"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// Fetch retrieves a dashboard using the Circonus API
func Fetch(client API, id, title string) (*api.Dashboard, error) {
	// logger := log.With().Str("cmd", "cosi dashboard fetch").Logger()

	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}

	if id != "" {
		db, err := FetchByID(client, id)
		if err != nil {
			return nil, errors.Wrap(err, "dashboard by id")
		}
		return db, nil
	} else if title != "" {
		db, err := FetchByTitle(client, title)
		if err != nil {
			return nil, errors.Wrap(err, "dashboard by title")
		}
		return db, nil
	}

	return nil, errors.Errorf("missing required argument identifying which dashboard to fetch")
}

// FetchByID retrieves a dashboard by CID using Circonus API
func FetchByID(client API, id string) (*api.Dashboard, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	cid := id
	if !strings.HasPrefix(cid, "/dashboard/") {
		cid = "/dashboard/" + id
	}
	if ok, err := regexp.MatchString(`^/dashboard/[0-9]+`, cid); err != nil {
		return nil, errors.Wrap(err, "compile dashboard id regexp")
	} else if !ok {
		return nil, errors.Errorf("invalid dashboard id (%s)", id)
	}
	db, err := client.FetchDashboard(api.CIDType(&cid))
	if err != nil {
		return nil, errors.Wrap(err, "fetch api")
	}
	return db, nil
}

// FetchByTitle retrieves a dashboard by Display Title using Circonus API
func FetchByTitle(client API, name string) (*api.Dashboard, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if name == "" {
		return nil, errors.Errorf("invalid title (empty)")
	}

	query := api.SearchQueryType("\"" + name + "\" (active:1)")

	dbs, err := client.SearchDashboards(&query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "search api")
	}

	if dbs == nil {
		return nil, errors.Errorf("no dashboard or error returned")
	}

	if len(*dbs) == 0 {
		return nil, errors.Errorf("no dashboards found matching (%s)", name)
	}

	if len(*dbs) > 1 {
		return nil, errors.Errorf("multiple dashboards matching (%s)", name)
	}

	db := (*dbs)[0]
	return &db, nil
}
