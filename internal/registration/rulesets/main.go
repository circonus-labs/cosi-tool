// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package rulesets

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NOTE: rulesets are not driven off of templates. the user must provide the
//       configuration files for rulesets in a specific directory.
//
//       /opt/circonus/cosi/rulesets - any JSON files in this directory will
//                                     be loaded as rulesets and submitted to
//                                     the API.

// Options defines rulsets registration configuration options
type Options struct {
	CheckInfo  *checks.CheckInfo
	Client     CircAPI
	Config     *options.Options
	RegDir     string
	RulesetDir string
}

// Rulesets defines a rulesets registration
type Rulesets struct {
	checkInfo  *checks.CheckInfo
	client     CircAPI
	config     *options.Options
	logger     zerolog.Logger
	regDir     string
	regPrefix  string
	regFiles   *[]string
	rulesetDir string
}

// New creates a new rulesets registration instance
func New(cfg *Options) (*Rulesets, error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg (nil)")
	}
	if cfg.CheckInfo == nil {
		return nil, errors.New("invalid check info (nil)")
	}
	if cfg.Client == nil {
		return nil, errors.New("invalid client (nil)")
	}
	if cfg.Config == nil {
		return nil, errors.New("invalid options config (nil)")
	}
	if cfg.RegDir == "" {
		return nil, errors.New("invalid registration directory (empty)")
	}
	if cfg.RulesetDir == "" {
		return nil, errors.New("invalid ruleset directory (empty)")
	}

	rs := &Rulesets{
		checkInfo:  cfg.CheckInfo,
		client:     cfg.Client,
		config:     cfg.Config,
		regDir:     cfg.RegDir,
		regPrefix:  "registration-ruleset-",
		rulesetDir: cfg.RulesetDir, //path.Join(defaults.BasePath, "rulesets"),
		logger:     log.With().Str("cmd", "register.rulesets").Logger(),
	}

	return rs, nil
}

// Register checks for rulesets configurations and creates registrations for any
// ruleset configurations found.
func (rs *Rulesets) Register() error {

	files, err := ioutil.ReadDir(rs.rulesetDir)
	if err != nil {
		if os.IsNotExist(err) {
			rs.logger.Warn().Err(err).Msg("skipping, rulesets")
			return nil // if the directory does not exist, just skip registering rulesets
		}
		return err
	}

	regs, err := regfiles.Find(rs.regDir, "ruleset")
	if err != nil {
		return err
	}
	rs.regFiles = regs

	for _, file := range files {
		if !file.Mode().IsRegular() {
			rs.logger.Debug().Str("file", file.Name()).Msg("not a regular file, skipping")
			continue
		}
		if filepath.Ext(file.Name()) != ".json" {
			rs.logger.Debug().Str("file", file.Name()).Msg("extension not '.json', skipping")
			continue
		}
		if strings.HasPrefix(file.Name(), rs.regPrefix) {
			continue
		}
		if rs.checkForRegistration(strings.Replace(file.Name(), ".json", "", 1)) {
			rs.logger.Debug().Str("file", file.Name()).Msg("registration found, skipping")
			continue
		}
		if err := rs.create(file.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (rs *Rulesets) create(cfgFile string) error {
	data, err := ioutil.ReadFile(path.Join(rs.rulesetDir, cfgFile))
	if err != nil {
		return err
	}

	var cfg circapi.RuleSet
	if err := json.Unmarshal(data, &cfg); err != nil {
		rs.logger.Error().Err(err).Msg("parsing ruleset config")
		return err
	}

	cfg.CheckCID = rs.checkInfo.CheckCID
	notes := rs.config.Common.Notes
	if cfg.Notes != nil {
		notes += *cfg.Notes
	}
	cfg.Notes = &notes
	if len(rs.config.Common.Tags) > 0 {
		cfg.Tags = append(cfg.Tags, rs.config.Common.Tags...)
	}

	rso, err := rs.client.CreateRuleSet(&cfg)
	if err != nil {
		return err
	}

	return regfiles.Save(path.Join(rs.regDir, rs.regPrefix+cfgFile), rso, true)
}

// checkForRegistration looks through existing registration files to see
// if the ruleset has already been created.
func (rs *Rulesets) checkForRegistration(id string) bool {
	if id == "" {
		return false
	}

	regFileSig := rs.regPrefix + id
	for _, rf := range *rs.regFiles {
		if strings.Contains(rf, regFileSig) {
			return true
		}
	}

	return false
}
