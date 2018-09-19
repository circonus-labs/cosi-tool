// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

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

func (c *Checks) parseTemplateConfig(cfgType, cfgName string, templateVars interface{}) (*circapi.CheckBundle, error) {
	if cfgType == "" {
		return nil, errors.New("invalid config type (empty)")
	}
	if cfgName == "" {
		return nil, errors.New("invalid config name (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	templateID := cfgType + "-" + cfgName
	t, _, err := c.templates.Load(c.regDir, templateID)
	if err != nil {
		return nil, err
	}
	if len(t.Configs) == 0 {
		return nil, errors.New("invalid template (zero configs)")
	}

	// get the named configuration, (e.g. in template-check-system.toml it would be [configs.system])
	conf, ok := t.Configs[cfgName]
	if !ok {
		return nil, errors.Errorf("unable to find 'configs.%s' in template file", cfgName)
	}
	if conf.Template == "" {
		return nil, errors.Errorf("empty 'configs.%s.template' found", cfgName)
	}

	// expand template w/data
	tmpl, err := template.New(templateID).Parse(conf.Template)
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	if err := tmpl.Execute(bw, templateVars); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Printf("%#v\n%#v\n", conf.Template, templateVars)
		}
		return nil, errors.Wrap(err, "executing template")
	}
	bw.Flush()

	// create check bundle config
	var cfg circapi.CheckBundle
	if err := json.Unmarshal(b.Bytes(), &cfg); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing template config")
	}

	return &cfg, nil
}
