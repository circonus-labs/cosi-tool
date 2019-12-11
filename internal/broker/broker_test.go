// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var testCosiBrokers = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "error") {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if strings.Contains(r.URL.Path, "badjson") {
		fmt.Fprintln(w, `{bad:"json"}`)
		return
	} else if strings.Contains(r.URL.Path, "nojson") {
		x := defaultBroker{Trap: 456}
		if !strings.Contains(r.URL.Path, "nofallback") {
			x.Fallback = 123
		}
		_ = json.NewEncoder(w).Encode(x)
		return
	} else if strings.Contains(r.URL.Path, "notrap") {
		x := defaultBroker{JSON: 789}
		if !strings.Contains(r.URL.Path, "nofallback") {
			x.Fallback = 123
		}
		_ = json.NewEncoder(w).Encode(x)
		return
	}
	fmt.Fprintln(w, `{"fallback":123,"httptrap":456,"json":789}`)
}))

func genMockClient() *APIMock {
	return &APIMock{
		FetchBrokerFunc: func(cid circapi.CIDType) (*circapi.Broker, error) {
			switch *cid {
			case "/broker/000":
				return nil, errors.New("forced mock api call error")
			case "/broker/123":
				return &circapi.Broker{
					CID:  "/broker/123",
					Name: "foo",
					Type: "xxx",
					Details: []circapi.BrokerDetail{
						{
							Status:  "active",
							Modules: []string{"abc", "selfcheck", "hidden:abc123", "abcdef", "abcdefghi", "abcdefghijkl", "abcdefghijklmnopqrstu"},
						},
						{
							Status: "foobar",
						},
					},
				}, nil
			case "/broker/456":
				return &circapi.Broker{
					CID:  "/broker/456",
					Name: "bar",
					Type: "yyy",
					Details: []circapi.BrokerDetail{
						{
							Status: "foobar",
						},
					},
				}, nil
			default:
				return nil, errors.Errorf("bad broker request cid (%s)", *cid)
			}
		},
		FetchBrokersFunc: func() (*[]circapi.Broker, error) {
			return &[]circapi.Broker{
				{CID: "/broker/123", Name: "foo", Type: "circonus"},
				{CID: "/broker/456", Name: "bar", Type: "enterprise"},
				{CID: "/broker/789", Name: "baz", Type: "circonus"},
			}, nil
		},
	}
}

func TestLoadCustomBroker(t *testing.T) {
	t.Log("Testing loadCustomBroker")

	zerolog.SetGlobalLevel(zerolog.Disabled)
	logger := zerolog.New(os.Stdout)

	t.Log("\tno options file")
	if bid, db, err := loadCustomBroker("", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if bid != 0 && db.Fallback != 0 && db.JSON != 0 && db.Trap != 0 {
		t.Fatalf("unexpected result (bid:%d db:%#v)", bid, db)
	}

	t.Log("\tvalid option file (empty)")
	if bid, db, err := loadCustomBroker("testdata/blank.json", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if bid != 0 && db.Fallback != 0 && db.JSON != 0 && db.Trap != 0 {
		t.Fatalf("unexpected result (bid:%d db:%#v)", bid, db)
	}

	t.Log("\tinvalid (bad json)")
	if _, _, err := loadCustomBroker("testdata/bad.json", logger); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "parsing custom options file json: invalid character 'b' looking for beginning of object key string" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid option file (explicit broker)")
	if bid, _, err := loadCustomBroker("testdata/explicit.json", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if bid != 1 {
		t.Fatalf("unexpected result (bid:%d)", bid)
	}

	t.Log("\tvalid option file (no fallback)")
	if bid, db, err := loadCustomBroker("testdata/nofallback.json", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if !(db.Fallback == 0 && (db.Trap > 0 && db.JSON > 0)) {
		t.Fatalf("unexpected result (bid:%d db:%#v)", bid, db)
	}

	t.Log("\tvalid option file (no json)")
	if bid, db, err := loadCustomBroker("testdata/nojson.json", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if !(db.JSON == 0 && (db.Trap > 0 && db.Fallback > 0)) {
		t.Fatalf("unexpected result (bid:%d db:%#v)", bid, db)
	}

	t.Log("\tvalid option file (no trap)")
	if bid, db, err := loadCustomBroker("testdata/notrap.json", logger); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if !(db.Trap == 0 && (db.Fallback > 0 && db.JSON > 0)) {
		t.Fatalf("unexpected result (bid:%d db:%#v)", bid, db)
	}

}

func TestFetchCosiDefaults(t *testing.T) {
	t.Log("Testing fetchCosiDefaults")

	zerolog.SetGlobalLevel(zerolog.Disabled)
	logger := zerolog.New(os.Stdout)

	t.Log("\tinvalid cosi url (empty)")
	if _, err := fetchCosiDefaults("", logger); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "invalid cosi url (empty)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (api error)")
	if _, err := fetchCosiDefaults(testCosiBrokers.URL+"/error/", logger); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "unexpected return code from cosi api (500 - 500 Internal Server Error)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid state (bad json)")
	if _, err := fetchCosiDefaults(testCosiBrokers.URL+"/badjson/", logger); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "parsing cosi api json: invalid character 'b' looking for beginning of object key string" {
		t.Fatalf("unexpected error (%v)", err)
	}
}

func TestGetBrokerFromList(t *testing.T) {
	t.Log("Testing getBrokerFromList")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("\tinvalid (empty list)")
	if _, err := getBrokerFromList([]uint{}, 0); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "empty list" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid (invalid index -2)")
	if _, err := getBrokerFromList([]uint{1}, -2); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "index (-2) out of range for list size (1)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tinvalid (invalid index 5)")
	if _, err := getBrokerFromList([]uint{1}, 5); err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "index (5) out of range for list size (1)" {
		t.Fatalf("unexpected error (%v)", err)
	}

	t.Log("\tvalid (indexed)")
	if id, err := getBrokerFromList([]uint{1}, 0); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if id != 1 {
		t.Fatalf("unexpected result %d", id)
	}

	t.Log("\tvalid (random, list len 1)")
	if id, err := getBrokerFromList([]uint{1}, -1); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if id != 1 {
		t.Fatalf("unexpected result %d", id)
	}

	t.Log("\tvalid (random, list len > 1)")
	if id, err := getBrokerFromList([]uint{1, 2, 3}, -1); err != nil {
		t.Fatalf("unexpected error (%v)", err)
	} else if id < 1 || id > 3 {
		t.Fatalf("unexpected result %d", id)
	}
}
