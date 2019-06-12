// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Checks defines the checks registration instance
type Checks struct {
	checkList map[string]*circapi.CheckBundle
	client    CircAPI
	config    *options.Options
	regDir    string
	templates *templates.Templates
	logger    zerolog.Logger
}

// Options defines the settings required to create a new checks instance
type Options struct {
	Client    CircAPI
	Config    *options.Options
	RegDir    string
	Templates *templates.Templates
}

// CheckInfo contains information used by graphs and dashboards
type CheckInfo struct {
	BundleCID     string
	CheckCID      string
	CheckID       uint
	CheckUUID     string
	SubmissionURL string // only applies to trap check (e.g. check-group)
}

const (
	statusActive = "active"
)

// New creates a new Checks instance
func New(o *Options) (*Checks, error) {
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
	c := Checks{
		checkList: make(map[string]*circapi.CheckBundle),
		client:    o.Client,
		config:    o.Config,
		regDir:    o.RegDir,
		templates: o.Templates,
		logger:    log.With().Str("cmd", "register.checks").Logger(),
	}
	return &c, nil
}

// Register creates checks using the Circonus API
func (c *Checks) Register() error {
	haveSystemCheck := false
	haveGroupCheck := false

	c.logger.Info().Str("reg_dir", c.regDir).Msg("finding registrations")
	regs, err := regfiles.Find(c.regDir, "check")
	if err != nil {
		return errors.Wrap(err, "finding check registrations")
	}
	c.logger.Debug().Strs("reg_files", *regs).Msg("check registration files")
	if len(*regs) > 0 {
		for _, regFile := range *regs {
			if strings.Contains(regFile, "registration-check-system") {
				var chk circapi.CheckBundle
				found, err := regfiles.Load(path.Join(c.regDir, regFile), &chk)
				if err != nil {
					return errors.Wrap(err, "loading check-system registration")
				}
				if found {
					c.logger.Info().Str("cid", chk.CID).Msg("found system check registration")
					c.logger.Info().Str("cid", chk.CID).Msg("fetching up-to-date check configuration via API")
					ck, err := check.FetchByID(c.client, chk.CID)
					if err != nil {
						return errors.Wrap(err, "fetching check")
					}
					if ck.Status != statusActive {
						return errors.Errorf("existing check bundle found (%s), INVALID - not active - (file:%s) (api:%s) -- please clean up artifacts from previous cosi registration", regFile, chk.Status, ck.Status)
					}
					if len(ck.Checks) == 0 {
						return errors.Errorf("existing check bundle found (%s), INVALID - has no checks (%s) -- please clean up artifacts from previous cosi registration", regFile, chk.CID)
					}
					if len(ck.CheckUUIDs) == 0 {
						return errors.Errorf("existing check bundle found (%s), INVALID - has no check uuids (%s) -- please clean up artifacts from previous cosi registration", regFile, chk.CID)
					}
					haveSystemCheck = true
					c.checkList["check-system"] = ck
				}
			} else if strings.Contains(regFile, "registration-check-group") {
				var chk circapi.CheckBundle
				found, err := regfiles.Load(path.Join(c.regDir, regFile), &chk)
				if err != nil {
					return errors.Wrap(err, "loading check-group registration")
				}
				if found {
					c.logger.Info().Str("cid", chk.CID).Msg("found group check registration")
					c.logger.Info().Str("cid", chk.CID).Msg("fetching up-to-date check configuration via API")
					ck, err := check.FetchByID(c.client, chk.CID)
					if err != nil {
						return errors.Wrap(err, "fetching check")
					}
					if ck.Status != statusActive {
						return errors.Errorf("existing check bundle found (%s), INVALID - not active (%s) -- please clean up artifacts from previous cosi registration", regFile, ck.Status)
					}
					if len(ck.Checks) == 0 {
						return errors.Errorf("existing check bundle found (%s), INVALID - has no checks (%s) -- please clean up artifacts from previous cosi registration", regFile, chk.CID)
					}
					if len(ck.CheckUUIDs) == 0 {
						return errors.Errorf("existing check bundle found (%s), INVALID - has no check uuids (%s) -- please clean up artifacts from previous cosi registration", regFile, chk.CID)
					}
					haveGroupCheck = true
					c.checkList["check-group"] = ck
				}
			} else {
				fmt.Println("unknown check type", regFile, "ignoring...")
			}
		}
	}

	if !haveSystemCheck {
		c.logger.Info().Msg("creating system check registrations")
		b, err := c.createSystemCheck()
		if err != nil {
			return err
		}
		c.logger.Info().Str("cid", b.CID).Msg("created system check")
		c.checkList["check-system"] = b
		if err := c.updateAgentConfig(); err != nil {
			return errors.Wrap(err, "updating agent config for reverse mode")
		}
	}

	if c.config.Checks.Group.Create && !haveGroupCheck {
		c.logger.Info().Msg("creating system check registrations")
		b, err := c.createGroupCheck()
		if err != nil {
			return err
		}
		c.logger.Info().Str("cid", b.CID).Msg("created group check")
		c.checkList["check-group"] = b
	}

	return nil
}

