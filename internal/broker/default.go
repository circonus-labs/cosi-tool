// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Default shows the default broker(s) (from COSI API and/or custom options file)
func Default(cosiURL string, w io.Writer, cosiBID uint, customOptionsFile string) error {
	logger := log.With().Str("cmd", "cosi broker default").Logger()

	if w == nil {
		return errors.New("invalid destination (nil)")
	}

	jsonID, trapID, err := getDefaultBrokers(cosiURL, customOptionsFile, cosiBID, logger)

	if err != nil {
		return errors.Wrap(err, "cosi broker default")
	}

	if jsonID > 0 && trapID > 0 {
		showDefault(w, jsonID, trapID)
		return nil
	}

	return errors.New("no default brokers defined")
}

func showDefault(w io.Writer, jsonID, trapID uint) {
	fmt.Fprintf(w, "JSON: %d\nTrap: %d\n", jsonID, trapID)
}

func getDefaultBrokers(cosiURL, optsFile string, cosiBID uint, logger zerolog.Logger) (uint, uint, error) {
	// explicit broker defined by user supplied (using cosi-install --broker option)
	if cosiBID > 0 {
		return cosiBID, cosiBID, nil
	}

	brokerID, custBroker, err := loadCustomBroker(optsFile, logger)
	if err != nil {
		return 0, 0, errors.Wrap(err, "custom options file")
	}

	// explicit broker specified in user supplied custom options file
	if brokerID > 0 {
		return brokerID, brokerID, nil
	}

	// explicit json and trap brokers specified in user supplied custom options file
	if custBroker.Trap > 0 && custBroker.JSON > 0 {
		return custBroker.JSON, custBroker.Trap, nil
	}

	cosiBroker, err := fetchCosiDefaults(cosiURL, logger)
	if err != nil {
		return 0, 0, errors.Wrap(err, "fetch cosi defaults")
	}

	var jsonID, trapID uint

	if custBroker.JSON != 0 {
		jsonID = custBroker.JSON
	} else if cosiBroker.JSON != 0 {
		jsonID = cosiBroker.JSON
	} else {
		if custBroker.Fallback != 0 {
			jsonID = custBroker.Fallback
		} else if cosiBroker.Fallback != 0 {
			jsonID = cosiBroker.Fallback
		} else {
			return 0, 0, errors.New("no json (or fallback) default brokers defined")
		}
	}

	if custBroker.Trap != 0 {
		trapID = custBroker.Trap
	} else if cosiBroker.Trap != 0 {
		trapID = cosiBroker.Trap
	} else {
		if custBroker.Fallback != 0 {
			trapID = custBroker.Fallback
		} else if cosiBroker.Fallback != 0 {
			trapID = cosiBroker.Fallback
		} else {
			return 0, 0, errors.New("no trap (or fallback) default brokers defined")
		}
	}

	if jsonID > 0 && trapID > 0 {
		return jsonID, trapID, nil
	}

	return 0, 0, errors.New("no default brokers defined")
}
