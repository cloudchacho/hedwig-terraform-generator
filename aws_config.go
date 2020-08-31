package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// AWSSubscription struct represents a subscription for a Hedwig consumer app
type AWSSubscription struct {
	AccountID string `json:"account_id"`
	Topic     string `json:"topic"`
}

func (s *AWSSubscription) UnmarshalJSON(data []byte) error {
	type AWSSubscriptionAlias AWSSubscription
	if err := json.Unmarshal(data, (*AWSSubscriptionAlias)(s)); err != nil {
		var topic string
		if strErr := json.Unmarshal(data, &topic); strErr != nil {
			return err
		}
		s.Topic = topic
	}
	return nil
}

// AWSQueueConsumer struct represents a Hedwig consumer app
type AWSQueueConsumer struct {
	Queue         string            `json:"queue"`
	Tags          map[string]string `json:"tags"`
	Subscriptions []AWSSubscription `json:"subscriptions"`
}

// AWSLambdaConsumer struct represents a Hedwig subscription for a lambda app
type AWSLambdaConsumer struct {
	FunctionARN       string   `json:"function_arn"`
	FunctionName      string   `json:"function_name,omitempty"`
	FunctionQualifier string   `json:"function_qualifier,omitempty"`
	Subscriptions     []string `json:"subscriptions"`
}

var lambdaARNRegexp = regexp.MustCompile(`^arn:aws:lambda:([^:]+):([^:]+):function:([^:]+)(:([^:]+))?$`)

// init initializes the data structure with function name and qualifier if required
func (ls *AWSLambdaConsumer) init() error {
	if ls.FunctionName != "" {
		return nil
	}

	if strings.Contains(ls.FunctionARN, "${") {
		return fmt.Errorf("unable to parse function ARN since it's an interpolated value")
	}

	matches := lambdaARNRegexp.FindStringSubmatch(ls.FunctionARN)
	if len(matches) > 0 {
		if ls.FunctionName == "" {
			ls.FunctionName = matches[3]
		}
		if ls.FunctionQualifier == "" && len(matches) >= 6 {
			ls.FunctionQualifier = matches[5]
		}
	}
	if ls.FunctionName == "" {
		return fmt.Errorf("unable to parse function ARN")
	}
	return nil
}

// AWSConfig struct represents the Hedwig configuration for AWS provider
type AWSConfig struct {
	Topics          []string             `json:"topics"`
	QueueConsumers  []*AWSQueueConsumer  `json:"queue_consumers,omitempty"`
	LambdaConsumers []*AWSLambdaConsumer `json:"lambda_consumers,omitempty"`
}

// newAWSConfig returns a new config read from a file
func newAWSConfig(filename string) (*AWSConfig, error) {
	configContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %q", err)
	}
	config := AWSConfig{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file as JSON object: %q", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	for _, ls := range config.LambdaConsumers {
		err := ls.init()
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

var snsTopicRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
var sqsQueueRegex = regexp.MustCompile(`^[A-Z0-9-]+$`)

// Validates that topic names are valid format
func (c *AWSConfig) validateTopics() error {
	for _, topic := range c.Topics {
		if !snsTopicRegex.MatchString(topic) {
			return fmt.Errorf("invalid topic name, must only contain: [a-z], [0-9], [-]: '%s'", topic)
		}
	}
	return nil
}

// Validates that consumer queues are valid format
func (c *AWSConfig) validateQueueConsumers() error {
	for _, consumer := range c.QueueConsumers {
		if !sqsQueueRegex.MatchString(consumer.Queue) {
			return fmt.Errorf("invalid queue name, must only contain: [A-Z], [0-9], [-]: '%s'", consumer.Queue)
		}

		if len(consumer.Subscriptions) == 0 {
			return fmt.Errorf("consumer must contain at least one subscription: '%s'", consumer.Queue)
		}

		for _, subscription := range consumer.Subscriptions {
			if subscription.AccountID != "" {
				// cross project topics will be created in a different config
				continue
			}
			// verify that topic was declared
			found := false
			for _, topic := range c.Topics {
				if topic == subscription.Topic {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("topic not declared: '%s'", subscription.Topic)
			}
		}
	}
	return nil
}

// Validates that lambda subscriptions refer to valid topics
func (c *AWSConfig) validateLambdaConsumers() error {
	for _, consumer := range c.LambdaConsumers {
		if len(consumer.Subscriptions) == 0 {
			return fmt.Errorf("consumer must contain at least one subscription: '%s'", consumer.FunctionARN)
		}

		for _, subscription := range consumer.Subscriptions {
			// verify that topic was declared
			found := false
			for _, topic := range c.Topics {
				if topic == subscription {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("topic not declared: '%s'", subscription)
			}
		}
	}
	return nil
}

// validate verifies that the input configuration is sane
func (c *AWSConfig) validate() error {
	if err := c.validateTopics(); err != nil {
		return err
	}

	if err := c.validateQueueConsumers(); err != nil {
		return err
	}

	if err := c.validateLambdaConsumers(); err != nil {
		return err
	}

	return nil
}
