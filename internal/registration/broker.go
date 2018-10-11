// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"time"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/broker"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// selectBroker uses one of several methods to determine the broker to use when creating a check.
// 1. explicit --broker on command line or  explicit broker set for check in config file
// 2. explicit list of brokers to select from for check type in config file
// 3. select from available enterprise brokers
// NOTE: if there are any enterprise brokers #5 will not be used even
//       if none of the enterprise brokers is valid for the check type.
//       to force use of a SaaS broker when enterprise brokers are available,
//       use one of the excplicit methods #1 or #2 above.
// 4. get broker from cosi-server for the check type
func (r *Registration) selectBroker(checkType string) (string, error) {
	logger := log.With().Str("cmd", "register.broker").Logger()

	if checkType == "" {
		return "", errors.New("invalid check type (empty)")
	}

	logger.Debug().Msg("fetching list of broker from Circonus API")
	// get list of brokers available to the current account
	brokers, err := broker.List(r.cliCirc)
	if err != nil {
		return "", err
	}

	//
	// explicit broker from command line or specific check section of config
	//
	{
		logger.Debug().Msg("checking config for explicit broker setting")
		valid, bid, err := r.getExplicit(checkType, brokers, &r.config.Checks)
		if err != nil {
			return "", err
		} else if valid {
			return bid, nil
		} // otherwise, invalid, fall-through to next seletion type
	}

	//
	// explicit list of brokers to use in registration configuration
	//
	{
		logger.Debug().Msg("checking config for explicit list of brokers setting")
		valid, bid, err := r.selectFromConfigList(checkType, brokers, &r.config.Brokers)
		if err != nil {
			return "", err
		} else if valid {
			return bid, nil
		} // otherwise, invalid, fall-through to next seletion type
	}

	//
	// select an enterprise broker, if any available
	{
		logger.Debug().Msg("checking for enterprise brokers")
		valid, bid, err := r.selectEnterprise(checkType, brokers)
		if err != nil {
			return "", err
		} else if valid {
			return bid, nil
		} // otherwise, invalid, fall-through to next seletion type
	}

	//
	// final fallback, select default from cosi (for SaaS, inside should NOT
	// get here since inside only uses enterprise brokers)
	//
	{
		logger.Debug().Msg("getting default broker from COSI API")
		valid, bid, err := r.getCosiDefault(checkType, brokers, r.cliCosi)
		if err != nil {
			return "", err
		} else if valid {
			return bid, nil
		}
	}

	return "", errors.New("unable to determine a valid broker to use")
}

// getExplicit broker from command line or specific check section of config
func (r *Registration) getExplicit(checkType string, brokers *[]api.Broker, cfg *options.Checks) (bool, string, error) {
	logger := log.With().Str("cmd", "register.broker").Logger()
	if checkType == "" {
		return false, "", errors.New("invalid check type (empty)")
	}
	if brokers == nil {
		return false, "", errors.New("invalid broker list (nil)")
	}
	if cfg == nil {
		return false, "", errors.New("invalid options config (nil)")
	}

	var brokerID string

	switch checkType {
	case jsonCheckType:
		brokerID = cfg.System.BrokerID
	case trapCheckType:
		brokerID = cfg.Group.BrokerID
	default:
		return false, "", errors.Errorf("unsupported check type (%s)", checkType)
	}
	if brokerID != "" {
		valid, bid, err := r.checkBroker(checkType, brokerID, brokers)
		if err != nil {
			return valid, bid, errors.Wrapf(err, "invalid broker id specified (%s)", brokerID)
		}
		logger.Debug().Str("check_type", checkType).Str("broker", bid).Msg("found broker in custom config")
		return valid, bid, err
	}

	logger.Debug().Str("check_type", checkType).Msg("no broker found in custom config")

	return false, "", nil // fall-through to next selection method
}

