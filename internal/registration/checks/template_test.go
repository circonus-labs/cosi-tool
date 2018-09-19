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

func TestParseTemplateConfig(t *testing.T) {
	t.Log("Testing parseTemplateConfig")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	type templateVars struct {
		GroupID    string
		HostName   string
		HostTarget string
	}
	tvars := templateVars{"foo", "bar", "baz"}
	tests := []struct {
		name         string
		checkType    string
		checkName    string
		templateVars interface{}
		shouldFail   bool
		expectedErr  string
	}{
		{"invalid (empty type)", "", "", nil, true, "invalid config type (empty)"},
		{"invalid (empty name)", "foo", "", nil, true, "invalid config name (empty)"},
		{"invalid (nil vars)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid (unknown id)", "foo", "bar", &tvars, true, "unknown template id (foo-bar)"},
		{"invalid (api error)", "check", "error", &tvars, true, "simulated api error response"},
		{"valid (system)", "check", "system", &tvars, false, ""},
		{"valid (group)", "check", "group", &tvars, false, ""},
	}

	tmpl, err := templates.New(genMockCosiAPI())
	if err != nil {
		t.Fatalf("unexpected error (%s)", err)
	}
	c, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: tmpl,
	})
	if err != nil {
		t.Fatalf("unable to create checks object (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			b, err := c.parseTemplateConfig(tst.checkType, tst.checkName, tst.templateVars)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
				if b != nil {
					t.Fatal("expected nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
				if b == nil {
					t.Fatal("expected NOT nil")
				}
				if tst.checkName == "system" {
					if b.DisplayName != tst.templateVars.(*templateVars).HostName+" cosi/system" {
						t.Fatalf("unexpected DisplayName (%s)", b.DisplayName)
					}
				} else if tst.checkName == "group" {
					if b.DisplayName != tst.templateVars.(*templateVars).GroupID+" cosi/group" {
						t.Fatalf("unexpected DisplayName (%s)", b.DisplayName)
					}
				}
			}
		})
	}
}
