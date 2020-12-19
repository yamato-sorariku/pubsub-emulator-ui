package controller

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Env struct {
	ProjectId      string
	TopicId        string
	SubscriptionId string
}

func SetUpPubSub() {
	err := createTopic()
	if err != nil {
		println(err.Error())
		return
	}
	err = createPushSubscription()
	if err != nil {
		println(err.Error())
		return
	}
}

type JsonRequest struct {
	Message string `json:"message"`
}

func Publish(c *gin.Context) {
	env := getPubSubEnv()

	var json JsonRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.ProjectId)
	if err != nil {
		fmt.Printf("pubsub.NewClient: %v\n", err)
		return
	}

	t := client.Topic(env.TopicId)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(json.Message),
	})
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Get: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complete message publish."})
}

func PullPubSubMessage() error {
	env := getPubSubEnv()

	PingToClients()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.ProjectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(env.SubscriptionId)
	sub.ReceiveSettings.Synchronous = false
	sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var counter int32

	var mu sync.Mutex
	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		atomic.AddInt32(&counter, 1)
		BroadcastMessagesToClients(msg)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("pubsub: Receive returned error: %v", err)
	}
	fmt.Sprintf("Received %d messages\n", counter)

	return nil
}

func createTopic() error {
	env := getPubSubEnv()

	// Topicの作成
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.ProjectId)
	if err != nil {
		return fmt.Errorf("Unable to create client to project %q: %s", env.ProjectId, err)
	}
	defer client.Close()

	topic := client.Topic(env.TopicId)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("Unable to topic %q for project %q: %s", env.TopicId, env.ProjectId, err)
	}
	if !ok {
		_, err = client.CreateTopic(ctx, env.TopicId)
		if err != nil {
			return fmt.Errorf("Unable to create topic %q for project %q: %s", env.TopicId, env.ProjectId, err)
		}
	}
	return nil
}
func createPushSubscription() error {
	env := getPubSubEnv()

	// Topicの作成
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, env.ProjectId)
	if err != nil {
		return fmt.Errorf("Unable to create client to project %q: %s", env.ProjectId, err)
	}
	defer client.Close()

	topic := client.Topic(env.TopicId)
	_, _ = client.CreateSubscription(
		ctx,
		env.SubscriptionId,
		pubsub.SubscriptionConfig{
			Topic: topic,
		},
	)
	return nil
}

func getPubSubEnv() Env {
	return Env{
		ProjectId:      os.Getenv("PUBSUB_PROJECT_ID"),
		TopicId:        os.Getenv("PUBSUB_TOPIC_ID"),
		SubscriptionId: os.Getenv("PUBSUB_SUBSCRIPTION_ID"),
	}
}
