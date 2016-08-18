package irkit

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInternetClient_GetKeys(t *testing.T) {
	mux := http.NewServeMux()
	tc := httptest.NewServer(mux)

	token := "123456789"
	mux.HandleFunc(pathKeys, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		buf, _ := ioutil.ReadAll(r.Body)
		if !strings.Contains(string(buf), fmt.Sprintf("clienttoken=%s", token)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, "./testdata/getkeys.json")
	})
	defer tc.Close()

	c, err := newInternetClient(tc.URL)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	deviceid, clientkey, err := c.GetKeys(context.Background(), token)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if deviceid != "IIII" {
		t.Fatalf("expect %q to be eq %q", deviceid, "IIII")
	}

	if clientkey != "KKKK" {
		t.Fatalf("expect %q to be eq %q", deviceid, "KKKK")
	}

}

func TestInternetClient_SendMessages(t *testing.T) {
	mux := http.NewServeMux()
	tc := httptest.NewServer(mux)

	clientkey := "keykey"
	deviceid := "idid"
	mux.HandleFunc(pathMessages, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		buf, _ := ioutil.ReadAll(r.Body)
		if !strings.Contains(string(buf), fmt.Sprintf("clientkey=%s", clientkey)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !strings.Contains(string(buf), fmt.Sprintf("deviceid=%s", deviceid)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	defer tc.Close()

	c, err := newInternetClient(tc.URL)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	msg := &Message{
		Format: "raw",
		Freq:   38,
		Data:   []int{1, 2, 3, 45},
	}

	if err := c.SendMessages(context.Background(), clientkey, deviceid, msg); err != nil {
		t.Fatalf("err: %s", err)
	}
}
