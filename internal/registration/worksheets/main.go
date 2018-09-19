// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"strings"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/circonus-labs/cosi-tool/internal/worksheet"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Worksheets defines the registration instance
type Worksheets struct {
	worksheetList map[string]*circapi.Worksheet
	client        CircAPI
	config        *options.Options
	regDir        string
	templates     *templates.Templates
	regFiles      *[]string
	logger        zerolog.Logger
}

// Options defines the settings required to create a new instance
type Options struct {
	Client    CircAPI
	Config    *options.Options
	RegDir    string
	Templates *templates.Templates
}

// New creates a new Worksheets instance
func New(o *Options) (*Worksheets, error) {
	if o == nil {
		return nil, errors.New("invalid options (nil)")
	}
	if o.Client == nil {
		return nil, errors.New("invalid client (nil)")
	}
	if o.Config == nil {
		return nil, errors.New("invalid options config (nil)")
	}
	if o.RegDir == "" {
		return nil, errors.New("invalid reg dir (empty)")
	}
	if o.Templates == nil {
		return nil, errors.New("invalid templates (nil)")
	}

	regs, err := regfiles.Find(o.RegDir, "worksheet")
	if err != nil {
		return nil, errors.Wrap(err, "finding worksheet registrations")
	}

	w := Worksheets{
		worksheetList: make(map[string]*circapi.Worksheet),
		client:        o.Client,
		config:        o.Config,
		regDir:        o.RegDir,
		templates:     o.Templates,
		regFiles:      regs,
		logger:        log.With().Str("cmd", "register.worksheets").Logger(),
	}

	return &w, nil
}

// Register creates worksheets using the Circonus API
func (w *Worksheets) Register(list map[string]bool) error {
	if len(list) == 0 {
		return errors.New("invalid list (empty)")
	}
	// list of _all_ things to create during registration
	// filter out non-worksheet items
	worksheetList := make(map[string]bool)
	for k, v := range list {
		if strings.HasPrefix(k, "worksheet-") {
			worksheetList[k] = v
		}
	}
	if len(worksheetList) == 0 {
		w.logger.Warn().Msg("0 worksheets found in list, not building ANY worksheets")
		return nil
	}

	for id, create := range worksheetList {
		if !create {
			continue
		}

		if err := w.create(id); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worksheets) create(id string) error {
	if id == "" {
		return errors.Errorf("invalid id (empty)")
	}

	t, _, err := w.templates.Load(w.regDir, id)
	if err != nil {
		return errors.Wrap(err, "loading template")
	}

	if len(t.Configs) == 0 {
		return errors.Errorf("%s invalid template (no configs)", id)
	}

	// set up the template expansion data
	type templateVars struct {
		HostName string
	}
	tvars := templateVars{
		HostName: w.config.Host.Name,
	}

	for cfgName, wcfg := range t.Configs {
		worksheetID := id + "-" + cfgName

		w.logger.Info().Str("id", worksheetID).Msg("building worksheet")

		if loaded, err := w.checkForRegistration(worksheetID); err != nil {
			return err
		} else if loaded {
			continue
		}

		cfg, err := parseTemplateConfig(worksheetID, wcfg.Template, tvars)
		if err != nil {
			return err
		}

		cfg.SmartQueries = []circapi.WorksheetSmartQuery{
			{
				Name:  "Circonus One Step Install",
				Order: []string{},
				Query: `(notes:"` + w.config.Common.Notes + `*")`,
			},
		}

		if len(w.config.Common.Tags) > 0 {
			cfg.Tags = append(cfg.Tags, w.config.Common.Tags...)
		}
		if len(w.config.Checks.System.Tags) > 0 {
			cfg.Tags = append(cfg.Tags, w.config.Checks.System.Tags...)
		}
		// add note
		notes := w.config.Common.Notes
		if cfg.Notes != nil {
			notes += *cfg.Notes
		}
		cfg.Notes = &notes

		if e := log.Debug(); e.Enabled() {
			cfgFile := path.Join(w.regDir, "config-"+worksheetID+".json")
			w.logger.Debug().Str("cfg_file", cfgFile).Msg("saving registration config")
			if err := regfiles.Save(cfgFile, cfg, true); err != nil {
				return errors.Wrapf(err, "saving config (%s)", cfgFile)
			}
		}

		sheet, err := worksheet.Create(w.client, cfg)
		if err != nil {
			return err
		}
		regFile := path.Join(w.regDir, "registration-"+worksheetID+".json")
		if err := regfiles.Save(regFile, sheet, true); err != nil {
			return errors.Wrapf(err, "saving %s registration", worksheetID)
		}
		w.worksheetList[worksheetID] = sheet
	}

	return nil
}

// checkForRegistration looks through existing registration files and
// if the id is found, it is loaded. Returns a boolean indicating
// if the registration was found+loaded successfully or an error
func (w *Worksheets) checkForRegistration(id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid id (empty)")
	}

	regFileSig := "registration-" + id
	for _, rf := range *w.regFiles {
		if !strings.Contains(rf, regFileSig) {
			continue
		}
		var ws circapi.Worksheet
		found, err := regfiles.Load(path.Join(w.regDir, rf), &ws)
		if err != nil {
			return false, errors.Wrapf(err, "loading %s", regFileSig)
		}
		if found {
			w.worksheetList[regFileSig] = &ws
			return found, nil
		}
		break // we already found it but there was an issue
	}

	return false, nil
}

func parseTemplateConfig(id, templateCfg string, templateVars interface{}) (*circapi.Worksheet, error) {
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	if templateCfg == "" {
		return nil, errors.New("invalid template config (empty)")
	}
	if templateVars == nil {
		return nil, errors.New("invalid template vars (nil)")
	}

	// parse template
	tmpl, err := template.New(id).Parse(templateCfg)
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
	var ws circapi.Worksheet
	if err := json.Unmarshal(b.Bytes(), &ws); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(b.String())
		}
		return nil, errors.Wrap(err, "parsing expanded template result")
	}

	return &ws, nil
}
