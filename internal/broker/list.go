// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"fmt"
	"io"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// DisplayList displays information about available broker(s)
func DisplayList(client API, w io.Writer) error {
	// logger := log.With().Str("cmd", "cosi broker list").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}
	if w == nil {
		return errors.New("invalid destination (nil)")
	}

	brokers, err := List(client)
	if err != nil {
		return errors.Wrap(err, "fetching brokers")
	}

	if len(*brokers) == 0 {
		return errors.New("no available brokers found")
	}

	format := "%5s %10s %-20s\n"
	fmt.Fprintf(w, format, "ID", "Type", "Name")

	for _, b := range *brokers {
		fmt.Fprintf(w, format, strings.Replace(b.CID, "/broker/", "", 1), b.Type, b.Name)
	}

	return nil
}

// List returns available brokers
func List(client API) (*[]api.Broker, error) {
	if client == nil {
		return nil, errors.New("invalid state, nil client")
	}

	brokers, err := client.FetchBrokers()
	if err != nil {
		return nil, errors.Wrap(err, "fetching brokers")
	}

	if len(*brokers) == 0 {
		return nil, errors.New("no brokers returned by API")
	}

	return brokers, nil
}
