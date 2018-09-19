// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestParseDashboardTemplate(t *testing.T) {
	t.Log("Testing parseDashboardTemplate")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tvars := struct {
		HostName  string
		CheckUUID string
	}{"foo", "bar"}
	tests := []struct {
		name        string
		dashID      string
		template    string
		vars        interface{}
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", "", nil, true, "invalid dashboard id (empty)"},
		{"invalid config (empty)", "foo", "", nil, true, "invalid template config (empty)"},
		{"invalid vars (nil)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid config (parse)", "foo", "baz{{}}", tvars, true, "parsing template: template: foo:1: missing value for command"},
		{"invalid config (exec)", "foo", "baz{{.BadVarName}}", tvars, true, "executing template: template: foo:1:5: executing \"foo\" at <.BadVarName>: can't evaluate field BadVarName in type struct { HostName string; CheckUUID string }"},
		{"invalid json (post exec)", "foo", "{", tvars, true, "parsing expanded template result: unexpected end of JSON input"},
		{"valid w/HostName", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
		{"valid w/CheckUUID", "foo", `{"item":"{{.CheckUUID}}"}`, tvars, false, ""},
		{"valid w/all", "foo", `{"host":"{{.HostName}}","item":"{{.CheckUUID}}"}`, tvars, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			v, err := parseDashboardTemplate(tst.dashID, tst.template, tst.vars)
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

func TestParseWidgetTemplate(t *testing.T) {
	t.Log("Testing parseWidgetTemplate")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tvars := struct {
		HostName  string
		CheckUUID string
		GraphUUID string
	}{"foo", "bar", "baz"}
	tests := []struct {
		name        string
		dashID      string
		template    string
		vars        interface{}
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", "", nil, true, "invalid dashboard id (empty)"},
		{"invalid config (empty)", "foo", "", nil, true, "invalid template config (empty)"},
		{"invalid vars (nil)", "foo", "bar", nil, true, "invalid template vars (nil)"},
		{"invalid config (parse)", "foo", "baz{{}}", tvars, true, "parsing template: template: foo:1: missing value for command"},
		{"invalid config (exec)", "foo", "baz{{.BadVarName}}", tvars, true, "executing template: template: foo:1:5: executing \"foo\" at <.BadVarName>: can't evaluate field BadVarName in type struct { HostName string; CheckUUID string; GraphUUID string }"},
		{"invalid json (post exec)", "foo", "{", tvars, true, "parsing expanded template result: unexpected end of JSON input"},
		{"valid w/HostName", "foo", `{"host":"{{.HostName}}"}`, tvars, false, ""},
		{"valid w/CheckUUID", "foo", `{"check_uuid":"{{.CheckUUID}}"}`, tvars, false, ""},
		{"valid w/GraphUUID", "foo", `{"graph_uuid":"{{.GraphUUID}}"}`, tvars, false, ""},
		{"valid w/all", "foo", `{"host":"{{.HostName}}","check_uuid":"{{.CheckUUID}}","graph_uuid":"{{.GraphUUID}}"}`, tvars, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			v, err := parseWidgetTemplate(tst.dashID, tst.template, tst.vars)
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
