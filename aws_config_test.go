package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAWSSchemaFail(t *testing.T) {
	schema := []byte(`
{
  "consumers": [
    {
      "queue": "dev-myapp",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      }
    }
  ],
  "topics": "not-a-list"
}
`)
	assert.EqualError(t, json.Unmarshal(schema, &AWSConfig{}),
		"json: cannot unmarshal string into Go struct field AWSConfig.topics of type []string")
}

func TestValidateAWSTopic(t *testing.T) {
	invalidTopics := []string{
		"UPPER",
		"under_score",
		"punctuation!",
	}

	config := AWSConfig{}
	for _, topic := range invalidTopics {
		config.Topics = []string{topic}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid topic name, must only contain: [a-z], [0-9], [-]: '%s'", topic),
			"Didn't fail validation for '%s'",
			topic,
		)
	}
}

func TestValidateAWSQueue(t *testing.T) {
	invalidQueues := []string{
		"lower",
		"UNDER_SCORE",
		"PUNCTUATION!",
	}

	config := AWSConfig{}
	for _, queue := range invalidQueues {
		config.QueueConsumers = []*AWSQueueConsumer{{Queue: queue}}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid queue name, must only contain: [A-Z], [0-9], [-]: '%s'", queue),
			"Didn't fail validation for '%s'",
			queue,
		)
	}
}

func TestValidateAWSSubscriptionTopic(t *testing.T) {
	config := AWSConfig{
		QueueConsumers: []*AWSQueueConsumer{{Queue: "QUEUE", Subscriptions: []string{"does-not-exist"}}},
	}
	assert.EqualError(
		t,
		config.validate(),
		"topic not declared: 'does-not-exist'",
		"Didn't fail validation for topic 'does-not-exist'",
	)
}

func TestValidateLambdaSubscriptionTopic(t *testing.T) {
	config := AWSConfig{
		LambdaConsumers: []*AWSLambdaConsumer{{FunctionARN: "function", Subscriptions: []string{"does-not-exist"}}},
	}
	assert.EqualError(
		t,
		config.validate(),
		"topic not declared: 'does-not-exist'",
		"Didn't fail validation for topic 'does-not-exist'",
	)
}

func TestValidAWSJSON(t *testing.T) {
	var validConfig = []byte(`{
  "topics": [
    "my-topic"
  ],
  "queue_consumers": [
    {
      "queue": "DEV-MYAPP",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      },
      "subscriptions": [
        "my-topic"
      ]
    }
  ],
  "lambda_consumers": [
    {
      "function_arn": "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
      "subscriptions": [
        "my-topic"
      ]
    }
  ]
}`)

	var validConfigObj = AWSConfig{
		QueueConsumers: []*AWSQueueConsumer{
			{
				"DEV-MYAPP",
				map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
				[]string{"my-topic"},
			},
		},
		LambdaConsumers: []*AWSLambdaConsumer{
			{
				FunctionARN:   "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
				Subscriptions: []string{"my-topic"},
			},
		},
		Topics: []string{"my-topic"},
	}

	config := AWSConfig{}
	json.Unmarshal(validConfig, &config)
	assert.Equal(t, validConfigObj, config)
}

func TestValidJSONNoLambda(t *testing.T) {
	var validConfig = []byte(`{
  "topics": [
    "my-topic"
  ],
  "queue_consumers": [
    {
      "queue": "DEV-MYAPP",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      },
      "subscriptions": [
        "my-topic"
      ]
    }
  ]
}`)

	var validConfigObj = AWSConfig{
		QueueConsumers: []*AWSQueueConsumer{
			{
				"DEV-MYAPP",
				map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
				[]string{"my-topic"},
			},
		},
		Topics: []string{"my-topic"},
	}

	config := AWSConfig{}
	json.Unmarshal(validConfig, &config)
	assert.Equal(t, validConfigObj, config)
}

func TestAWSValidNoConsumers(t *testing.T) {
	var validConfig = []byte(`{
  "topics": [
    "my-topic"
  ],
  "lambda_consumers": [
    {
      "function_arn": "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
      "subscriptions": ["my-topic"]
    }
  ]
}`)

	var validConfigObj = AWSConfig{
		LambdaConsumers: []*AWSLambdaConsumer{
			{
				FunctionARN:   "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
				Subscriptions: []string{"my-topic"},
			},
		},
		Topics: []string{"my-topic"},
	}

	config := AWSConfig{}
	json.Unmarshal(validConfig, &config)
	assert.Equal(t, validConfigObj, config)
}

func TestLambdaConsumer_Init(t *testing.T) {
	ls := AWSLambdaConsumer{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "deployed", ls.FunctionQualifier)
}

func TestLambdaConsumer_Init_Fail(t *testing.T) {
	ls := AWSLambdaConsumer{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:foo:myFunction:deployed",
	}
	assert.Error(t, ls.init(), "unable to parse function ARN")
}

func TestLambdaConsumer_Init_NoQualifier(t *testing.T) {
	ls := AWSLambdaConsumer{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "", ls.FunctionQualifier)
}

func TestLambdaConsumer_Init_Interpolated(t *testing.T) {
	ls := AWSLambdaConsumer{
		FunctionARN:  "${aws_lambda_function.myFunction.arn}",
		FunctionName: "myFunction",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "", ls.FunctionQualifier)
}

func TestLambdaConsumer_Init_InterpolatedFail(t *testing.T) {
	ls := AWSLambdaConsumer{
		FunctionARN: "${aws_lambda_function.myFunction.arn}",
	}
	assert.Error(t, ls.init(), "unable to parse function ARN since it's an interpolated value")
}