// GetCheckInfo returns information used by graphs and dashboards
func (c *Checks) GetCheckInfo(checkID string) (*CheckInfo, error) {
	if checkID == "" {
		return nil, errors.New("invalid check id (empty)")
	}
	if !strings.HasPrefix(checkID, "check-") {
		checkID = "check-" + checkID
	}

	for chkID, chk := range c.checkList {
		if chkID != checkID {
			continue
		}

		if len(chk.Checks) == 0 {
			return nil, errors.Errorf("invalid check bundle, has no checks (%s)", chk.CID)
		}

		if len(chk.CheckUUIDs) == 0 {
			return nil, errors.Errorf("invalid check bundle, has no check uuids (%s)", chk.CID)
		}

		info := CheckInfo{
			BundleCID: chk.CID,
			CheckCID:  chk.Checks[0],
			CheckUUID: chk.CheckUUIDs[0],
		}

		id, err := strconv.ParseUint(strings.Replace(chk.Checks[0], "/check/", "", 1), 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "coverting check id to uint")
		}
		info.CheckID = uint(id)

		if chk.Type == "httptrap" {
			info.SubmissionURL = chk.Config["submission_url"]
		}

		return &info, nil
	}

	return nil, errors.Errorf("check id not found (%s)", checkID)
}

func (c *Checks) createCheck(id string, cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
	if id == "" {
		return nil, errors.Errorf("invalid id (empty)")
	}
	if cfg == nil {
		return nil, errors.Errorf("invalid check bundle config (nil)")
	}

	if e := log.Debug(); e.Enabled() {
		cfgFile := path.Join(c.regDir, "config-"+id+".json")
		c.logger.Debug().Str("cfg_file", cfgFile).Msg("saving registration config")
		if err := regfiles.Save(cfgFile, cfg, true); err != nil {
			return nil, errors.Wrapf(err, "saving config (%s)", cfgFile)
		}
	}

	b, err := check.Create(c.client, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s", id)
	}

	regFile := path.Join(c.regDir, "registration-"+id+".json")
	if err := regfiles.Save(regFile, b, true); err != nil {
		return nil, errors.Wrapf(err, "saving %s registration", id)
	}

	return b, nil
}

func (c *Checks) updateAgentConfig() error {
	// TODO: TBD update agent config here or in cosi install shell script
	// Set reverse mode on
	// Set api key to "cosi"
	// /opt/circonus/agent/sbin/circonus-agentd --reverse --api-key="cosi" --api-app="cosi" --show-config=toml > /opt/circonus/agent/etc/circonus-agent.toml
	// Restart agent
	return nil
}

func (c *Checks) UpdateSystemCheck(metrics *map[string]string) error {
	cfg, ok := c.checkList["check-system"]
	if !ok {
		return errors.New("no system check found in check list")
	}

	// Do not attempt to update checks using metric_filters. Whether metrics
	// are active is handled via regular expressions in the broker itself.
	// The two methods of handling metric status are mutually exclusive.
	if len(cfg.MetricFilters) > 0 {
		return nil
	}

	// The check DOES NOT use metric_filters, metric status is being managed
	// manually by updating the check configuration. (DEPRECATED - use filters)
	updateCheck := false
	for mn, mt := range *metrics {
		c.logger.Debug().Str("metric_name", mn).Msg("find")
		found := false
		for i := 0; i < len(cfg.Metrics); i++ {
			// disable cosi_placeholder metric
			if cfg.Metrics[i].Name == "cosi_placeholder" && cfg.Metrics[i].Status == statusActive {
				cfg.Metrics[i].Status = "available"
				continue
			}
			if cfg.Metrics[i].Name != mn {
				continue
			}
			if cfg.Metrics[i].Status == statusActive {
				c.logger.Debug().Str("metric_name", mn).Msg("already active, skipping")
				found = true
				break
			}
			c.logger.Debug().Str("metric_name", mn).Msg("activating metric")
			cfg.Metrics[i].Status = statusActive
			updateCheck = true
			found = true
			break
		}
		if !found {
			c.logger.Debug().Str("metric_name", mn).Msg("adding metric")
			cfg.Metrics = append(cfg.Metrics, circapi.CheckBundleMetric{
				Name:   mn,
				Type:   mt,
				Status: statusActive,
			})
			updateCheck = true
		}
	}
	if updateCheck {
		for i := 0; i < len(cfg.Metrics); i++ {
			// disable cosi_placeholder metric
			if cfg.Metrics[i].Name == "cosi_placeholder" && cfg.Metrics[i].Status == statusActive {
				cfg.Metrics[i].Status = "available"
				break
			}
		}
		c.logger.Debug().Msg("updating check to activate new metrics")
		newCfg, err := check.Update(c.client, cfg)
		if err != nil {
			return err
		}
		c.checkList["check-system"] = newCfg
	}

	return nil
}
