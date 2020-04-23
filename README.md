# Konnek AWS
Transform AWS events into CloudEvents – and send them somewhere.

# Getting Started

## Create a Local Receiver
Let's create a setup to receive events directly in your machine.

Open one terminal and run the following Docker container and keep it running, as the event will be shown in the logs:
```bash
docker run --rm -p 8080:8080 jonatasbaldin/konnek-consumer
```

Open another terminal and use [https://ngrok.com/] to expose the Docker container to the Internet, so Konnek can reach you:
```bash
ngrok http 8080
```

Take note on your ngrok address (https://xxxxxxxx.ngrok.io), we will use it in a bit.

## Installing with Serverless Framework
Assuming you already have the Serverless Framework [installed](https://serverless.com/framework/docs/getting-started/) and the AWS credentials [configured](https://serverless.com/framework/docs/providers/aws/cli-reference/config-credentials/), get the latest Konnek version:
```bash
wget https://github.com/jonatasbaldin/konnek-aws/releases/download/v0.0.2/konnek-aws-0.0.2.zip -O konnek.zip
```

Get the official Konnek `serverless.yml` file: 
```bash
wget https://raw.githubusercontent.com/jonatasbaldin/konnek-aws/master/config/serverless-framework/serverless.yml
```

Set the `KONNEK_CONSUMER` environment variable to the Ngrok address generated before and deploy it:
```bash
KONNEK_CONSUMER="https://xxxxxxxx.ngrok.io" serverless deploy
```

Get a SQS mock data file:
```bash
wget https://raw.githubusercontent.com/jonatasbaldin/konnek-aws/master/testdata/sqs.json
```

Finally, test it:
```bash
serverless invoke -f konnek -p sqs.json
```

Look in your Docker terminal, you should see something like:
```bash
Context Attributes,
  specversion: 1.0
  type: com.amazon.sqs
  source: arn:aws:sqs:eu-central-1:123456789012:MyQueue
  id: cfbb5e8d-f025-4a9c-9b7a-55d10a4b42e2
  time: 2020-04-23T16:56:44.066357394Z
  datacontenttype: application/json
Extensions,
  traceparent: 00-a0b1fe0032d5ad22af09319d51271ded-916075d55bd96bfe-00
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

That's your event!


## Installing Manually
Use the following steps to install Konnek manually (using AWS CLI):
```bash
# Get Konnek latest version
wget https://github.com/jonatasbaldin/konnek-aws/releases/download/v0.0.2/konnek-aws-0.0.2.zip -O konnek.zip

# Get a basic AWS Role template for the Lambda
wget https://raw.githubusercontent.com/jonatasbaldin/konnek-aws/master/config/aws-basic-role.json

# Create a basic AWS Lambda Role
# Note down the Role.Arn output, we will use in a bit
aws iam create-role --role-name konnek-lambda-role --assume-role-policy-document file://aws-basic-role.json

# Give it the AWSLambdaBasicExecutionRole policy
aws iam attach-role-policy --role-name konnek-lambda-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

# Deploy the function!
# Add the Role.Arn in the <id> field on the --role option
aws lambda create-function --function-name konnek --runtime go1.x --zip-file fileb://konnek.zip --environment "Variables={KONNEK_CONSUMER=<your-ngrok-address>}" --handler main --role arn:aws:iam::<id>:role/konnek-lambda-role
```

Once deployed, test it:
```bash
# Get a SQS mock data file:
wget https://raw.githubusercontent.com/jonatasbaldin/konnek-aws/master/testdata/sqs.json

aws lambda invoke --function-name konnek --payload fileb://sqs.json out.txt
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