// selectFromConfigList will use an explicit list of brokers configured for a check type in the configuraion file
func (r *Registration) selectFromConfigList(checkType string, brokers *[]api.Broker, cfg *options.Brokers) (bool, string, error) {
	logger := log.With().Str("cmd", "register.broker").Logger()
	if checkType == "" {
		return false, "", errors.New("invalid check type (empty)")
	}
	if brokers == nil {
		return false, "", errors.New("invalid broker list (nil)")
	}
	if cfg == nil {
		return false, "", errors.New("invalid options config (nil)")
	}

	var brokerID string

	if checkType == jsonCheckType && len(cfg.System.List) > 0 {
		if cfg.System.Default >= 0 {
			if len(cfg.System.List) > cfg.System.Default {
				brokerID = cfg.System.List[cfg.System.Default]
			} else {
				return false, "", errors.New("invalid system check broker config in regconf (default out of list range)")
			}
		} else if cfg.System.Default == -1 {
			brokerID = cfg.System.List[rand.Intn(len(cfg.System.List))]
		} else {
			return false, "", errors.New("invalid system check broker config in regconf (default invalid)")
		}
	} else if checkType == trapCheckType && len(cfg.Group.List) > 0 {
		if cfg.Group.Default >= 0 {
			if len(cfg.Group.List) > cfg.Group.Default {
				brokerID = cfg.Group.List[cfg.Group.Default]
			} else {
				return false, "", errors.New("invalid group check broker config in regconf (default out of list range)")
			}
		} else if cfg.Group.Default == -1 {
			brokerID = cfg.Group.List[rand.Intn(len(cfg.Group.List))]
		} else {
			return false, "", errors.New("invalid group check broker config in regconf (default invalid)")
		}
	}
	if brokerID != "" {
		valid, bid, err := r.checkBroker(checkType, brokerID, brokers)
		if err != nil {
			return false, "", errors.Wrapf(err, "invalid broker id (%s) in regconf for check type %s", brokerID, checkType)
		} else if valid {
			logger.Debug().Str("check_type", checkType).Str("broker", bid).Msg("found broker list in custom config")
			return valid, bid, nil
		}
	}

	logger.Debug().Str("check_type", checkType).Msg("no broker list found in custom config")

	return false, "", nil // fall-through to next selection method
}

func (r *Registration) selectEnterprise(checkType string, brokers *[]api.Broker) (bool, string, error) {
	logger := log.With().Str("cmd", "register.broker").Logger()

	if checkType == "" {
		return false, "", errors.New("invalid check type (empty)")
	}
	if brokers == nil {
		return false, "", errors.New("invalid broker list (nil)")
	}

	//
	// create list of active enterprise brokers, if there are any
	//
	haveEnterpriseBrokers := false
	enterpriseBrokerList := []string{}
	for _, broker := range *brokers {
		if broker.Type != "enterprise" {
			continue
		}
		// ensure there is at least ONE instance which is active
		noneActive := true
		for _, detail := range broker.Details {
			if detail.Status == "active" {
				noneActive = false
				break
			}
		}
		if noneActive {
			continue
		}
		haveEnterpriseBrokers = true
		valid, bid, err := r.checkBroker(checkType, broker.CID, brokers)
		if err != nil {
			logger.Warn().Err(err).Msg("checking enterprise broker, skipping")
		}
		if valid {
			enterpriseBrokerList = append(enterpriseBrokerList, bid)
		}
	}
	if haveEnterpriseBrokers {
		// NOTE: if there *are* enterprise brokers, only select from enterprise brokers.
		//       this enforcement also makes cosi work for inside setups which are using
		//       the public cosi-server. (an inside setup only has enterprise brokers...)
		if len(enterpriseBrokerList) == 0 {
			return false, "", errors.New("available enterprise brokers found, none valid")
		} else if len(enterpriseBrokerList) == 1 { // only one, return it
			bid := enterpriseBrokerList[0]
			logger.Debug().Str("check_type", checkType).Str("broker", bid).Msg("found enterprise broker")
			return true, bid, nil
		} else { // otherwise, select a random one
			bid := enterpriseBrokerList[rand.Intn(len(enterpriseBrokerList))]
			logger.Debug().Str("check_type", checkType).Str("broker", bid).Msg("found more than one enterprise broker, using random one")
			return true, bid, nil
		}
	}

	logger.Debug().Msg("no viable enterprise brokers found")
	return false, "", nil
}

