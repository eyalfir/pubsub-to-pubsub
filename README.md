# pubsub-to-pubsub

## Overview

This docker image is pulling messages from a pubsub subscription and pubslished to a pubsub topic

## Configuration

| environment variable | optional | default | description | example
|-|-|-|-|-|
| SOURCE_PROJECT_NAME | no | | GCP project name ||
| SOURCE_SUBSCRIPTION_NAME | no || name of subscription to subscribe to | my-subsbcription |
| SINK_PROJECT_NAME | no | | GCP project name ||
| SINK_TOPIC_NAME | no || name of topic to publish to | my-topic |
| LOG_LEVEL | yes | info | log level | debug |
| METRICS_PORT | yes | 2121 | port to expose prometheus metrics. set to empty string to skip metrics endpoint | 3434 |
