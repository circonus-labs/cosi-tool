// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

import (
	"testing"
)

func TestUpdateFromFile(t *testing.T) {
	t.Log("Test UpdateFromFile")

	t.Log("\tinvalid client")
	if err := UpdateFromFile(nil, "", "", false); err == nil {
		t.Fatal("expected error")
	}

	client := genMockClient()

	t.Log("\tinvalid input file (empty)")
	if err := UpdateFromFile(client, "", "", false); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid input file (empty)" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tinvalid input file (missing)")
	if err := UpdateFromFile(client, "testdata/missing.json", "", false); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "reading configuration file: open testdata/missing.json: no such file or directory" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tinvalid input file (parsing)")
	if err := UpdateFromFile(client, "testdata/bad.json", "", false); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "loading configuration: unexpected end of JSON input" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tvalid input file (api error)")
	if err := UpdateFromFile(client, "testdata/api-error.json", "", false); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "Circonus API error updating check bundle: forced mock api call error" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tvalid input file")
	if err := UpdateFromFile(client, "testdata/registration-check-system.json", "testdata/will_be_overwritten.json", true); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid input file (no force)")
	if err := UpdateFromFile(client, "testdata/registration-check-system.json", "testdata/will_be_overwritten.json", false); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "saving updated check bundle: testdata/will_be_overwritten.json already exists, see --force" {
		t.Fatalf("expected different error, got (%v)", err)
	}

}
