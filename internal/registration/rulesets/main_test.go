// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package rulesets

import (
	"errors"
	"os"
	"testing"

	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/rs/zerolog"
)

func genMockCircAPI() CircAPI {
	return &CircAPIMock{
		CreateRuleSetFunc: func(cfg *circapi.RuleSet) (*circapi.RuleSet, error) {
			if cfg.CheckCID == "error" {
				return nil, errors.New("forced mock api error")
			}
			return cfg, nil
		},
	}
}

func TestNew(t *testing.T) {
	t.Log("Testing New")

	client := genMockCircAPI()

	tests := []struct {
		name        string
		cfg         *Options
		shouldFail  bool
		expectedErr string
	}{
		{"invalid config (nil)", nil, true, "invalid cfg (nil)"},
		{"invalid check info (nil)", &Options{}, true, "invalid check info (nil)"},
		{"invalid client (nil)", &Options{
			CheckInfo: &checks.CheckInfo{},
		}, true, "invalid client (nil)"},
		{"invalid options config (nil)", &Options{
			CheckInfo: &checks.CheckInfo{},
			Client:    client,
		}, true, "invalid options config (nil)"},
		{"invalid registration dir (empty)", &Options{
			CheckInfo: &checks.CheckInfo{},
			Client:    client,
			Config:    &options.Options{},
		}, true, "invalid registration directory (empty)"},
		{"invalid ruleset dir (empty)", &Options{
			CheckInfo: &checks.CheckInfo{},
			Client:    client,
			Config:    &options.Options{},
			RegDir:    "testdata",
		}, true, "invalid ruleset directory (empty)"},
		{"valid", &Options{
			CheckInfo:  &checks.CheckInfo{},
			Client:     client,
			Config:     &options.Options{},
			RegDir:     "testdata",
			RulesetDir: "testdata",
		}, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tst.cfg)
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

	client := genMockCircAPI()
	options := options.Options{Common: options.Common{Tags: []string{"c1:v1"}}}
	ci := checks.CheckInfo{CheckCID: "/check/1234"}
	ciErr := checks.CheckInfo{CheckCID: "error"}

	tests := []struct {
		name        string
		cfg         *Config
		shouldFail  bool
		expectedErr string
	}{
		{"missing rulesetdir", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ci,
			RegDir:     "testdata",
			RulesetDir: "testdata/missing"}, false, ""},
		{"empty rulesetdir", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ci,
			RegDir:     "testdata",
			RulesetDir: "testdata/empty"}, false, ""},
		{"invalid ruleset config", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ci,
			RegDir:     "testdata",
			RulesetDir: "testdata/invalid"}, true, "unexpected end of JSON input"},
		{"api error", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ciErr,
			RegDir:     "testdata",
			RulesetDir: "testdata/apierror"}, true, "forced mock api error"},
		{"valid", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ci,
			RegDir:     "testdata",
			RulesetDir: "testdata"}, false, ""},
		{"valid (exists)", &Config{
			Client:     client,
			Config:     &options,
			CheckInfo:  &ci,
			RegDir:     "testdata",
			RulesetDir: "testdata"}, false, ""},
	}

	// remove the test registration so it will be created, then not created because it exists
	os.Remove("testdata/registration-ruleset-valid-ignore.json")

	for _, test := range tests {
		tst := test

		t.Run(tst.name, func(t *testing.T) {
			r, err := New(tst.cfg)
			if err != nil {
				t.Fatalf("unexpected error (%s)", err)
			}
			rerr := r.Register()
			if tst.shouldFail {
				if rerr == nil {
					t.Fatal("expected error")
				} else if rerr.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", rerr)
				}
			} else {
				if rerr != nil {
					t.Fatalf("unexpected error (%s)", rerr)
				}
			}
		})
	}
}
