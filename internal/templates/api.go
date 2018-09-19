// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package templates

import "github.com/circonus-labs/cosi-server/api"

//go:generate moq -out api_test.go . API

// API interface abstraction of cosi server api (for mocking)
type API interface {
	FetchTemplate(id string) (*api.Template, error)
}
