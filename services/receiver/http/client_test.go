package http

import (
	"bytes"
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	pbcloudevents "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/gopaytech/istio-upgrade-consumer/settings"
	"github.com/gopaytech/istio-upgrade-consumer/types"
	"github.com/gopaytech/istio-upgrade-consumer/usecases/upgrade"
	model "github.com/gopaytech/istio-upgrade-proto/upgrade"
)

type fakeStorage struct {
	created []types.Upgrade
}

func (f *fakeStorage) Create(_ context.Context, u types.Upgrade) error {
	f.created = append(f.created, u)
	return nil
}

func (f *fakeStorage) Get(context.Context) (*types.Upgrade, error) {
	return nil, errors.New("not found")
}

func (f *fakeStorage) Update(context.Context, types.Upgrade) error {
	return errors.New("not found")
}

func newTestServer(t *testing.T) (*httptest.Server, *fakeStorage, *nethttp.Server) {
	t.Helper()

	s := settings.Settings{
		ClusterName:        "my-cluster",
		ClusterEnvironment: "production",
		ReceiverHTTPPort:   8080,
		TimeLocation:       "UTC",
		TimeFormat:         "2006-01-02",
	}
	storage := &fakeStorage{}
	server, err := NewUpgradeHTTP(s, storage, upgrade.NewUpgradeImplementation(s, storage)).Server()
	if err != nil {
		t.Fatalf("failed to build server: %v", err)
	}

	ts := httptest.NewServer(server.Handler)
	t.Cleanup(ts.Close)
	return ts, storage, server
}

func sendUpgradeEvent(t *testing.T, target string) int {
	t.Helper()

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	e := cloudevents.NewEvent()
	e.SetID("1")
	e.SetType("upgrade-event")
	e.SetSource("test")
	if err := e.SetData(pbcloudevents.ContentTypeProtobuf, &model.Upgrade{IstioVersion: "1.27.1", ClusterName: "my-cluster"}); err != nil {
		t.Fatalf("failed to set data: %v", err)
	}

	res := c.Send(cloudevents.ContextWithTarget(context.Background(), target), e)
	var httpResult *cehttp.Result
	if !cloudevents.ResultAs(res, &httpResult) {
		t.Fatalf("send did not return an HTTP result: %v", res)
	}
	return httpResult.StatusCode
}

func TestServerStoresUpgradeEventOnUpgradePath(t *testing.T) {
	ts, storage, _ := newTestServer(t)

	if code := sendUpgradeEvent(t, ts.URL+upgradePath); code != nethttp.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}

	if len(storage.created) != 1 {
		t.Fatalf("expected 1 stored upgrade, got %d", len(storage.created))
	}
	if got := storage.created[0]; got.Version != "1.27.1" || got.ClusterName != "my-cluster" {
		t.Fatalf("unexpected stored upgrade: %+v", got)
	}
}

func TestServerRejectsOtherPaths(t *testing.T) {
	ts, storage, _ := newTestServer(t)

	for _, path := range []string{"/", "/v1/upgrade/", "/v1/upgrade/x", "/v1"} {
		if code := sendUpgradeEvent(t, ts.URL+path); code != nethttp.StatusNotFound {
			t.Errorf("%s: expected 404, got %d", path, code)
		}
	}

	if len(storage.created) != 0 {
		t.Fatalf("expected nothing stored, got %d", len(storage.created))
	}
}

func TestServerRejectsGet(t *testing.T) {
	ts, _, _ := newTestServer(t)

	res, err := nethttp.Get(ts.URL + upgradePath)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	if res.StatusCode != nethttp.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.StatusCode)
	}
}

func TestServerRejectsOversizedBody(t *testing.T) {
	ts, storage, _ := newTestServer(t)

	req, err := nethttp.NewRequest(nethttp.MethodPost, ts.URL+upgradePath, bytes.NewReader(make([]byte, maxBodyBytes+1)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("ce-specversion", "1.0")
	req.Header.Set("ce-id", "1")
	req.Header.Set("ce-type", "upgrade-event")
	req.Header.Set("ce-source", "test")
	req.Header.Set("content-type", pbcloudevents.ContentTypeProtobuf)

	res, err := nethttp.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	if res.StatusCode < 400 {
		t.Fatalf("expected an error status, got %d", res.StatusCode)
	}
	if len(storage.created) != 0 {
		t.Fatalf("expected nothing stored, got %d", len(storage.created))
	}
}

func TestServerTimeouts(t *testing.T) {
	_, _, server := newTestServer(t)

	if server.ReadHeaderTimeout != readHeaderTimeout || server.ReadTimeout != readTimeout ||
		server.WriteTimeout != writeTimeout || server.IdleTimeout != idleTimeout {
		t.Fatalf("unexpected timeouts: %+v", server)
	}
	if server.MaxHeaderBytes != maxHeaderBytes {
		t.Fatalf("expected MaxHeaderBytes %d, got %d", maxHeaderBytes, server.MaxHeaderBytes)
	}
}
