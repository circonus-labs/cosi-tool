// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestParseGraphTemplate(t *testing.T) {
	t.Log("Testing parseGraphTemplate")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tvars := struct {
		HostName string
		Item     string
	}{"foo", "bar"}
	tests := []struct {
		name        string
		graphID     string
		template    string
		vars        interface{}
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", "", nil, true, "invalid graph id (empty)"},
		{"invalid config (empty)", "foo", "", nil, true, "invalid template config (empty)"},
		{"invalid vars (nil)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid config (parse)", "foo", "baz{{}}", tvars, true, "parsing template: template: foo:1: missing value for command"},
		{"invalid config (exec)", "foo", "baz{{.BadVarName}}", tvars, true, "executing template: template: foo:1:5: executing \"foo\" at <.BadVarName>: can't evaluate field BadVarName in type struct { HostName string; Item string }"},
		{"invalid json (post exec)", "foo", "{", tvars, true, "parsing expanded template result: unexpected end of JSON input"},
		{"valid w/HostName", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
		{"valid w/Item", "foo", `{"item":"{{.Item}}"}`, tvars, false, ""},
		{"valid w/all", "foo", `{"host":"{{.HostName}}","item":"{{.Item}}"}`, tvars, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			v, err := parseGraphTemplate(tst.graphID, tst.template, tst.vars)
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

func TestParseDatapointTemplate(t *testing.T) {
	t.Log("Testing parseDatapointTemplate")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tvars := struct {
		HostName   string
		Item       string
		MetricName string
	}{"foo", "bar", "baz"}
	tests := []struct {
		name        string
		graphID     string
		template    string
		vars        interface{}
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", "", nil, true, "invalid graph id (empty)"},
		{"invalid config (empty)", "foo", "", nil, true, "invalid template config (empty)"},
		{"invalid vars (nil)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid config (parse)", "foo", "baz{{}}", tvars, true, "parsing template: template: foo:1: missing value for command"},
		{"invalid config (exec)", "foo", "baz{{.BadVarName}}", tvars, true, "executing template: template: foo:1:5: executing \"foo\" at <.BadVarName>: can't evaluate field BadVarName in type struct { HostName string; Item string; MetricName string }"},
		{"invalid json (post exec)", "foo", "{", tvars, true, "parsing expanded template result: unexpected end of JSON input"},
		{"valid w/HostName", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
		{"valid w/Item", "foo", `{"item":"{{.Item}}"}`, tvars, false, ""},
		{"valid w/MetricName", "foo", `{"metric":"{{.MetricName}}"}`, tvars, false, ""},
		{"valid w/all", "foo", `{"host":"{{.HostName}}","item":"{{.Item}}","metric":"{{.MetricName}}"}`, tvars, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			v, err := parseDatapointTemplate(tst.graphID, tst.template, tst.vars)
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
