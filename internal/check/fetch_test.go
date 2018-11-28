// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

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
		{"invalid (foo)", client, "foo", "invalid check bundle id (foo)"},
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

func TestFetchByType(t *testing.T) {
	t.Log("Testing [fetch] FetchByType")

	viper.Reset()
	client := genMockClient()
	regDir := "testdata/"

	tests := []struct {
		desc   string
		cli    CircAPI
		dir    string
		ctype  string
		errMsg string
	}{
		{"invalid state", nil, "", "", "invalid state, nil client"},
		{"invalid (empty dir)", client, "", "", "invalid registration directory (empty)"},
		{"invalid (empty type)", client, regDir, "", "invalid check type (empty)"},
		{"invalid (type)", client, regDir, "foo", "invalid check type (foo)"},
		{"invalid (missing)", client, regDir, "group", "loading check type: open testdata/registration-check-group.json: no such file or directory"},
		{"valid", client, regDir, "system", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := FetchByType(test.cli, test.dir, test.ctype)
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

func TestFetchByName(t *testing.T) {
	t.Log("Testing [fetch] FetchByName")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		desc   string
		cli    CircAPI
		name   string
		errMsg string
	}{
		{"invalid state", nil, "", "invalid state, nil client"},
		{"invalid (empty)", client, "", "invalid display name (empty)"},
		{"invalid (apierr)", client, "apierror", "search api: \"apierror\" (active:1)"},
		{"invalid (none found)", client, "none", "no checks found matching (none)"},
		{"invalid (multiple found)", client, "multi", "multiple checks matching (multi)"},
		{"valid", client, "valid", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := FetchByName(test.cli, test.name)
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

func TestFetchByTarget(t *testing.T) {
	t.Log("Testing [fetch] FetchByTarget")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		desc   string
		cli    CircAPI
		target string
		errMsg string
	}{
		{"invalid state", nil, "", "invalid state, nil client"},
		{"invalid (empty)", client, "", "invalid target (empty)"},
		{"invalid (apierr)", client, "apierror", "search api: (host:\"apierror\")(active:1)"},
		{"invalid (none found)", client, "none", "no checks found matching (none)"},
		{"invalid (multiple found)", client, "multi", "multiple checks matching (multi)"},
		{"valid", client, "valid", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := FetchByTarget(test.cli, test.target)
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
	regDir := "testdata/"

	tests := []struct {
		desc   string
		cli    CircAPI
		id     string
		ctype  string
		name   string
		target string
		errMsg string
	}{
		{"invalid state", nil, "", "", "", "", "invalid state, nil client"},
		{"no arguments", client, "", "", "", "", "missing required argument identifying which check to fetch"},
		{"invalid id", client, "foo", "", "", "", "check by id: invalid check bundle id (foo)"},
		{"valid id", client, "123", "", "", "", ""},
		{"invalid type", client, "", "none", "", "", "check by type: invalid check type (none)"},
		{"valid type", client, "", "system", "", "", ""},
		{"invalid name", client, "", "", "none", "", "check by name: no checks found matching (none)"},
		{"valid name", client, "", "", "foo", "", ""},
		{"invalid target", client, "", "", "", "none", "check by target: no checks found matching (none)"},
		{"valid target", client, "", "", "", "foo", ""},
	}

	for _, test := range tests {
		t.Log("\t", test.desc)
		_, err := Fetch(test.cli, regDir, test.id, test.ctype, test.name, test.target)
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
