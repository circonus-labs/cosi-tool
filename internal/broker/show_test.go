// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"io/ioutil"
	"testing"
)

func TestShow(t *testing.T) {
	t.Log("Testing Show")

	t.Log("\tinvalid state (nil client)")
	if err := Show(nil, nil, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid state, nil client" {
		t.Fatalf("unexpected error (%v)", err)
	}

	client := genMockClient()

	t.Log("\tinvalid state (nil dest)")
	if err := Show(client, nil, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid destination (nil)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid cid (empty)")
	if err := Show(client, ioutil.Discard, ""); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid broker id (empty)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid cid (foo)")
	if err := Show(client, ioutil.Discard, "foo"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid broker id (foo)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid (api error)")
	if err := Show(client, ioutil.Discard, "000"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "fetch api: forced mock api call error" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid (inactive)")
	if err := Show(client, ioutil.Discard, "/broker/456"); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "Broker /broker/456 is not active" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid (active)")
	if err := Show(client, ioutil.Discard, "123"); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}
