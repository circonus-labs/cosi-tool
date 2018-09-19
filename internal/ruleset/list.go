// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

// List local cosi graphs
import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type detail struct {
	metric  string // main metric
	checkID string // the check id
	numRule int
	uiurl   string
}

// List local cosi rulesets
func List(client API, w io.Writer, uiURL, regDir string, quiet, long bool) error {
	logger := log.With().Str("cmd", "cosi ruleset list").Logger()

	// NOTE: the rule_set api object does not contain sufficient information to
	//       translate directly into a ui url. the actual id of the rule set is
	//       NOT included in the api object.

	// the best we can do is return a search:
	// uiURL+"/fault-detection/rules?search=cpu`idle (check_id:253229)"

	if regDir == "" {
		return errors.Errorf("invalid registration directory (empty)")
	}

	regs, err := regfiles.Find(regDir, "ruleset")
	if err != nil {
		return errors.Wrap(err, "loading ruleset registrations")
	}

	format := "%-5s %s\n"

	if !long && !quiet {
		fmt.Fprintf(w, format, "#Rule", "Metric")
	}

	for _, regFile := range *regs {
		d, err := getDetail(client, regDir, regFile, uiURL)
		if err != nil {
			logger.Warn().Err(err).Str("dir", regDir).Str("file", regFile).Msg("skipping...")
			continue
		}
		show(w, d, format, long)
	}

	return nil
}

func getDetail(client API, regDir, regFile, uiURL string) (*detail, error) {
	if regDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if regFile == "" {
		return nil, errors.New("invalid ruleset registration file (empty)")
	}

	var rs api.RuleSet

	data, err := ioutil.ReadFile(path.Join(regDir, regFile))
	if err != nil {
		return nil, errors.Wrap(err, "reading ruleset registration file")
	}

	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, errors.Wrap(err, "parsing ruleset registration json")
	}

	d := detail{
		metric:  rs.MetricName,
		checkID: strings.Replace(rs.CheckCID, "/check/", "", 1),
		numRule: len(rs.Rules),
	}
	d.uiurl = uiURL + "fault-detection/rules?search=" + url.QueryEscape(fmt.Sprintf("%s (check_id:%s)", d.metric, d.checkID))

	return &d, nil
}

func show(w io.Writer, d *detail, format string, long bool) error {

	if !long {
		fmt.Fprintf(w, format, fmt.Sprintf("%d", d.numRule), d.metric)
		return nil
	}

	fmt.Fprintln(w, "===========")
	fmt.Fprintf(w, "No. Rules : %d\n", d.numRule)
	fmt.Fprintf(w, "Metric    : %s\n", d.metric)
	fmt.Fprintf(w, "URL       : %s\n", d.uiurl)

	return nil
}
