// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package templates

import (
	"errors"
	"strings"
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-server/api"
)

var (
	checkTemplate     *api.Template
	dashboardTemplate *api.Template
	graphTemplate     *api.Template
	worksheetTemplate *api.Template
)

func genMockClient() *APIMock {
	return &APIMock{
		FetchTemplateFunc: func(id string) (*api.Template, error) {
			if strings.Contains(id, "error") {
				return nil, errors.New("forced mock api call error")
			} else if strings.HasPrefix(id, "check-") {
				return checkTemplate, nil
			} else if strings.HasPrefix(id, "dashboard-") {
				return dashboardTemplate, nil
			} else if strings.HasPrefix(id, "graph-") {
				return graphTemplate, nil
			} else if strings.HasPrefix(id, "worksheet-") {
				return worksheetTemplate, nil
			}
			return nil, errors.New("unknown template id")
		},
	}
}

func TestFetch(t *testing.T) {
	t.Log("Testing Fetch")

	tests := []struct {
		name       string
		id         string
		shouldFail bool
		expected   string
	}{
		{"invalid (empty)", "", true, "invalid id (empty)"},
		{"invalid (invalid)", "test-invalid", true, "unknown template id"},
		{"invalid (apierror)", "test-error", true, "forced mock api call error"},
		{"valid", "check-test", false, ""},
	}

	tmpl := Templates{client: genMockClient()}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := tmpl.Fetch(tst.id)
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

func TestFetchAll(t *testing.T) {
	t.Log("Testing FetchAll")

	tests := []struct {
		name       string
		ids        []string
		shouldFail bool
		expected   string
	}{
		{"invalid (empty)", []string{}, true, "invalid template list (empty)"},
		{"invalid (invalid)", []string{"test-invalid"}, true, "unknown template id"},
		{"invalid (apierror)", []string{"test-error"}, true, "forced mock api call error"},
		{"valid", []string{"check-test"}, false, ""},
	}

	tmpl := Templates{client: genMockClient()}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			tlist, err := tmpl.FetchAll(tst.ids)
			if tst.shouldFail {
				if err == nil {
					if tlist == nil {
						t.Fatal("expected error")
					} else {
						if (*tlist)[0].Err.Error() != tst.expected {
							t.Fatalf("unexpected error (%s)", err)
						}
					}
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

func TestLoad(t *testing.T) {
	t.Log("Testing Load")

	dir := "testdata"
	tests := []struct {
		name       string
		dir        string
		id         string
		shouldFind bool
		shouldFail bool
		expected   string
	}{
		{"invalid (dir empty)", "", "", false, true, "invalid directory (empty)"},
		{"invalid (id empty)", dir, "", false, true, "invalid id (empty)"},
		{"invalid (readerr)", dir, "test-readerr", true, true, "reading template: read testdata/template-test-readerr.toml: is a directory"},
		{"invalid (syntax)", dir, "test-invalid", true, true, "parsing template: (1, 6): was expecting token =, but got keys cannot contain : character instead"},
		{"valid", dir, "test-valid", true, false, ""},
	}

	tmpl := Templates{client: genMockClient()}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, found, err := tmpl.Load(tst.dir, tst.id)
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
			if tst.shouldFind {
				if !found {
					t.Fatal("expected template to be found")
				}
			} else {
				if found {
					t.Fatal("expected template to NOT be found")
				}
			}
		})
	}
}

func TestDefaultTemplateList(t *testing.T) {
	t.Log("Testing DefaultTemplateList")

	tests := []struct {
		name        string
		metrics     *agentapi.Metrics
		shouldFail  bool
		expectedErr string
	}{
		{"invalid", nil, true, "invalid metric list (nil)"},
		{"empty", &agentapi.Metrics{}, true, "invalid metric list (empty)"},
		{"valid (bare)", &agentapi.Metrics{"foo": agentapi.Metric{Type: "n", Value: 1.2}}, false, ""},
		{"valid (w/delim)", &agentapi.Metrics{"bar`baz": agentapi.Metric{Type: "n", Value: 1.2}}, false, ""},
		{"valid (dedupe)", &agentapi.Metrics{"baz`bar": agentapi.Metric{Type: "n", Value: 1.2}, "baz`foo": agentapi.Metric{Type: "n", Value: 2.1}}, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			_, err := DefaultTemplateList(tst.metrics)
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
