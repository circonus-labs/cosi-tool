// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import "testing"

func TestDelete(t *testing.T) {
	t.Log("Test Delete")

	client := genMockClient()

	tests := []struct {
		name        string
		cid         string
		infile      string
		shouldFail  bool
		expectedErr string
	}{
		{"invalid (no args)", "", "", true, "missing required argument identifying ruleset"},
		{"invalid cid (foo)", "foo", "", true, "invalid ruleset id (foo)"},
		{"valid cid", "123_test_metric", "", false, ""},
		{"invalid in file (missing)", "", "testdata/missing.json", true, "reading configuration file: open testdata/missing.json: no such file or directory"},
		{"invalid in file (parsing)", "", "testdata/bad.json", true, "loading configuration: unexpected end of JSON input"},
		{"valid in file (apierr)", "", "testdata/api-error.json", true, "Circonus API error deleting ruleset: forced mock api call error"},
		{"valid in file", "", "testdata/registration-ruleset-test.json", false, ""},
	}

	t.Log("\tinvalid client")
	if err := Delete(nil, "", ""); err == nil {
		t.Fatal("expected error")
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			err := Delete(client, tst.cid, tst.infile)
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
