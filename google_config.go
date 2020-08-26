package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
)

// GoogleCrossProjectSubscription struct represents a cross-project subscription for a Hedwig consumer app
type GoogleCrossProjectSubscription struct {
	Project string `json:"project"`
	Topic   string `json:"topic"`
}

// GooglePullConsumer struct represents a Hedwig consumer app
type GooglePullConsumer struct {
	Queue                     string                           `json:"queue"`
	Subscriptions             []string                         `json:"subscriptions"`
	CrossProjectSubscriptions []GoogleCrossProjectSubscription `json:"cross_project_subscriptions"`
	ServiceAccount            string                           `json:"service_account"`
	Labels                    map[string]string                `json:"labels"`
}

// GoogleTopic struct represents a Hedwig topic
type GoogleTopic struct {
	Name            string   `json:"name"`
	EnableFirehose  bool     `json:"enable_firehose"`
	ServiceAccounts []string `json:"service_accounts"`
}

// GoogleConfig struct represents the Hedwig configuration for Google Cloud
type GoogleConfig struct {
	Topics        []*GoogleTopic        `json:"topics"`
	PullConsumers []*GooglePullConsumer `json:"pull_consumers,omitempty"`
	// TODO
	// PushConsumers  []*PushConsumer `json:"push_consumers,omitempty"`
}

// newGoogleConfig returns a new config read from a file
func newGoogleConfig(filename string) (*GoogleConfig, error) {
	configContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %q", err)
	}
	config := GoogleConfig{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file as JSON object: %q", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

var googleTopicRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
var googleSubscriptionNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
var labelKeyRegex = regexp.MustCompile("^[a-z][a-z0-9-_]*$")
var labelValueRegex = regexp.MustCompile("^[a-z0-9-_]*$")

// Validates that topic names are valid format
func (c *GoogleConfig) validateTopics() error {
	for _, topic := range c.Topics {
		if !googleTopicRegex.MatchString(topic.Name) {
			return fmt.Errorf("invalid topic name: |%s|, must match regex: %s", topic.Name, googleTopicRegex)
		}
	}
	return nil
}

// Validates that consumer queues are valid format
func (c *GoogleConfig) validateQueueConsumers() error {
	for _, consumer := range c.PullConsumers {
		if !googleSubscriptionNameRegex.MatchString(consumer.Queue) {
			return fmt.Errorf(
				"invalid subscription name: |%s|, must match regex: %s", consumer.Queue, googleSubscriptionNameRegex)
		}

		if len(consumer.Subscriptions) == 0 && len(consumer.CrossProjectSubscriptions) == 0 {
			return fmt.Errorf("consumer must contain at least one subscription: '%s'", consumer.Subscriptions)
		}

		for _, subscription := range consumer.Subscriptions {
			// verify that topic was declared
			found := false
			for _, topic := range c.Topics {
				if topic.Name == subscription {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("topic not declared: |%s|", subscription)
			}
		}

		for k, v := range consumer.Labels {
			if !labelKeyRegex.MatchString(k) {
				return fmt.Errorf("invalid label key: |%s|, must match regex: %s", k, labelKeyRegex)
			}
			if !labelValueRegex.MatchString(v) {
				return fmt.Errorf("invalid label value: |%s|, must match regex: %s", v, labelValueRegex)
			}
		}
	}
	return nil
}

// validate verifies that the input configuration is sane
func (c *GoogleConfig) validate() error {
	if err := c.validateTopics(); err != nil {
		return err
	}

	if err := c.validateQueueConsumers(); err != nil {
		return err
	}

	return nil
}
