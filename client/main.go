package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	pbcloudevents "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	model "github.com/gopaytech/istio-upgrade-proto/upgrade"
)

func main() {
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/v1/upgrade")

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	// Load the certificate
	p.Client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{loadCertificate()},
		},
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	data := &model.Upgrade{
		IstioVersion: "1.22.4",
		ClusterName:  "s-go-sy-primary-gke-01",
	}
	e := cloudevents.NewEvent()
	e.SetType("upgrade-event")
	e.SetSource("testing-client")
	_ = e.SetData(pbcloudevents.ContentTypeProtobuf, data)

	res := c.Send(ctx, e)
	if cloudevents.IsUndelivered(res) {
		log.Printf("Failed to send: %v", res)
	} else {
		var httpResult *cehttp.Result
		if cloudevents.ResultAs(res, &httpResult) {
			log.Printf("Sent with status code %d", httpResult.StatusCode)
		} else {
			log.Printf("Send did not return an HTTP response: %s", res)
		}
	}
}

func loadCertificate() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("path/to/certificate.crt", "path/to/private.key")
	if err != nil {
		log.Fatalf("failed to load certificate: %v", err)
	}
	return cert
}
