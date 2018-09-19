// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func parseGraphTemplate(graphID, templateCfg string, templateVars interface{}) (*circapi.Graph, error) {
	if graphID == "" {
		return nil, errors.New("invalid graph id (empty)")
	}
	if templateCfg == "" {
		return nil, errors.New("invalid template config (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	// parse template
	tmpl, err := template.New(graphID).Parse(templateCfg)
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

	// create graph config
	var graph circapi.Graph
	if err := json.Unmarshal(b.Bytes(), &graph); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing expanded template result")
	}

	return &graph, nil
}

func parseDatapointTemplate(graphID, templateCfg string, templateVars interface{}) (*circapi.GraphDatapoint, error) {
	if graphID == "" {
		return nil, errors.New("invalid graph id (empty)")
	}
	if templateCfg == "" {
		return nil, errors.New("invalid template config (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	// parse template
	tmpl, err := template.New(graphID).Parse(templateCfg)
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

	// create graph datapoint config
	var dp circapi.GraphDatapoint
	if err := json.Unmarshal(b.Bytes(), &dp); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing expanded template result")
	}

	return &dp, nil
}
