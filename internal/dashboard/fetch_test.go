// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"testing"

	"github.com/spf13/viper"
)

func TestFetchByID(t *testing.T) {
	t.Log("Testing [fetch] FetchByID")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		desc   string
		cli    CircAPI
		id     string
		errMsg string
	}{
		{"invalid state", nil, "", "invalid state, nil client"},
		{"invalid (empty)", client, "", "invalid id (empty)"},
		{"invalid (foo)", client, "foo", "invalid dashboard id (foo)"},
		{"invalid (apierror)", client, "000", "fetch api: forced mock api call error"},
		{"valid", client, "123", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := FetchByID(test.cli, test.id)
		if test.errMsg == "" {
			if err != nil {
				t.Fatalf("unexpected error (%s)", err)
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			} else if err.Error() != test.errMsg {
				t.Fatalf("unexpected error (%s)", err)
			}
		}
	}
}

func TestFetchByTitle(t *testing.T) {
	t.Log("Testing [fetch] FetchByTitle")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		desc   string
		cli    CircAPI
		title  string
		errMsg string
	}{
		{"invalid state", nil, "", "invalid state, nil client"},
		{"invalid (empty)", client, "", "invalid title (empty)"},
		{"invalid (apierr)", client, "apierror", "search api: \"apierror\" (active:1)"},
		{"invalid (none found)", client, "none", "no dashboards found matching (none)"},
		{"invalid (multiple found)", client, "multi", "multiple dashboards matching (multi)"},
		{"valid", client, "valid", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := FetchByTitle(test.cli, test.title)
		if test.errMsg == "" {
			if err != nil {
				t.Fatalf("unexpected error (%s)", err)
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			} else if err.Error() != test.errMsg {
				t.Fatalf("unexpected error (%s)", err)
			}
		}
	}
}

func TestFetch(t *testing.T) {
	t.Log("Testing Fetch")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		desc   string
		cli    CircAPI
		id     string
		title  string
		errMsg string
	}{
		{"invalid state", nil, "", "", "invalid state, nil client"},
		{"no arguments", client, "", "", "missing required argument identifying which dashboard to fetch"},
		{"invalid id", client, "foo", "", "dashboard by id: invalid dashboard id (foo)"},
		{"valid id", client, "123", "", ""},
		{"invalid title", client, "", "none", "dashboard by title: no dashboards found matching (none)"},
		{"valid title", client, "", "foo", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := Fetch(test.cli, test.id, test.title)
		if test.errMsg == "" {
			if err != nil {
				t.Fatalf("unexpected error (%s)", err)
			}
		} else {
			if err == nil {
				t.Fatal("expected error")
			} else if err.Error() != test.errMsg {
				t.Fatalf("unexpected error (%s)", err)
			}
		}
	}
}
