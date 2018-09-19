// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	FetchBroker(cid api.CIDType) (*api.Broker, error)
	FetchBrokers() (*[]api.Broker, error)
}
