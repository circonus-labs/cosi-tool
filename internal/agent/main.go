// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/pkg/errors"
)

// Client defines the circonus-agent api client configuration
type Client struct {
	apiClient AgentAPI
}

// New creates a new circonus-agent api client
func New(agentURL string) (*Client, error) {
	if agentURL == "" {
		return nil, errors.New("invalid agent URL (empty)")
	}

	c, err := agentapi.New(agentURL)
	if err != nil {
		return nil, err
	}

	return &Client{apiClient: c}, nil
}

// ActivePluginList returns a list of active plugin IDs (e.g. cpu, vm, if, disk, etc.)
// or an error based on the inventory returned from the running agent. The plugin ID
// may be used to build a list of graph templates (e.g. graph-cpu, etc.).
func (c *Client) ActivePluginList() (*[]string, error) {

	inventory, err := c.apiClient.Inventory()
	if err != nil {
		return nil, errors.Wrap(err, "fetching inventory from agent")
	}

	if len(*inventory) == 0 {
		return nil, errors.New("no plugins are active, 0 returned from agent")
	}

	// NOTE: this removes duplicates (if any plugins have instances, only want the base plugin id)
	plist := make(map[string]bool)
	for _, plugin := range *inventory {
		plist[plugin.Name] = true
	}

	plugins := make([]string, len(plist))
	i := 0
	for pid := range plist {
		plugins[i] = pid
		i++
	}

	return &plugins, nil
}

// AvailableMetrics returns a list of metrics available from the agent
func (c *Client) AvailableMetrics(pluginID string) (*agentapi.Metrics, error) {
	return c.apiClient.Metrics(pluginID)
}
