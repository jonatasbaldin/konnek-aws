package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func readJSONFile(filePath string) (content map[string]interface{}) {
	fileContent, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(fileContent, &content)

	return content
}

func Test_getEventMetadata(t *testing.T) {
	lambdaCtx := lambdacontext.LambdaContext{
		AwsRequestID: "69fc7648-d849-53f9-a871-215de0e1ec0e",
	}

	testCases := []struct {
		name          string
		event         map[string]interface{}
		eventMetadata *EventMetadata
	}{
		{
			name:  SQSType,
			event: readJSONFile("testdata/sqs.json"),
			eventMetadata: &EventMetadata{
				Type:   SQSType,
				Source: "arn:aws:sqs:eu-central-1:123456789012:MyQueue",
				Id:     "69fc7648-d849-53f9-a871-215de0e1ec0e",
			},
		},
		{
			name:  S3PutType,
			event: readJSONFile("testdata/s3-put.json"),
			eventMetadata: &EventMetadata{
				Type:   S3PutType,
				Source: "arn:aws:s3:::example-bucket",
				Id:     "69fc7648-d849-53f9-a871-215de0e1ec0e",
			},
		},
		{
			name:  S3DeleteType,
			event: readJSONFile("testdata/s3-delete.json"),
			eventMetadata: &EventMetadata{
				Type:   S3DeleteType,
				Source: "arn:aws:s3:::example-bucket",
				Id:     "69fc7648-d849-53f9-a871-215de0e1ec0e",
			},
		},
		{
			name:  SNSType,
			event: readJSONFile("testdata/sns.json"),
			eventMetadata: &EventMetadata{
				Type:   SNSType,
				Source: "arn:aws:sns:EXAMPLE",
				Id:     "69fc7648-d849-53f9-a871-215de0e1ec0e",
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			eventMetadata, err := getEventMetadata(&lambdaCtx, tC.event)

			if err != nil {
				t.Errorf("got err %v", err)
			}

			if !reflect.DeepEqual(eventMetadata, tC.eventMetadata) {
				t.Errorf("expected %+v, got %+v", tC.eventMetadata, eventMetadata)
			}
		})
	}
}
