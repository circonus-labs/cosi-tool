// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

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
	cid         string // cid of graph
	description string // description of cosi graph
	title       string // title of graph
	uiurl       string
}

// List local cosi graphs
func List(client CircAPI, w io.Writer, uiURL, regDir string, quiet, long bool) error {
	logger := log.With().Str("cmd", "cosi graph list").Logger()

	if regDir == "" {
		return errors.Errorf("invalid registration directory (empty)")
	}

	graphs, err := regfiles.Find(regDir, "graph")
	if err != nil {
		return errors.Wrap(err, "loading graph registrations")
	}

	// TODO: workout short list formatting during testing
	format := "%-20s %-40s %-40s\n"

	if !long && !quiet {
		fmt.Fprintf(w, format, "ID", "Title", "Description")
	}

	for _, regFile := range *graphs {
		d, err := getDetail(client, regDir, regFile, uiURL)
		if err != nil {
			logger.Warn().Err(err).Str("dir", regDir).Str("file", regFile).Msg("skipping...")
			continue
		}
		show(w, d, format, long)
	}

	return nil
}

func getDetail(client CircAPI, regDir, regFile, uiURL string) (*detail, error) {
	if regDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if regFile == "" {
		return nil, errors.New("invalid graph registration file (empty)")
	}

	var g circapi.Graph

	data, err := ioutil.ReadFile(path.Join(regDir, regFile))
	if err != nil {
		return nil, errors.Wrap(err, "reading graph registration file")
	}

	if err := json.Unmarshal(data, &g); err != nil {
		return nil, errors.Wrap(err, "parsing graph registration json")
	}

	d := detail{
		cid:         strings.Replace(g.CID, "/graph/", "", 1),
		description: g.Description,
		title:       g.Title,
	}
	d.uiurl = uiURL + "trending/graphs/view/" + d.cid

	return &d, nil
}

func show(w io.Writer, d *detail, format string, long bool) error {

	if !long {
		fmt.Fprintf(w, format, d.cid, d.title, d.description)
		return nil
	}

	fmt.Fprintln(w, "================")
	fmt.Fprintf(w, "CID        : %s\n", d.cid)
	fmt.Fprintf(w, "Title      : %s\n", d.title)
	fmt.Fprintf(w, "Description: %s\n", d.description)
	fmt.Fprintf(w, "URL        : %s\n", d.uiurl)

	return nil
}
