package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

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
		PullConsumers: []*GooglePullConsumer{{Queue: "myapp", Subscriptions: []string{"does-not-exist"}}},
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
			{Queue: "myapp", Subscriptions: []string{"topic"}, Labels: map[string]string{"UPPER": ""}},
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
      "name": "topic"
	}
  ],
  "pull_consumers": [
    {
      "queue": "dev-myapp",
      "labels": {
        "app": "myapp",
        "env": "dev"
      },
      "subscriptions": [
        "topic"
      ]
    }
  ]
}`)

	var validConfigObj = GoogleConfig{
		PullConsumers: []*GooglePullConsumer{
			{
				"dev-myapp",
				[]string{"topic"},
				map[string]string{
					"app": "myapp",
					"env": "dev",
				},
			},
		},
		Topics: []*GoogleTopic{{Name: "topic"}},
	}

	config := GoogleConfig{}
	err := json.Unmarshal(validConfig, &config)
	require.NoError(t, err)
	assert.Equal(t, validConfigObj, config)
}
