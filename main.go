package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

const (
	SQSType      = "com.amazon.sqs"
	S3PutType    = "com.amazon.s3.put"
	S3DeleteType = "com.amazon.s3.delete"
	SNSType      = "com.amazon.sns"
)

type EnvConfig struct {
	CloudEventsConsumer string `envconfig:"KONNEK_CE_CONSUMER" required:"true"`
}

func _main(ctx context.Context, event map[string]interface{}) {
	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Fatal("could not load lambda context")
	}

	var env EnvConfig
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("could not load environment variables: %v", err)
	}

	eventMetadata, err := getEventMetadata(lambdaCtx, event)
	if err != nil {
		log.Fatalf("could not get eventMetadata: %v", err)
	}
	log.Printf("eventMetadata is: %+v", eventMetadata)

	cloudEventsClient, err := newCloudEventsClient(env.CloudEventsConsumer)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	cloudEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			ID:     eventMetadata.Id,
			Source: *types.ParseURIRef(eventMetadata.Source),
			Type:   eventMetadata.Type,
		},
		Data: event,
	}

	_, _, err = cloudEventsClient.Send(context.Background(), cloudEvent)
	if err != nil {
		log.Fatalf("could not send event: %v", err)
	}

	log.Printf("event with id %s send to %s", eventMetadata.Id, env.CloudEventsConsumer)
}

func main() {
	lambda.Start(_main)
}
