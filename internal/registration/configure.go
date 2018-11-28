// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"math/rand"
	"os"
	"time"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (r *Registration) configure() error {
	if err := verifyRegDir(r.regDir, r.logger); err != nil {
		return err
	}

	if cfg, err := options.LoadConfigFile(viper.GetString(config.KeyRegConf)); err == nil {
		r.config = cfg
	} else {
		return err
	}

	if viper.GetString(KeyShowConfig) != "" {
		options.DumpConfig(r.config, viper.GetString(KeyShowConfig), os.Stdout)
		os.Exit(0)
	}

	// Set brokers for checks if they are not already set
	r.logger.Info().Msg("selecting system check broker")
	if bid, err := r.selectBroker("json"); err == nil {
		r.logger.Info().Str("broker", bid).Msg("system check broker")
		r.config.Checks.System.BrokerID = bid
	} else {
		return errors.Wrap(err, "selecting system check broker")
	}

	if r.config.Checks.Group.Create {
		r.logger.Info().Msg("selecting group check broker")
		if bid, err := r.selectBroker("httptrap"); err == nil {
			r.logger.Info().Str("broker", bid).Msg("group check broker")
			r.config.Checks.System.BrokerID = bid
		} else {
			return errors.Wrap(err, "selecting group check broker")
		}
	}

	// available metrics
	if metrics, err := getAvailableMetrics(viper.GetString(config.KeyAgentURL), r.logger); err == nil {
		r.availableMetrics = metrics
	} else {
		return err
	}

	// set the templates to attempt to create
	// NOTE: the system check will ALWAYS be created
	r.logger.Info().Msg("setting template list")
	list := viper.GetStringSlice(KeyTemplateList)
	if len(list) == 0 {
		r.logger.Info().Msg("using default template list")
		if lst, err := templates.DefaultTemplateList(r.availableMetrics); err == nil {
			list = *lst
		} else {
			return errors.Wrap(err, "setting default template list")
		}
	}
	for _, t := range list {
		r.templateList[t] = true
	}
	r.logger.Info().Strs("list", list).Msg("templates")
	return nil
}

// verifyRegDir verifies the registration directrory exists, attempst to create
// if the directory is not found.
func verifyRegDir(regDir string, logger zerolog.Logger) error {
	logger.Info().Str("reg_dir", regDir).Msg("verify registration directory")
	fi, err := os.Stat(regDir)
	if err != nil {
		logger.Warn().Err(err).Msg("attempting to create")
		if os.IsNotExist(err) {
			if merr := os.MkdirAll(regDir, 0755); merr != nil {
				return errors.Wrap(merr, "registration directory")
			}
		} else {
			return errors.Wrap(err, "registration directory")
		}
	} else if !fi.IsDir() {
		return errors.Errorf("registration directory (%s) not a directory", regDir)
	}

	return nil
}

// Get available metrics from agent
func getAvailableMetrics(agentURL string, logger zerolog.Logger) (*agentapi.Metrics, error) {
	if agentURL == "" {
		return nil, errors.New("invalid agent URL (empty)")
	}

	logger.Info().Str("agent_url", agentURL).Msg("fetching available metrics from agent")

	client, err := agentapi.New(agentURL)
	if err != nil {
		return nil, errors.New("creating agent API client")
	}

	metrics, err := client.Metrics("")
	if err != nil {
		return nil, errors.Wrap(err, "fetching available metrics from agent")
	}

	return metrics, nil
}
