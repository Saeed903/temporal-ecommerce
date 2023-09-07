package main

import (
	"log"
	"net"
	"os"

	app "github.com/saeed903/temporal-ecommerce"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var(
	stripKey = os.Getenv("STRIPE_PRIVATE_KEY")
	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunKey = os.Getenv("MAILGUN_KEY")
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: net.JoinHostPort("localhost", "7233"),
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client ", err)
	}
	defer c.Close()

	w := worker.New(c, "CART_TASK_QUEUE", worker.Options{})

	if stripKey == "" {
		log.Fatalln("Must set STRIPE_PRIVATE_KEY environment variable")
	}

	if mailgunDomain == "" {
		log.Fatalln("Must set MAILGUN_KEY environment variable")
	}

	if mailgunKey == "" {
		log.Fatalln("Must set MAILGUN_KEY environment variable")
	}

	a := &app.Activities{
		StripeKey: stripKey,
		MailgunDomain: mailgunDomain,
		MailgunKey: mailgunKey,
	}

	w.RegisterActivity(a.CreateStripeCharge)
	w.RegisterActivity(a.SendAbandonedCartEmail)

	w.RegisterWorkflow(app.CartWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	} 

}