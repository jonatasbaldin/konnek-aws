package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventsclient "github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

const (
	ceSqsType      = "com.aws.sqs"
	ceS3PutType    = "com.aws.s3.put"
	ceS3DeleteType = "com.aws.s3.delete"
)

type EnvConfig struct {
	CloudEventsConsumer string `envconfig:"KONNEK_CE_CONSUMER" required:"true"`
}

func NewCloudEventsClient(cloudEventConsumer string) (cloudeventsclient.Client, error) {
	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(cloudEventConsumer),
		cloudeventshttp.WithEncoding(cloudeventshttp.Default),
	)
	if err != nil {
		return nil, err
	}

	client, err := cloudeventsclient.New(
		transport,
		cloudevents.WithDataContentType(cloudevents.ApplicationJSON),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getEventType(event interface{}) (eventType string, err error) {
	var sqs events.SQSEvent
	sqsByte, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	json.Unmarshal(sqsByte, &sqs)

	if len(sqs.Records) > 0 && strings.Contains(sqs.Records[0].EventSource, "aws:sqs") {
		return ceSqsType, nil
	}

	var s3 events.S3Event
	s3Byte, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	json.Unmarshal(s3Byte, &s3)

	if len(s3.Records) > 0 && strings.Contains(s3.Records[0].EventSource, "aws:s3") {
		if s3.Records[0].EventName == "ObjectCreated:Put" {
			return ceS3PutType, nil
		}
		if s3.Records[0].EventName == "ObjectRemoved:Delete" {
			return ceS3DeleteType, nil
		}
	}

	return "", fmt.Errorf("unable to get event type")
}

func _main(ctx context.Context, event interface{}) {
	log.Printf("event is: %+v\n", event)

	lambdaContext, _ := lambdacontext.FromContext(ctx)
	log.Printf("context is: %+v\n", lambdaContext)

	var envConfig EnvConfig
	err := envconfig.Process("", &envConfig)
	if err != nil {
		log.Fatalf("could not load environment variables: %v\n", err)
	}

	// Parser
	eventType, err := getEventType(event)
	if err != nil {
		log.Fatalf("could not get eventType: %v\n", err)
	}
	log.Printf("eventType is: %+v\n", eventType)

	// CE
	cloudEventsClient, err := NewCloudEventsClient(envConfig.CloudEventsConsumer)
	if err != nil {
		log.Fatalf("could not create client: %v\n", err)
	}

	parsedEventSource, err := url.Parse(lambdaContext.InvokedFunctionArn)
	if err != nil {
		log.Fatalf("could not parse eventSource: %v\n", err)
	}

	cloudEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			ID:     lambdaContext.AwsRequestID,
			Source: types.URIRef{URL: *parsedEventSource},
			Type:   eventType,
		},
		Data: event,
	}

	_, _, err = cloudEventsClient.Send(context.Background(), cloudEvent)
	if err != nil {
		log.Fatalf("could not send event: %v\n", err)
	}
}

func main() {
	lambda.Start(_main)
}
