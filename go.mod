module github.com/konnek/konnek-aws

go 1.13

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.1 // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.1.1 // indirect
	github.com/aws/aws-lambda-go v1.16.0
	github.com/cloudevents/sdk-go v1.1.2
	github.com/cloudevents/sdk-go/v2 v2.0.0-preview8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/knative/eventing-sources v0.13.3 // indirect
	knative.dev/eventing v0.13.5 // indirect
	knative.dev/pkg v0.0.0-20200402224918-0cf29f826c40 // indirect
	sigs.k8s.io/controller-runtime v0.5.2
)
