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
	settings, err := NewSettings()
	if err != nil {
		log.Fatalf("failed to load settings: %v", err)
	}

	url := settings.Host + "/v1/upgrade"
	ctx := cloudevents.ContextWithTarget(context.Background(), url)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	if settings.ClientCertPath != "" && settings.ClientKeyPath != "" {
		p.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{loadCertificate(settings.ClientCertPath, settings.ClientKeyPath)},
				InsecureSkipVerify: true,
			},
		}
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	data := &model.Upgrade{
		IstioVersion: settings.IstioVersion,
		ClusterName:  settings.ClusterName,
	}

	e := cloudevents.NewEvent()
	e.SetType("upgrade-event")
	e.SetSource(settings.EventSource)
	e.SetSubject(settings.EventSubject)

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

func loadCertificate(certPath, keyPath string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load certificate: %v", err)
	}
	return cert
}