// getCosiDefault contacts the cosi-server to get the default broker for the check type
func (r *Registration) getCosiDefault(checkType string, brokers *[]api.Broker, client CosiAPI) (bool, string, error) {
	if checkType == "" {
		return false, "", errors.New("invalid check type (empty)")
	}
	if brokers == nil {
		return false, "", errors.New("invalid broker list (nil)")
	}
	if client == nil {
		return false, "", errors.New("invalid cosi API client (nil)")
	}

	brokerID, err := client.FetchBroker(checkType)
	if err != nil {
		return false, "", err
	}

	return r.checkBroker(checkType, brokerID, brokers)
}

// checkBroker verifies a broker has correct module for check type and is reachable
func (r *Registration) checkBroker(checkType, brokerID string, brokers *[]api.Broker) (bool, string, error) {
	if checkType == "" {
		return false, "", errors.New("invalid check type (empty)")
	}
	if brokerID == "" {
		return false, "", errors.New("invalid broker id (empty)")
	}
	if brokers == nil {
		return false, "", errors.New("invalid broker list (nil)")
	}

	if checkType != jsonCheckType && checkType != trapCheckType {
		return false, "", errors.Errorf("unknown check type (%s)", checkType)
	}

	rxBrokerCID := regexp.MustCompile(`^/broker/[0-9]+$`)
	if !rxBrokerCID.MatchString(brokerID) {
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(brokerID) {
			return false, "", errors.Errorf("invalid broker id specified (%s) - format should be '#' or '/broker/#'", brokerID)
		}
		brokerID = "/broker/" + brokerID
	}

	var connErr error
	hasModule := false
	for _, broker := range *brokers {
		if broker.CID != brokerID {
			continue
		}
		for _, instance := range broker.Details {
			if instance.Status != "active" {
				continue
			}
			for _, module := range instance.Modules {
				if module != checkType {
					continue
				}
				hasModule = true
				ip := instance.ExternalHost
				if ip == nil {
					ip = instance.IP
				}
				if ip == nil {
					continue // no external host or ip set - unreachable, active w/o ip...wtf!?
				}
				port := instance.ExternalPort
				if port == 0 {
					if instance.Port != nil {
						port = *instance.Port
					}
				}
				if port == 0 {
					port = uint16(43191) // default broker port
				}

				ok, err := brokerConnectionTest(*ip, port, r.maxBrokerResponseTime)
				if err != nil {
					// stack up the conn errors (may have tested 1-n broker instances, messy but complete...)
					if connErr != nil {
						connErr = errors.WithMessage(err, fmt.Sprintf("%s | %s", connErr.Error(), broker.Name))
					} else {
						connErr = errors.WithMessage(err, broker.Name)
					}
					break
				}
				if ok {
					return true, brokerID, nil
				}
			}
		}
	}

	if !hasModule {
		return false, "", errors.Errorf("broker %s has no instance with module (%s) loaded", brokerID, checkType)
	}
	if connErr != nil {
		return false, "", errors.Errorf("broker %s has no instance with adequate connectivity - %s", brokerID, connErr)
	}

	return false, "", errors.Errorf("broker %s has no viable instance", brokerID)
}

func brokerConnectionTest(ip string, port uint16, deadline time.Duration) (bool, error) {
	if ip == "" {
		return false, errors.New("invalid ip (empty)")
	}
	if port == 0 {
		return false, errors.New("invalid port (empty|0)")
	}
	if deadline == time.Duration(0) {
		return false, errors.New("invalid duration (0)")
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), deadline)
	if err != nil {
		return false, err
	}
	if err := conn.Close(); err != nil {
		return false, err
	}

	return true, nil
}
