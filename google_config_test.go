package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaFail(t *testing.T) {
	schema := []byte(`
{
  "consumers": [
    {
      "name": "dev-myapp",
      "labels": {
        "App": "myapp",
        "Env": "dev"
      },
      "subscriptions": [
        "mytopic"
      ]
    }
  ],
  "topics": "not-a-list"
}
`)
	assert.EqualError(t, json.Unmarshal(schema, &GoogleConfig{}),
		"json: cannot unmarshal string into Go struct field GoogleConfig.topics of type []*main.GoogleTopic")
}

func TestValidateTopic(t *testing.T) {
	invalidTopics := []string{
		"UPPER",
		"under_score",
		"punctuation!",
	}

	config := GoogleConfig{}
	for _, topic := range invalidTopics {
		config.Topics = []*GoogleTopic{{Name: topic}}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid topic name: |%s|, must match regex: %s", topic, googleTopicRegex),
			"Didn't fail validation for '%s'",
			topic,
		)
	}
}

func TestValidateName(t *testing.T) {
	invalidQueues := []string{
		"UPPER",
		"UNDER_SCORE",
		"PUNCTUATION!",
	}

	config := GoogleConfig{}
	for _, queue := range invalidQueues {
		config.PullConsumers = []*GooglePullConsumer{{Queue: queue}}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid subscription name: |%s|, must match regex: %s", queue, googleSubscriptionNameRegex),
			"Didn't fail validation for '%s'",
			queue,
		)
	}
}

func TestValidateSubscriptionTopic(t *testing.T) {
	config := GoogleConfig{
		PullConsumers: []*GooglePullConsumer{{Queue: "myapp", Subscriptions: []GoogleSubscription{{Topic: "does-not-exist"}}}},
	}
	assert.EqualError(
		t,
		config.validate(),
		"topic not declared: |does-not-exist|",
	)
}

func TestValidateSubscriptionLabel(t *testing.T) {
	config := GoogleConfig{
		PullConsumers: []*GooglePullConsumer{
			{Queue: "myapp", Subscriptions: []GoogleSubscription{{Topic: "topic"}}, Labels: map[string]string{"UPPER": ""}},
		},
		Topics: []*GoogleTopic{{Name: "topic"}},
	}
	assert.EqualError(
		t,
		config.validate(),
		"invalid label key: |UPPER|, must match regex: ^[a-z][a-z0-9-_]*$",
	)
}

func TestValidJSON(t *testing.T) {
	var validConfig = []byte(`{
  "topics": [
	{
      "name": "topic",
      "service_accounts": [
        "myapp@project.iam.gserviceaccount.com"
      ]
	}
  ],
  "pull_consumers": [
    {
      "queue": "dev-myapp",
      "service_account": "myapp@project.iam.gserviceaccount.com",
      "labels": {
        "app": "myapp",
        "env": "dev"
      },
      "subscriptions": [
        "topic",
        {
          "project": "other-project",
          "topic": "topic2"
        },
        {
          "topic": "topic3",
          "enable_ordering": true
        }
      ],
      "high_message_count_threshold": 10000
    }
  ]
}`)

	var validConfigObj = GoogleConfig{
		PullConsumers: []*GooglePullConsumer{
			{
				"dev-myapp",
				[]GoogleSubscription{{Topic: "topic"}, {Project: "other-project", Topic: "topic2"}, {EnableOrdering: true, Topic: "topic3"}},
				"myapp@project.iam.gserviceaccount.com",
				map[string]string{
					"app": "myapp",
					"env": "dev",
				},
				10000,
			},
		},
		Topics: []*GoogleTopic{{"topic", false, []string{"myapp@project.iam.gserviceaccount.com"}}},
	}

	config := GoogleConfig{}
	err := json.Unmarshal(validConfig, &config)
	require.NoError(t, err)
	assert.Equal(t, validConfigObj, config)
}
