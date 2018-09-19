// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"io/ioutil"
	"testing"
)

func TestDisplayList(t *testing.T) {
	t.Log("Testing DisplayList")

	t.Log("\tinvalid state (nil client)")
	if err := DisplayList(nil, nil); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid state, nil client" {
		t.Fatalf("unexpected error (%v)", err)
	}

	client := genMockClient()

	t.Log("\tinvalid state (nil dest)")
	if err := DisplayList(client, nil); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid destination (nil)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid")
	if err := DisplayList(client, ioutil.Discard); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}

func TestList(t *testing.T) {
	t.Log("Testing List")

	t.Log("\tinvalid state (nil client)")
	if _, err := List(nil); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid state, nil client" {
		t.Fatalf("unexpected error (%v)", err)
	}

	client := genMockClient()

	t.Log("\tvalid")
	if _, err := List(client); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}
