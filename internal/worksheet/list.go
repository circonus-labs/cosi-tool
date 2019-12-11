// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type detail struct {
	cid         string // cid of worksheet
	description string // description of cosi worksheet
	title       string // title of worksheet
	uiurl       string
}

// List local cosi worksheets
func List(client CircAPI, w io.Writer, uiURL, regDir string, quiet, long bool) error {
	logger := log.With().Str("cmd", "cosi worksheet list").Logger()

	if regDir == "" {
		return errors.Errorf("invalid registration directory (empty)")
	}

	worksheets, err := regfiles.Find(regDir, "worksheet")
	if err != nil {
		return errors.Wrap(err, "loading worksheet registrations")
	}

	// TODO: workout short list formatting during testing
	format := "%-20s %-40s %-40s\n"

	if !long && !quiet {
		fmt.Fprintf(w, format, "ID", "Title", "Description")
	}

	for _, regFile := range *worksheets {
		d, err := getDetail(regDir, regFile, uiURL)
		if err != nil {
			logger.Warn().Err(err).Str("dir", regDir).Str("file", regFile).Msg("skipping...")
			continue
		}
		show(w, d, format, long)
	}

	return nil
}

func getDetail(regDir, regFile, uiURL string) (*detail, error) {
	if regDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if regFile == "" {
		return nil, errors.New("invalid worksheet registration file (empty)")
	}

	var w circapi.Worksheet

	data, err := ioutil.ReadFile(path.Join(regDir, regFile))
	if err != nil {
		return nil, errors.Wrap(err, "reading worksheet registration file")
	}

	if err := json.Unmarshal(data, &w); err != nil {
		return nil, errors.Wrap(err, "parsing worksheet registration json")
	}

	d := detail{
		cid:         strings.Replace(w.CID, "/worksheet/", "", 1),
		description: *w.Description,
		title:       w.Title,
	}

	d.uiurl = uiURL + "trending/worksheets/" + d.cid

	return &d, nil
}

func show(w io.Writer, d *detail, format string, long bool) {

	if !long {
		fmt.Fprintf(w, format, d.cid, d.title, d.description)
		return
	}

	fmt.Fprintln(w, "================")
	fmt.Fprintf(w, "CID        : %s\n", d.cid)
	fmt.Fprintf(w, "Title      : %s\n", d.title)
	fmt.Fprintf(w, "Description: %s\n", d.description)
	fmt.Fprintf(w, "URL        : %s\n", d.uiurl)
}
