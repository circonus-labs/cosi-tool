// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/circonus-labs/cosi-tool/internal/dashboard"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (d *Dashboards) create(id string) error {
	if id == "" {
		return errors.Errorf("invalid id (empty)")
	}

	t, _, err := d.templates.Load(d.regDir, id)
	if err != nil {
		return errors.Wrap(err, "loading template")
	}

	if len(t.Configs) == 0 {
		return errors.Errorf("%s invalid template (no configs)", id)
	}

	// TODO: will need to revisit and carve out the "dashboard_instance"
	//       vars from the old cosi plugin cmd support.
	//       the below meta data only adds partial support for the way the
	//       plugin command does its job. (postgress and cassandra use
	//       instances - meaning multiple copies of the same dashboard template)

	// set up the template expansion data
	// NOTE: all dashboards in a single template get the SAME set of basic
	//       template vars. a combination of local system items as well as
	//       any k:v data from a meta configuration file.
	tvars := map[string]string{}

	if m, err := d.loadMeta(id); err == nil {
		for k, v := range m {
			tvars[k] = v
		}
	} else {
		d.logger.Warn().Err(err).Msg("loading dashobard meta data file")
	}
	// do system vars after to ensure they are not overwritten
	tvars["HostName"] = d.config.Host.Name
	tvars["CheckUUID"] = d.checkInfo.CheckUUID

	for dashName, cfg := range t.Configs {
		dashID := id + "-" + dashName

		d.logger.Info().Str("id", dashID).Msg("building dashboard")

		dcfg, err := parseDashboardTemplate(dashID, cfg.Template, tvars)
		if err != nil {
			return err
		}

		for widx, wcfg := range cfg.Widgets {
			delete(tvars, "GraphUUID")
			if wcfg.GraphName != "" {
				if g, found := d.graphInfo[wcfg.GraphName]; found {
					tvars["GraphUUID"] = g.UUID
				} else {
					d.logger.Fatal().Str("dashboard_id", dashID).Str("graph_id", wcfg.GraphName).Interface("ginfo", d.graphInfo).Msg("graph not found")
				}
			}
			widget, werr := parseWidgetTemplate(fmt.Sprintf("%s-%d", dashID, widx), wcfg.Template, tvars)
			if werr != nil {
				return werr
			}
			dcfg.Widgets = append(dcfg.Widgets, *widget)
		}

		if e := log.Debug(); e.Enabled() {
			cfgFile := path.Join(d.regDir, "config-"+dashID+".json")
			d.logger.Debug().Str("cfg_file", cfgFile).Msg("saving registration config")
			if err := regfiles.Save(cfgFile, cfg, true); err != nil {
				return errors.Wrapf(err, "saving config (%s)", cfgFile)
			}
		}

		dash, err := dashboard.Create(d.client, dcfg)
		if err != nil {
			return err
		}
		regFile := path.Join(d.regDir, "registration-"+dashID+".json")
		if err := regfiles.Save(regFile, dash, true); err != nil {
			return errors.Wrapf(err, "saving %s registration", id)
		}
		d.dashList[dashID] = dash
	}

	return nil
}

// loadMeta reads a dashboard meta data file if found ('meta-<id>' e.g. meta-dashboard-system.json).
// NOTE: keys in meta data must start with an uppercase character and match *EXACTLY* what is used
//       in the template values should be strings - regardless of how they are used in the template:
//       e.g. meta `"Key":"1"`, template `"something":{{.Key}}` or `"sometthing":"{{.Key}}"`
func (d *Dashboards) loadMeta(id string) (map[string]string, error) {
	if id == "" {
		return nil, errors.New("invalid id (empty)")
	}
	metaFile := path.Join(d.regDir, id+".json")
	data, err := ioutil.ReadFile(metaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		if e := log.Debug(); e.Enabled() {
			fmt.Println(string(data))
			d.logger.Error().Err(err).Msg("parsing meta config")
		}
		return nil, err
	}
	return v, nil
}
