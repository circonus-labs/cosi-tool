// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	agentapi "github.com/circonus-labs/circonus-agent/api"
)

//go:generate moq -out api_agent_test.go . AgentAPI

// AgentAPI interface abstraction of circonus api (for mocking)
type AgentAPI interface {
	Inventory() (*agentapi.Inventory, error)
	Metrics(pluginID string) (*agentapi.Metrics, error)
}
