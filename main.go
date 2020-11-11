package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func getEnvOrQuit(name string) string {
	res := os.Getenv(name)
	if res == "" {
		log.Fatal("Missing environmental variable: ", name)
	}
	return res
}

var (
	topic    *pubsub.Topic
	consumed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "consumed_messages",
		Help: "The total number of messages consumed from subscription, acked or not",
	})
	acked = promauto.NewCounter(prometheus.CounterOpts{
		Name: "acked_messages",
		Help: "The total number of messages acked",
	})
)

func handleMessage(ctx context.Context, m *pubsub.Message) {
	consumed.Inc()
	log.Debug("publishing...")
	res := topic.Publish(ctx, m)
	_, err := res.Get(ctx)
	if err != nil {
		log.Error(err)
		m.Nack()
	}
	log.Debug("published")
	m.Ack()
	acked.Inc()
	log.Debug("acked")
}

func main() {
	metricsPortStr, isSet := os.LookupEnv("METRICS_PORT")
	if !(isSet && metricsPortStr == "") {
		if metricsPortStr == "" {
			metricsPortStr = "2121"
		}
		if _, err := strconv.Atoi(metricsPortStr); err != nil {
			log.Fatal("bad port given in METRICS_PORT: ", metricsPortStr)
		}
		http.Handle("/metrics", promhttp.Handler())
		log.Info("serving metrics on port ", metricsPortStr, ", url /metrics")
		go http.ListenAndServe(":"+metricsPortStr, nil)
	}
	ctx := context.Background()
	sourceClient, err := pubsub.NewClient(ctx, getEnvOrQuit("SOURCE_PROJECT_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	sinkClient, sinkerr := pubsub.NewClient(ctx, getEnvOrQuit("SINK_PROJECT_NAME"))
	if err != nil {
		log.Fatal(sinkerr)
	}
	topic = sinkClient.Topic(getEnvOrQuit("SINK_TOPIC_NAME"))
	sub := sourceClient.Subscription(getEnvOrQuit("SOURCE_SUBSCRIPTION_NAME"))
	err = sub.Receive(ctx, handleMessage)
	if err != nil {
		log.Fatal("Unable to receive messages")
	}
}
