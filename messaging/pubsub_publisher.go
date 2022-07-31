package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/util"
)

type PubsubClient struct {
	log     *util.Logger
	client  *pubsub.Client
	topicID string
}

// NewPubsubPublisher create a new instance of PubsubPublisher
func NewPubsubPublisher(projectID, topicID string, log *util.Logger) (Publisher, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create the topic if it doesn't exist.
	if err := createTopic(ctx, client, topicID); err != nil {
		return nil, err
	}

	return &PubsubClient{client: client, topicID: topicID, log: log}, nil
}

func createTopic(ctx context.Context, client *pubsub.Client, topicID string) error {
	if exists, err := client.Topic(topicID).Exists(ctx); err != nil {
		return err
	} else if !exists {
		if _, err := client.CreateTopic(ctx, topicID); err != nil {
			return err
		}
	}
	return nil
}

func (p *PubsubClient) PublishMerchant(merchant *dto.MerchantPublish) error {
	return p.publish(merchant.MerchantID, merchant)
}

// publish a message with specified name and data (json serializable) to cloud pub/sub
func (p *PubsubClient) publish(merchantID string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx := context.Background()

	topic := p.client.Topic(p.topicID)
	message := &pubsub.Message{
		Data:       jsonData,
		Attributes: map[string]string{"merchantId": merchantID},
	}

	publishedMsg, err := topic.Publish(ctx, message).Get(ctx)
	p.log.Log().Str("data", fmt.Sprintf("%+v", data)).Msgf("[published] msgId=%v, merchantId=%v", publishedMsg, merchantID)
	return err
}
