// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	"testing"

	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func TestCreateGroupCheck(t *testing.T) {
	t.Log("Testing createGroupCheck")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	c, err := New(&Options{
		Client: genMockCircAPI(),
		Config: &options.Options{
			Common: options.Common{
				Tags: []string{"cosi"},
			},
			Host: options.Host{
				Name: "foo",
			},
			Checks: options.Checks{
				Group: options.GroupCheck{
					Create:   true,
					ID:       "foobar",
					Tags:     []string{"bar", "baz"},
					BrokerID: "/broker/1",
				},
			},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create checks object (%s)", err)
	}

	if _, err := c.createGroupCheck(); err != nil {
		t.Fatalf("unexpected error (%s)", err)
	}
	// fmt.Printf("%#v\n", b)

}
