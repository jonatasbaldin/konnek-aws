# Konnek AWS
Transform AWS events into CloudEvents – and send them somewhere.

> **This is a proof of concept, so everything – from code to instructions – are not production ready. If the idea is valid, I'll bring it to an alpha version soon.**

## Idea
Konnek is a PoC trying to encapsulate cloud provider events into the [CloudEvents](https://cloudevents.io/) specification and forward them through the CloudEvents HTTP protocol-binding. The idea is to receive the events inside the cloud provider FaaS platform, parse, format and them send off.

The original idea was to feed these events into the Knative Eventing system, see [here](https://github.com/jonatasbaldin/konnek-event-receiver) for more info.

This repository contains the AWS Lambda implementation.

## Installing

### Setting up a local receiver 
Before we deploy the Lambda function, let's setup a place for these events to arrive. You'll need [ngrok](https://ngrok.com/).
```bash
docker run --rm -p 8080:8080 jonatasbaldin/konnek-knative-consumer

# Open another terminal and:
ngrok http 8080
```

Let ngrok run in this terminal. Note down your Ngrok address `https://xxxxxxx.ngrok.io`.

### Deploying the Lambda function
In _yet_ another terminal:

```bash
# Build Go binary
GOOS=linux go build -o main

# Zip it up!
zip function.zip main

# Create a basic AWS Lambda Role
# Note down the Role.Arn output, we will use in a bit
aws iam create-role --role-name konnek-lambda-role --assume-role-policy-document file://aws-basic-role.json

# Give it the AWSLambdaBasicExecutionRole policy
aws iam attach-role-policy --role-name konnek-lambda-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

# Deploy the function!
aws lambda create-function --function-name konnek --runtime go1.x --zip-file fileb://function.zip --environment "Variables={KONNEK_CE_CONSUMER=<your-ngrok-address>}" --handler main --role arn:aws:iam::<id>:role/konnek-lambda-role
```

Once deployed, test it:
```bash
aws lambda invoke --function-name konnek --payload fileb://testdata/sqs.json out.txt
```

Cool! Just see the logs in the terminal where you started the `docker` container!
```bash
2020/04/03 17:58:29 Validation: valid
Context Attributes,
  specversion: 1.0
  type: com.aws.sqs
  source: arn:aws:lambda:eu-central-1:580317195889:function:konnek
  id: 7cb16f6f-c1df-4632-8df6-2f583310e454
  time: 2020-04-03T17:58:29.07743932Z
  datacontenttype: application/json
Extensions,
  traceparent: 00-8280f5336cdec53b8fff591cd950c18b-9500a90e3d556f65-00
Data,
  {
    "Records": [
      {
        "attributes": {
          "ApproximateFirstReceiveTimestamp": "1523232000001",
          "ApproximateReceiveCount": "1",
          "SenderId": "123456789012",
          "SentTimestamp": "1523232000000"
        },
        "awsRegion": "eu-central-1",
        "body": "Hello from SQS!",
        "eventSource": "aws:sqs",
        "eventSourceARN": "arn:aws:sqs:eu-central-1:123456789012:MyQueue",
        "md5OfBody": "7b270e59b47ff90a553787216d55d91d",
        "messageAttributes": {},
        "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
        "receiptHandle": "MessageReceiptHandle"
      }
    ]
  }
```

Testing with a real SQS queue:
```bash
# Attach the SQS policy to the Lambda function
aws iam attach-role-policy --role-name konnek-lambda-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole

# Create a new SQS queue
# Note down the queue URL!
aws sqs create-queue --queue-name konnek

# Subscribe the Lambda function to the SQS queue
# The only way I found to get the SQS ARN is through the AWS dashboard :(
aws lambda create-event-source-mapping --function-name konnek --event-source-arn <sqs-queue-arn> --batch-size 1

# Send a message to SQS queue
aws sqs send-message --queue-url <queue-url> --message-body "sup!"
```

## After PoC
List of things to have for a more production ready piece of software:
- Implement more events, **only SQS and S3 are implemented**
- Implement in other cloud providers, like Google Cloud Functions and Azure Functions
- Deal with failure when delivering events (maybe use the cloud provider's queue system)
- Implement a proper controller for Knative
- Some form of authentication between konnek and receiver

So, what do you think?