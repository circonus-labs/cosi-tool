// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/graphs"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func TestCreate(t *testing.T) {
	t.Log("Testing create")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name        string
		id          string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", true, "invalid id (empty)"},
		{"valid", "dashboard-test", false, ""},
	}

	d, err := New(&Options{
		Client: genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		GraphInfo: &map[string]graphs.GraphInfo{"graph-test": {CID: "/graphs/abcd-efgh-0123-4567", UUID: "abcd-efgh-0123-4567"}},
		Metrics:   &agentapi.Metrics{"test": {}},
	})
	if err != nil {
		t.Fatalf("unable to create dashboards object (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			err := d.create(tst.id)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
		})
	}
}

func TestLoadMeta(t *testing.T) {
	t.Log("Testing loadMeta")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name        string
		id          string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", true, "invalid id (empty)"},
		{"error (parsing)", "meta-error", true, "unexpected end of JSON input"},
		{"invalid (value not string)", "meta-invalid", true, "json: cannot unmarshal number into Go value of type string"},
		{"missing", "meta-missing", false, ""},
		{"valid", "meta-valid", false, ""},
	}

	d, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		GraphInfo: &map[string]graphs.GraphInfo{},
		Metrics:   &agentapi.Metrics{"test": {}},
	})
	if err != nil {
		t.Fatalf("unable to create dashboards object (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := d.loadMeta(tst.id)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}
		})
	}
}
