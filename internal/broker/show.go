// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/pkg/errors"
)

// Show displays information about a specific broker
func Show(client API, w io.Writer, id string) error {
	// logger := log.With().Str("cmd", "cosi broker show").Logger()

	if client == nil {
		return errors.New("invalid state, nil client")
	}
	if w == nil {
		return errors.New("invalid destination (nil)")
	}
	if id == "" {
		return errors.New("invalid broker id (empty)")
	}

	cid := id
	if !strings.HasPrefix(cid, "/broker/") {
		cid = "/broker/" + id
	}
	if ok, err := regexp.MatchString(`^/broker/[0-9]+`, cid); err != nil {
		return errors.Wrap(err, "compile broker id regexp")
	} else if !ok {
		return errors.Errorf("invalid broker id (%s)", id)
	}
	b, err := client.FetchBroker(api.CIDType(&cid))
	if err != nil {
		return errors.Wrap(err, "fetch api")
	}

	var modules []string
	maxWidth := uint(0)
	active := false

	for _, d := range b.Details {
		if d.Status != "active" {
			continue
		}
		active = true
		for _, m := range d.Modules {
			if m == "selfcheck" || strings.HasPrefix(m, "hidden:") {
				continue
			}
			maxWidth = uint(math.Max(float64(maxWidth), float64(len(m))))
			modules = append(modules, m)
		}
	}

	if active {
		numModules := len(modules)
		maxWidth += 2
		numColumns := int(math.Floor(70 / float64(maxWidth)))
		format := fmt.Sprintf("%%-%ds", maxWidth)
		spacer := strings.Repeat(" ", 8)

		fmt.Fprintf(w, "ID    : %s\n", strings.Replace(b.CID, "/broker/", "", 1))
		fmt.Fprintf(w, "Name  : %s\n", b.Name)
		fmt.Fprintf(w, "Type  : %s\n", b.Type)
		fmt.Fprintf(w, "Checks: %d types supported\n", numModules)

		for i := 0; i < numModules; i += numColumns {
			line := ""
			for j := 0; j < numColumns; j++ {
				m := ""
				if i+j < numModules {
					m = modules[i+j]
				}
				line += fmt.Sprintf(format, m)
			}
			fmt.Fprintf(w, "%s%s\n", spacer, line)
		}

		return nil
	}

	return errors.Errorf("Broker %s is not active", id)
}
