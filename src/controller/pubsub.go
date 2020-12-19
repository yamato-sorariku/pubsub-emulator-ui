package controller

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const projectID = "my-project-id"
const topicID = "my-topic"
const subscriptionID = "my-topic-subscription"

func Pubsub(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello world!!"})
}

func SetPushEndpoint(c *gin.Context) {
	err := createTopic()
	if err != nil {
		fatalf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed set push endpoint."})
		return
	}
	err = createPushSubscription()
	if err != nil {
		fatalf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed set push endpoint."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complete set endpoint to PubSub."})
}

func Publish(c *gin.Context) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("pubsub.NewClient: %v\n", err)
		return
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("TEST NOW!!"),
	})
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Get: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complete message publish."})
}

func createTopic() error {
	// Topicの作成
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Unable to create client to project %q: %s", projectID, err)
	}
	defer client.Close()

	topic := client.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("Unable to topic %q for project %q: %s", topicID, projectID, err)
	}
	if !ok {
		_, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			return fmt.Errorf("Unable to create topic %q for project %q: %s", topicID, projectID, err)
		}
	}
	return nil
}
func createPushSubscription() error {
	// Topicの作成
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Unable to create client to project %q: %s", projectID, err)
	}
	defer client.Close()
	debugf("Client connected with project ID %q", projectID)

	topic := client.Topic(topicID)
	_, _ = client.CreateSubscription(
		ctx,
		subscriptionID,
		pubsub.SubscriptionConfig{
			Topic: topic,
			PushConfig: pubsub.PushConfig{
				Endpoint: "http://host.docker.internal:8080/api/v1/pubsub",
			},
		},
	)
	return nil
}

func debugf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
}

func fatalf(format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, os.Args[0]+": "+format+"\n", params...)
}
