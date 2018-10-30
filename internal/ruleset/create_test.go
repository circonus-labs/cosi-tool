// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import "testing"

func TestCreateFromFile(t *testing.T) {
	t.Log("Test CreateFromFile")

	client := genMockClient()

	tests := []struct {
		name        string
		client      *APIMock
		inFile      string
		outFile     string
		force       bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid input file (empty)", client, "", "", false, true, "invalid input file (empty)"},
		{"invalid input file (missing)", client, "testdata/missing.json", "", false, true, "reading configuration file: open testdata/missing.json: no such file or directory"},
		{"invalid input file (parsing)", client, "testdata/bad.json", "", false, true, "loading configuration: unexpected end of JSON input"},
		{"valid input file (apierr)", client, "testdata/api-error.json", "", false, true, "Circonus API error creating ruleset: forced mock api call error"},
		{"valid input file", client, "testdata/valid-ignore.json", "testdata/will_be_overwritten.json", true, false, ""},
		{"valid input file (no force)", client, "testdata/valid-ignore.json", "testdata/will_be_overwritten.json", false, true, "saving created ruleset: testdata/will_be_overwritten.json already exists, see --force"},
	}

	t.Log("\tinvalid client")
	if err := CreateFromFile(nil, "", "", false); err == nil {
		t.Fatal("expected error")
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			err := CreateFromFile(tst.client, tst.inFile, tst.outFile, tst.force)
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
