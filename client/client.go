package client

import (
	"github.com/go-resty/resty/v2"
	"gitlab.com/apus-backend/base-service/config"
	"time"
)

type Client struct {
	MembershipClient *resty.Client
}

func NewClient(c config.Config) (*Client, error) {
	membershipClient, err := initRestyClient(c, c.HostUrl.Membership)
	if err != nil {
		return nil, err
	}

	return &Client{
		MembershipClient: membershipClient,
	}, nil
}

func initRestyClient(c config.Config, hostUrl string) (*resty.Client, error) {
	// Create a Resty Client
	client := resty.New()

	// Unique settings at Client level
	//--------------------------------
	// Enable debug mode
	if c.Resty.Debug == nil || !*c.Resty.Debug {
		client.SetDebug(false)
	} else {
		client.SetDebug(true)
	}

	// Set client timeout as per your need
	duration, err := time.ParseDuration(c.Resty.Timeout)
	if err != nil {
		return nil, err
	}
	client.SetTimeout(duration)

	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	client.SetHostURL(hostUrl)

	// Headers for all request
	client.SetHeader("Accept", "application/json")
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	return client, nil
}
