// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import (
	"testing"

	"github.com/spf13/viper"
)

// func TestFetchByID(t *testing.T) {
// 	t.Log("Testing [fetch] FetchByID")
//
// 	viper.Reset()
// 	client := genMockClient()
//
// 	tests := []struct {
// 		desc   string
// 		cli    API
// 		id     string
// 		errMsg string
// 	}{
// 		{"invalid state", nil, "", "invalid state, nil client"},
// 		{"invalid (empty)", client, "", "invalid id (empty)"},
// 		{"invalid (foo)", client, "foo", "invalid graph id (foo)"},
// 		{"invalid (apierror)", client, "000", "fetch api: forced mock api call error"},
// 		{"valid", client, "123", ""},
// 	}
//
// 	for _, test := range tests {
// 		t.Log("\t", test.desc)
// 		_, err := FetchByID(test.cli, test.id)
// 		if test.errMsg == "" {
// 			if err != nil {
// 				t.Fatalf("unexpected error (%s)", err)
// 			}
// 		} else {
// 			if err == nil {
// 				t.Fatal("expected error")
// 			} else if err.Error() != test.errMsg {
// 				t.Fatalf("unexpected error (%s)", err)
// 			}
// 		}
// 	}
// }
//
// func TestFetchByTitle(t *testing.T) {
// 	t.Log("Testing [fetch] FetchByTitle")
//
// 	viper.Reset()
// 	client := genMockClient()
//
// 	tests := []struct {
// 		desc   string
// 		cli    API
// 		title  string
// 		errMsg string
// 	}{
// 		{"invalid state", nil, "", "invalid state, nil client"},
// 		{"invalid (empty)", client, "", "invalid title (empty)"},
// 		{"invalid (apierr)", client, "apierror", "search api: \"apierror\" (active:1)"},
// 		{"invalid (none found)", client, "none", "no graphs found matching (none)"},
// 		{"invalid (multiple found)", client, "multi", "multiple graphs matching (multi)"},
// 		{"valid", client, "valid", ""},
// 	}
//
// 	for _, test := range tests {
// 		t.Log("\t", test.desc)
// 		_, err := FetchByTitle(test.cli, test.title)
// 		if test.errMsg == "" {
// 			if err != nil {
// 				t.Fatalf("unexpected error (%s)", err)
// 			}
// 		} else {
// 			if err == nil {
// 				t.Fatal("expected error")
// 			} else if err.Error() != test.errMsg {
// 				t.Fatalf("unexpected error (%s)", err)
// 			}
// 		}
// 	}
// }

func TestFetch(t *testing.T) {
	t.Log("Testing Fetch")

	viper.Reset()
	client := genMockClient()

	tests := []struct {
		name        string
		cli         API
		id          string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid state (nil)", nil, "", true, "invalid state, nil client"},
		{"invalid id (empty)", client, "", true, "invalid id (empty)"},
		{"invalid id", client, "foo", true, "invalid ruleset id (foo)"},
		{"valid (apierr)", client, "000_error", true, "fetch api: forced mock api call error"},
		{"valid id", client, "123_test_metric", false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			_, err := Fetch(tst.cli, tst.id)
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
