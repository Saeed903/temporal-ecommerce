package main

import (
	"os"

	"go.temporal.io/sdk/client"
)

type (
	ErrorResponse struct {
		Message string
	}

	UpdateEmailRequest struct {
		Email string
	}

	CheckoutRequest struct {
		Email string
	}
)

var (
	HttpPort = os.Getenv("PORT")
	temporal client.Client
)
