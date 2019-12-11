// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package broker provides the various supporting functions for the
// `cosi broker *` commands
package broker

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	// KeyCID is the broker id
	KeyCID = "check.id"
	// CIDDefault is the default value for the id show option
	CIDDefault = ""
)

type customBroker struct {
	ID           uint   `json:"id"`            // >0 (there is no broker id 0)
	TrapList     []uint `json:"httptrap_list"` // list of valid broker ids
	TrapIdx      int    `json:"httptrap_idx"`  // index into list or -1 for random
	JSONList     []uint `json:"json_list"`     // list of valid broker ids
	JSONIdx      int    `json:"json_idx"`      // index into list or -1 for random
	FallbackList []uint `json:"fallback_list"` // list of valid broker ids
	FallbackIdx  int    `json:"fallback_idx"`  // index into list or -1 for random
}

type customOptions struct {
	Broker customBroker `json:"broker"`
}

type defaultBroker struct {
	Trap     uint `json:"httptrap,omitempty"`
	JSON     uint `json:"json,omitempty"`
	Fallback uint `json:"fallback,omitempty"`
}

func init() {
	// for random broker selection
	rand.Seed(time.Now().UnixNano())
}

func fetchCosiDefaults(cosiURL string, logger zerolog.Logger) (*defaultBroker, error) {
	if cosiURL == "" {
		return nil, errors.New("invalid cosi url (empty)")
	}

	var cosiBroker defaultBroker

	reqURL := cosiURL
	if !strings.HasSuffix(reqURL, "/") {
		reqURL += "/"
	}
	reqURL += "brokers/"
	r, err := http.Get(reqURL) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "http request")
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		xtra, _ := ioutil.ReadAll(r.Body)
		logger.Debug().Int("code", r.StatusCode).Str("body", string(xtra)).Msg("unexpected response from COSI API")
		return nil, errors.Errorf("unexpected return code from cosi api (%d - %s)", r.StatusCode, r.Status)
	}

	if err := json.NewDecoder(r.Body).Decode(&cosiBroker); err != nil {
		return nil, errors.Wrap(err, "parsing cosi api json")
	}

	return &cosiBroker, nil
}

func loadCustomBroker(file string, logger zerolog.Logger) (uint, *defaultBroker, error) {
	if file == "" {
		return 0, &defaultBroker{}, nil // this is ok, custom options are ... "optional"
	}

	var custBroker defaultBroker

	f, err := os.Open(file)
	if err != nil {
		return 0, nil, errors.Wrap(err, "custom options file")
	}

	var cb customOptions
	if err := json.NewDecoder(f).Decode(&cb); err != nil {
		return 0, nil, errors.Wrap(err, "parsing custom options file json")
	}

	if cb.Broker.ID != 0 {
		return cb.Broker.ID, nil, nil
	}

	{
		// get setting for each type

		var id uint
		var err error

		id, err = getBrokerFromList(cb.Broker.FallbackList, cb.Broker.FallbackIdx)
		if err != nil {
			logger.Warn().Err(err).Msg("no fallback setting, ignoring")
		} else {
			custBroker.Fallback = id
		}

		id, err = getBrokerFromList(cb.Broker.TrapList, cb.Broker.TrapIdx)
		if err != nil {
			logger.Warn().Err(err).Msg("no trap setting, ignoring")
		} else {
			custBroker.Trap = id
		}

		id, err = getBrokerFromList(cb.Broker.JSONList, cb.Broker.JSONIdx)
		if err != nil {
			logger.Warn().Err(err).Msg("no json setting, ignoring")
		} else {
			custBroker.JSON = id
		}
	}

	return 0, &custBroker, nil
}

func getBrokerFromList(list []uint, idx int) (uint, error) {
	listLen := len(list)

	if listLen == 0 {
		return 0, errors.New("empty list")
	}

	if idx == -1 {
		idx = rand.Intn(listLen)
		return list[idx], nil
	}

	if idx < 0 || idx >= listLen {
		return 0, errors.Errorf("index (%d) out of range for list size (%d)", idx, listLen)
	}

	return list[idx], nil
}
