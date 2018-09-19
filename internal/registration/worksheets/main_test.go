// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheets

import (
	"errors"
	"path"
	"testing"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func genMockCircAPI() CircAPI {
	return &CircAPIMock{
		CreateWorksheetFunc: func(cfg *circapi.Worksheet) (*circapi.Worksheet, error) {
			if cfg.CID == "error" {
				return nil, errors.New("forced mock api error")
			}
			return cfg, nil
		},
	}
}

func TestNew(t *testing.T) {
	t.Log("Testing New")

	client := genMockCircAPI()
	opts := options.Options{}
	regDir := "testdata"
	tmpl := templates.Templates{}

	tests := []struct {
		name        string
		client      CircAPI
		config      *options.Options
		regDir      string
		templates   *templates.Templates
		shouldFail  bool
		expectedErr string
	}{
		{"invalid (client)", nil, nil, "", nil, true, "invalid client (nil)"},
		{"invalid (options)", client, nil, "", nil, true, "invalid options config (nil)"},
		{"invalid (regdir)", client, &opts, "", nil, true, "invalid reg dir (empty)"},
		{"invalid (templates)", client, &opts, regDir, nil, true, "invalid templates (nil)"},
		{"invalid (missing regdir)", client, &opts, path.Join(regDir, "missing"), &tmpl, true, "finding worksheet registrations: reading registration directory: open testdata/missing: no such file or directory"},
		{"valid", client, &opts, regDir, &tmpl, false, ""},
	}

	if _, err := New(nil); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid options (nil)" {
		t.Fatalf("unexpected error (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := New(&Options{tst.client, tst.config, tst.regDir, tst.templates})
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

func TestRegister(t *testing.T) {
	t.Log("Testing Register")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	w, err := New(&Options{
		Client: genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{
				Name: "foo",
			},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	empty := map[string]bool{}
	nosheets := map[string]bool{"graph-vm": true, "dashboard-system": true}
	list := map[string]bool{"worksheet-test": true, "worksheet-disabled": false, "dashboard-system": true}

	{
		t.Log("empty")
		err := w.Register(empty)
		if err == nil {
			t.Fatal("epxected error")
		}
		if err.Error() != "invalid list (empty)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	{
		t.Log("nosheets")
		err := w.Register(nosheets)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	{
		t.Log("list")
		err := w.Register(list)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
}

func TestCreate(t *testing.T) {
	t.Log("Testing create")

	tests := []struct {
		name        string
		id          string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", true, "invalid id (empty)"},
		{"valid", "worksheet-test", false, ""},
	}

	w, err := New(&Options{
		Client: genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{
				Name: "foo",
			},
			Common: options.Common{
				Notes: "foobar",
				Tags:  []string{"foo:bar", "baz:qux"},
			},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}
	// reset regFiles (ignore the test reg file)
	w.regFiles = &[]string{}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			// t.Parallel() - creates are serial, not concurrent
			err := w.create(tst.id)
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

func TestCheckForRegistration(t *testing.T) {
	t.Log("Testing checkForRegistration")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	w, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	tests := []struct {
		name        string
		sheetID     string
		shouldFind  bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", false, true, "invalid id (empty)"},
		{"missing", "worksheet-missing", false, false, ""},
		{"error (parsing)", "worksheet-error", false, true, "loading registration-worksheet-error: parsing registration (testdata/registration-worksheet-error.json): unexpected end of JSON input"},
		{"valid", "worksheet-valid", true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			loaded, err := w.checkForRegistration(tst.sheetID)
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
			if tst.shouldFind {
				if !loaded {
					t.Fatal("expected found+loaded")
				}
			} else {
				if loaded {
					t.Fatal("expected not found/loaded")
				}
			}
		})
	}
}

func TestParseTemplateConfig(t *testing.T) {
	t.Log("Testing parseTemplateConfig")

	tvars := struct {
		HostName string
	}{"foo"}

	tests := []struct {
		name        string
		sheetID     string
		template    string
		vars        interface{}
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", "", nil, true, "invalid id (empty)"},
		{"invalid config (empty)", "foo", "", nil, true, "invalid template config (empty)"},
		{"invalid vars (nil)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid config (parse)", "foo", "baz{{}}", tvars, true, "parsing template: template: foo:1: missing value for command"},
		{"invalid config (exec)", "foo", "baz{{.BadVarName}}", tvars, true, "executing template: template: foo:1:5: executing \"foo\" at <.BadVarName>: can't evaluate field BadVarName in type struct { HostName string }"},
		{"invalid json (post exec)", "foo", "{", tvars, true, "parsing expanded template result: unexpected end of JSON input"},
		{"valid w/HostName", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
		{"valid w/all", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			v, err := parseTemplateConfig(tst.sheetID, tst.template, tst.vars)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				} else if v == nil {
					t.Fatal("expected valid return (not nil)")
				}
			}
		})
	}

}
