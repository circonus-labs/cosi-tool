// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func parseDashboardTemplate(dashID, templateCfg string, templateVars interface{}) (*circapi.Dashboard, error) {
	if dashID == "" {
		return nil, errors.New("invalid dashboard id (empty)")
	}
	if templateCfg == "" {
		return nil, errors.New("invalid template config (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	// parse template
	tmpl, err := template.New(dashID).Parse(templateCfg)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	tmpl = tmpl.Option("missingkey=error")
	// expand template w/data
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	if err := tmpl.Execute(bw, templateVars); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Printf("%#v\n%#v\n", templateCfg, templateVars)
		}
		return nil, errors.Wrap(err, "executing template")
	}
	bw.Flush()

	// create config
	var dash circapi.Dashboard
	if err := json.Unmarshal(b.Bytes(), &dash); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing expanded template result")
	}

	return &dash, nil
}

func parseWidgetTemplate(dashID, templateCfg string, templateVars interface{}) (*circapi.DashboardWidget, error) {
	if dashID == "" {
		return nil, errors.New("invalid dashboard id (empty)")
	}
	if templateCfg == "" {
		return nil, errors.New("invalid template config (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	// parse template
	tmpl, err := template.New(dashID).Parse(templateCfg)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	tmpl = tmpl.Option("missingkey=error")
	// expand template w/data
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	if err := tmpl.Execute(bw, templateVars); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Printf("%#v\n%#v\n", templateCfg, templateVars)
		}
		return nil, errors.Wrap(err, "executing template")
	}
	bw.Flush()

	// create config
	var widget circapi.DashboardWidget
	if err := json.Unmarshal(b.Bytes(), &widget); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing expanded template result")
	}

	return &widget, nil
}
