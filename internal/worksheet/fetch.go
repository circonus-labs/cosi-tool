// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

import (
	"regexp"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// Fetch retrieves a worksheet using the Circonus CircAPI
func Fetch(client CircAPI, id, title string) (*circapi.Worksheet, error) {
	// logger := log.With().Str("cmd", "cosi worksheet fetch").Logger()

	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}

	if id != "" {
		w, err := FetchByID(client, id)
		if err != nil {
			return nil, errors.Wrap(err, "worksheet by id")
		}
		return w, nil
	} else if title != "" {
		w, err := FetchByTitle(client, title)
		if err != nil {
			return nil, errors.Wrap(err, "worksheet by title")
		}
		return w, nil
	}

	return nil, errors.Errorf("missing required argument identifying which worksheet to fetch")
}

// FetchByID retrieves a worksheet by CID using Circonus CircAPI
func FetchByID(client CircAPI, id string) (*circapi.Worksheet, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	cid := id
	if !strings.HasPrefix(cid, "/worksheet/") {
		cid = "/worksheet/" + id
	}
	if ok, err := regexp.MatchString(`^/worksheet/[0-9]+`, cid); err != nil {
		return nil, errors.Wrap(err, "compile worksheet id regexp")
	} else if !ok {
		return nil, errors.Errorf("invalid worksheet id (%s)", id)
	}
	db, err := client.FetchWorksheet(circapi.CIDType(&cid))
	if err != nil {
		return nil, errors.Wrap(err, "fetch api")
	}
	return db, nil
}

// FetchByTitle retrieves a worksheet by Display Title using Circonus CircAPI
func FetchByTitle(client CircAPI, name string) (*circapi.Worksheet, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}
	if name == "" {
		return nil, errors.Errorf("invalid title (empty)")
	}

	query := circapi.SearchQueryType("\"" + name + "\" (active:1)")

	dbs, err := client.SearchWorksheets(&query, nil)
	if err != nil {
		return nil, errors.Wrap(err, "search api")
	}

	if dbs == nil {
		return nil, errors.Errorf("no worksheet or error returned")
	}

	if len(*dbs) == 0 {
		return nil, errors.Errorf("no worksheets found matching (%s)", name)
	}

	if len(*dbs) > 1 {
		return nil, errors.Errorf("multiple worksheets matching (%s)", name)
	}

	db := (*dbs)[0]
	return &db, nil
}
