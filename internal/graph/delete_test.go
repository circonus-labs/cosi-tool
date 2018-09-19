// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

import "testing"

func TestDelete(t *testing.T) {
	t.Log("Test Delete")

	t.Log("\tinvalid client")
	if err := Delete(nil, "", ""); err == nil {
		t.Fatal("expected error")
	}

	client := genMockClient()

	t.Log("\tinvalid (no args)")
	if err := Delete(client, "", ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "missing required argument identifying graph" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tinvalid id (foo)")
	if err := Delete(client, "foo", ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid graph id (foo)" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tvalid id (123)")
	if err := Delete(client, "123", ""); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid input file (missing)")
	if err := Delete(client, "", "testdata/missing.json"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "reading configuration file: open testdata/missing.json: no such file or directory" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tinvalid input file (parsing)")
	if err := Delete(client, "", "testdata/bad.json"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "loading configuration: unexpected end of JSON input" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tvalid input file (api error)")
	if err := Delete(client, "", "testdata/api-error.json"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "Circonus API error deleting graph: forced mock api call error" {
		t.Fatalf("expected different error, got (%v)", err)
	}

	t.Log("\tvalid input file")
	if err := Delete(client, "", "testdata/registration-graph-test.json"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}
