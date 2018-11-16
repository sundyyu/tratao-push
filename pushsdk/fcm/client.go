package fcm

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type Client struct {
	CfgPath string
}

func NewClient(cfgPath string) *Client {
	return &Client{cfgPath}
}

func (c *Client) Push(n *Notification, TPR string) (*Response, error) {
	opt := option.WithCredentialsFile(c.CfgPath) // "/Users/traveltao/xcpushdemo-firebase-adminsdk-iy30t-340376f811.json"
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := TPR

	// See documentation on defining a message payload.
	message := &messaging.Message{
		// Data: map[string]string{
		// 	"body":  n.Title,
		// 	"title": n.Body,
		// },
		Notification: &messaging.Notification{
			Title: n.Title,
			Body:  n.Body,
		},
		Token: registrationToken,
	}

	resp := &Response{}

	// Send a message to the device corresponding to the provided registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	resp.Message = response
	return resp, nil
}
