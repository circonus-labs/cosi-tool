// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

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
	cosiType     string    // type of cosi check (e.g. system|group)
	name         string    // display name of check
	checkType    string    // type of check (e.g. json:nad, httptrap, etc.)
	metrics      uint      // number of active metrics
	lastModified time.Time // last modified time
	status       string    // whether check has been modified
	checkURLs    []string  // list of check URLs
}

// List local cosi checks
func List(client CircAPI, w io.Writer, uiURL, regDir string, quiet, verify, long bool) error {
	logger := log.With().Str("cmd", "cosi check list").Logger()

	if regDir == "" {
		return errors.Errorf("invalid registration directory (empty)")
	}

	checks, err := regfiles.Find(regDir, "check")
	if err != nil {
		return errors.Wrap(err, "loading check registrations")
	}

	// quiet := viper.GetBool(KeyQuiet)
	// verifyCheck := viper.GetBool(KeyVerify)
	// long := viper.GetBool(KeyLong)
	format := "%-6s %-40s %-10s %8s %-21s %-8s\n"

	if !long && !quiet {
		fmt.Fprintf(w, format, "ID", "Name", "Type", "#Active", "Modified", "Status")
	}

	for _, regFile := range *checks {
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
		return nil, errors.New("invalid check registration file (empty)")
	}

	var b circapi.CheckBundle

	data, err := ioutil.ReadFile(path.Join(regDir, regFile))
	if err != nil {
		return nil, errors.Wrap(err, "reading check registration file")
	}

	if err := json.Unmarshal(data, &b); err != nil {
		return nil, errors.Wrap(err, "parsing check registration json")
	}

	cosiType := strings.Replace(regFile, "registration-", "", 1)
	cosiType = strings.Replace(cosiType, "check-", "", 1)
	cosiType = strings.Replace(cosiType, ".json", "", 1)

	lm := time.Unix(int64(b.LastModified), 0)

	status := "n/a"
	if verifyCheck {
		modified, err := isModified(client, b.CID, b.LastModified)
		if err != nil {
			return nil, errors.Wrapf(err, "verifying check %s", b.DisplayName)
		}
		status = "OK"
		if modified {
			status = "Modified"
		}
	}

	checkURLs := []string{}
	for _, cid := range b.Checks {
		checkURLs = append(checkURLs, uiURL+strings.Replace(cid, "/check/", "checks/", 1))
	}

	d := detail{
		cosiType:     cosiType,
		name:         b.DisplayName,
		checkType:    b.Type,
		metrics:      uint(len(b.Metrics)),
		lastModified: lm,
		status:       status,
		checkURLs:    checkURLs,
	}

	return &d, nil
}

func show(w io.Writer, d *detail, format string, long bool) {

	if !long {
		fmt.Fprintf(w, format, d.cosiType, d.name, d.checkType, fmt.Sprintf("%d", d.metrics), d.lastModified.Format(time.RFC822Z), d.status)
		return
	}

	fmt.Fprintln(w, "================")
	fmt.Fprintf(w, "ID        : %s\n", d.cosiType)
	fmt.Fprintf(w, "Name      : %s\n", d.name)
	fmt.Fprintf(w, "Type      : %s\n", d.checkType)
	fmt.Fprintf(w, "Metrics   : %d\n", d.metrics)
	fmt.Fprintf(w, "Modified  : %s\n", d.lastModified.Format(time.RFC1123Z))
	fmt.Fprintf(w, "Status    : %s\n", d.status)
	for _, chkURL := range d.checkURLs {
		fmt.Fprintf(w, "Check URL : %s\n", chkURL)
	}
}
