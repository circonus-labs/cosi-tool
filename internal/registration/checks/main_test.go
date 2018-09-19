// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	"io/ioutil"
	"testing"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func genMockCosiAPI() CosiAPI {
	return &CosiAPIMock{
		FetchTemplateFunc: func(id string) (*cosiapi.Template, error) {
			switch id {
			case "check-system":
				data, err := ioutil.ReadFile("testdata/template-check-system.toml")
				if err != nil {
					return nil, errors.Wrapf(err, "reading (%s) template", id)
				}
				var tmpl cosiapi.Template
				if err := toml.Unmarshal(data, &tmpl); err != nil {
					return nil, errors.Wrapf(err, "parsing (%s) template", id)
				}
				return &tmpl, nil
			case "check-group":
				data, err := ioutil.ReadFile("testdata/template-check-group.toml")
				if err != nil {
					return nil, errors.Wrapf(err, "reading (%s) template", id)
				}
				var tmpl cosiapi.Template
				if err := toml.Unmarshal(data, &tmpl); err != nil {
					return nil, errors.Wrapf(err, "parsing (%s) template", id)
				}
				return &tmpl, nil
			case "check-error":
				return nil, errors.Errorf("simulated api error response")
			}
			return nil, errors.Errorf("unknown template id (%s)", id)
		},
	}
}

func genMockCircAPI() CircAPI {
	return &CircAPIMock{
		CreateCheckBundleFunc: func(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
			if cfg.CID == "error" {
				return nil, errors.New("forced mock api error")
			}
			return cfg, nil
		},
		UpdateCheckBundleFunc: func(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
			panic("TODO: mock out the UpdateCheckBundle method")
		},
	}
}

func TestCreateCheck(t *testing.T) {
	t.Log("Testing createCheck")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name       string
		id         string
		obj        *circapi.CheckBundle
		shouldFail bool
		expected   string
	}{
		{"invalid (empty id)", "", nil, true, "invalid id (empty)"},
		{"invalid (nil obj)", "foo", nil, true, "invalid check bundle config (nil)"},
		{"invalid (apierr)", "apierr", &circapi.CheckBundle{CID: "error"}, true, "creating apierr: forced mock api error"},
		{"invalid (write)", "isdir", &circapi.CheckBundle{CID: "foo"}, true, "saving isdir registration: testdata/registration-isdir.json is a directory"},
		{"valid", "will_be_overwritten", &circapi.CheckBundle{CID: "/check_bundle/1234"}, false, ""},
	}

	c, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create checks object (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := c.createCheck(tst.id, tst.obj)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expected {
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

func TestGetCheckInfo(t *testing.T) {
	t.Log("Testing GetCheckInfo")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	c, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create checks object (%s)", err)
	}
	c.checkList = map[string]*circapi.CheckBundle{
		"check-system": {
			CID:        "/check_bundle/123",
			Checks:     []string{"/check/456"},
			CheckUUIDs: []string{"abcd-1234-efgh-5678"},
			Type:       "json:nad",
		},
		"check-group": {
			CID:        "/check_bundle/567",
			Checks:     []string{"/check/890"},
			CheckUUIDs: []string{"5678-efgh-1234-abcd"},
			Type:       "httptrap",
			Config: circapi.CheckBundleConfig{
				"submission_url": "http://127.0.0.1/foo",
			},
		},
		"check-bad": {
			CID:        "/check_bundle/bad",
			Checks:     []string{"/check/bad"},
			CheckUUIDs: []string{"bad-bad-bad"},
			Type:       "bad",
		},
	}

	{
		t.Log("empty")
		ci, err := c.GetCheckInfo("")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "invalid check id (empty)" {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci != nil {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
	{
		t.Log("invalid")
		ci, err := c.GetCheckInfo("invalid")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "check id not found (check-invalid)" {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci != nil {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
	{
		t.Log("bad")
		ci, err := c.GetCheckInfo("bad")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != `coverting check id to uint: strconv.ParseUint: parsing "bad": invalid syntax` {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci != nil {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
	{
		t.Log("system")
		ci, err := c.GetCheckInfo("system")
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci.CheckID != uint(456) {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
	{
		t.Log("check-system")
		ci, err := c.GetCheckInfo("check-system")
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci.CheckID != uint(456) {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
	{
		t.Log("check-group")
		ci, err := c.GetCheckInfo("check-group")
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
		if ci.CheckID != uint(890) {
			t.Fatalf("unexpected return (%#v)", ci)
		}
		if ci.SubmissionURL == "" {
			t.Fatalf("unexpected return (%#v)", ci)
		}
	}
}
