// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

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
		{"invalid (foo)", client, "foo", "invalid worksheet id (foo)"},
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
		{"invalid (none found)", client, "none", "no worksheets found matching (none)"},
		{"invalid (multiple found)", client, "multi", "multiple worksheets matching (multi)"},
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
		{"no arguments", client, "", "", "missing required argument identifying which worksheet to fetch"},
		{"invalid id", client, "foo", "", "worksheet by id: invalid worksheet id (foo)"},
		{"valid id", client, "123", "", ""},
		{"invalid title", client, "", "none", "worksheet by title: no worksheets found matching (none)"},
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

	// regDir := "testdata/"
	//
	// t.Log("\tinvalid state")
	// if err := Fetch(nil, "", "", "", false); err == nil {
	// 	t.Fatal("expected error")
	// } else if err.Error() != "invalid state, nil client" {
	// 	t.Fatalf("expected different error, got (%v)", err)
	// }
	//
	// client := genMockClient()
	//
	// t.Log("\tno arguments")
	// if err := Fetch(client, "", "", "", false); err == nil {
	// 	t.Fatal("expected error")
	// } else if err.Error() != "missing required argument identifying which worksheet to fetch" {
	// 	t.Fatalf("expected different error, got (%v)", err)
	// }
	//
	// out := regDir + "will_be_overwritten.json"
	// force := true
	//
	// t.Log("\tid")
	// if err := Fetch(client, "123", "", out, force); err != nil {
	// 	t.Fatalf("expected NO error, got (%v)", err)
	// }
	// t.Log("\tinvalid id (foo)")
	// if err := Fetch(client, "foo", "", out, force); err == nil {
	// 	t.Fatal("expected error")
	// } else if err.Error() != "worksheet by id: invalid worksheet id (foo)" {
	// 	t.Fatalf("expected different error, got (%v)", err)
	// }
	//
	// t.Log("\ttitle")
	// if err := Fetch(client, "", "foo", out, force); err != nil {
	// 	t.Fatalf("expected NO error, got (%v)", err)
	// }
	// t.Log("\tvalid name, no worksheets")
	// if err := Fetch(client, "", "nodash", out, force); err == nil {
	// 	t.Fatal("expected error")
	// } else if err.Error() != "worksheet by title: no worksheets found matching (nodash)" {
	// 	t.Fatalf("expected different error, got (%v)", err)
	// }
}
