// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type detail struct {
	cid          string    // cid of dashboard
	uuid         string    // uuid used for dashboard url in long
	cosiType     string    // type of cosi dashboard
	title        string    // title of dashboard
	lastModified time.Time // last modified time
	status       string    // whether dashboard has been modified
	uiurl        string
}

// List local cosi dashboards
func List(client CircAPI, w io.Writer, uiURL, regDir string, quiet, verify, long bool) error {
	logger := log.With().Str("cmd", "cosi dashboard list").Logger()

	if regDir == "" {
		return errors.Errorf("invalid registration directory (empty)")
	}

	dashboards, err := regfiles.Find(regDir, "dashboard")
	if err != nil {
		return errors.Wrap(err, "loading dashboard registrations")
	}

	// TODO: workout short list formatting during testing
	format := "%-20s %-10s %-21s %-8s %-40s\n"

	if !long && !quiet {
		fmt.Fprintf(w, format, "ID", "CID", "Modified", "Status", "Title")
	}

	for _, regFile := range *dashboards {
		d, err := getDetail(client, regDir, regFile, uiURL, verify)
		if err != nil {
			logger.Warn().Err(err).Str("dir", regDir).Str("file", regFile).Msg("skipping...")
			continue
		}
		show(w, d, format, long)
	}

	return nil
}

func getDetail(client CircAPI, regDir, regFile, uiURL string, verifyCheck bool) (*detail, error) {
	if regDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if regFile == "" {
		return nil, errors.New("invalid dashboard registration file (empty)")
	}

	var db circapi.Dashboard

	data, err := ioutil.ReadFile(path.Join(regDir, regFile))
	if err != nil {
		return nil, errors.Wrap(err, "reading dashboard registration file")
	}

	if err := json.Unmarshal(data, &db); err != nil {
		return nil, errors.Wrap(err, "parsing dashboard registration json")
	}

	cosiType := strings.Replace(regFile, "registration-", "", 1)
	cosiType = strings.Replace(cosiType, "dashboard-", "", 1)
	cosiType = strings.Replace(cosiType, ".json", "", 1)

	lm := time.Unix(int64(db.LastModified), 0)

	status := "n/a"
	if verifyCheck {
		modified, err := isModified(client, db.CID, db.LastModified)
		if err != nil {
			return nil, errors.Wrapf(err, "verifying dashboard %s", db.Title)
		}
		status = "OK"
		if modified {
			status = "Modified"
		}
	}

	d := detail{
		cid:          strings.Replace(db.CID, "/dashboard/", "", 1),
		uuid:         db.UUID,
		cosiType:     cosiType,
		title:        db.Title,
		lastModified: lm,
		status:       status,
	}

	d.uiurl = uiURL + "dashboards/view/" + d.uuid

	return &d, nil
}

func show(w io.Writer, d *detail, format string, long bool) error {

	if !long {
		fmt.Fprintf(w, format, d.cosiType, d.cid, d.lastModified.Format(time.RFC822Z), d.status, d.title)
		return nil
	}

	fmt.Fprintln(w, "================")
	fmt.Fprintf(w, "ID       : %s\n", d.cosiType)
	fmt.Fprintf(w, "CID      : %s\n", d.cid)
	fmt.Fprintf(w, "Title    : %s\n", d.title)
	fmt.Fprintf(w, "Modified : %s\n", d.lastModified.Format(time.RFC1123Z))
	fmt.Fprintf(w, "Status   : %s\n", d.status)
	fmt.Fprintf(w, "URL      : %s\n", d.uiurl)

	return nil
}
