// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"github.com/circonus-labs/circonus-agent/api"
)

//go:generate moq -out api_agent_test.go . API

// API interface abstraction of circonus api (for mocking)
type API interface {
	Inventory() (*api.Inventory, error)
	Metrics(pluginID string) (*api.Metrics, error)
}
