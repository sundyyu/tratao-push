package ios

import (
	"encoding/json"
	"github.com/RobotsAndPencils/buford/certificate"
	"github.com/RobotsAndPencils/buford/payload"
	"github.com/RobotsAndPencils/buford/payload/badge"
	"github.com/RobotsAndPencils/buford/push"
)

const (
	host = push.Development
)

type Client struct {
	CerPath string
	CerPass string
}

func NewClient(cerPath string, cerPass string) *Client {
	return &Client{cerPath, cerPass}
}

func (c *Client) Push(n *Notification, TPR string) (*Response, error) {
	resp := &Response{}

	// load a certificate and use it to connect to the APN service:
	cert, err := certificate.Load(c.CerPath, c.CerPass)
	if err != nil {
		return nil, err
	}

	client, err := push.NewClient(cert)
	if err != nil {
		return nil, err
	}

	service := push.NewService(client, host)

	// construct a payload to send to the device:
	p := payload.APS{
		Alert: payload.Alert{Title: n.Title, Body: n.Body},
		Badge: badge.New(0),
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	topic := certificate.TopicFromCert(cert)
	headers := push.Headers{
		Topic: topic,
	}

	// push the notification:
	id, err := service.Push(TPR, &headers, b)
	if err != nil {
		return nil, err
	}

	resp.Message = id
	return resp, nil
}
