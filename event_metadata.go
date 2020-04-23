package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type EventMetadata struct {
	Type   string
	Source string
	Id     string
}

func getEventMetadata(lambdaCxt *lambdacontext.LambdaContext, event map[string]interface{}) (eventMetadata *EventMetadata, err error) {
	eventMetadata = &EventMetadata{}

	eventMetadata.Id = lambdaCxt.AwsRequestID

	eventBody, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("could not marshal event into bytes: %v", err)
	}

	var sqs events.SQSEvent
	err = json.Unmarshal(eventBody, &sqs)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal event: %v", err)
	}

	if len(sqs.Records) > 0 && strings.Contains(sqs.Records[0].EventSource, "aws:sqs") {
		eventMetadata.Type = SQSType
		eventMetadata.Source = sqs.Records[0].EventSourceARN

		return
	}

	var s3 events.S3Event
	err = json.Unmarshal(eventBody, &s3)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal event: %v", err)
	}

	if len(s3.Records) > 0 && strings.Contains(s3.Records[0].EventSource, "aws:s3") {
		eventMetadata.Source = s3.Records[0].S3.Bucket.Arn

		if s3.Records[0].EventName == "ObjectCreated:Put" {
			eventMetadata.Type = S3PutType
		}
		if s3.Records[0].EventName == "ObjectRemoved:Delete" {
			eventMetadata.Type = S3DeleteType
		}

		return
	}

	var sns events.SNSEvent
	err = json.Unmarshal(eventBody, &sns)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal event: %v", err)
	}

	if len(sns.Records) > 0 && strings.Contains(sns.Records[0].EventSource, "aws:sns") {
		eventMetadata.Source = sns.Records[0].EventSubscriptionArn
		eventMetadata.Type = SNSType

		return
	}

	return
}
