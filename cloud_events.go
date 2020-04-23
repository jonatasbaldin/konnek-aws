package main

import (
	cloudevents "github.com/cloudevents/sdk-go"
)

func newCloudEventsClient(cloudEventConsumer string) (cloudevents.Client, error) {
	transport, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(cloudEventConsumer),
	)
	if err != nil {
		return nil, err
	}

	client, err := cloudevents.NewClient(
		transport,
		cloudevents.WithDataContentType(cloudevents.ApplicationJSON),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
