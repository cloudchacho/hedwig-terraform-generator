package main

import (
	"path/filepath"
	"text/template"

	"gopkg.in/urfave/cli.v1"
)

type awsConfigWriter struct {
	c      *cli.Context
	config *AWSConfig
}

func (w *awsConfigWriter) shouldSkipFile(file string) bool {
	switch file {
	case queuesFile:
		return len(w.config.QueueConsumers) == 0
	case variablesFile:
		return len(w.config.QueueConsumers) == 0
	default:
		return false
	}
}

func (w *awsConfigWriter) initTemplates() (*template.Template, error) {
	actions := map[string][]string{
		"QueueAlertAlarmActions": w.c.StringSlice(queueAlertAlarmActionsFlag),
		"QueueAlertOKActions":    w.c.StringSlice(queueAlertOKActionsFlag),
		"DLQAlertAlarmActions":   w.c.StringSlice(dlqAlertAlarmActionsFlag),
		"DLQAlertOKActions":      w.c.StringSlice(dlqAlertOKActionsFlag),
	}
	variables := map[string]string{
		"AwsAccountID": w.c.String(awsAccountIDFlag),
		"AwsRegion":    w.c.String(awsRegionFlag),
	}

	files := []string{
		queuesTmplFile,
		subscriptionsTmplFile,
		topicsTmplFile,
		variablesTmplFile,
	}
	templates := template.New(files[0]) // need an arbitrary name
	templates = templates.Funcs(template.FuncMap{
		"generator_version": func() string { return VERSION },
		"actions":           func() map[string][]string { return actions },
		"variables":         func() map[string]string { return variables },
		"hclvalue":          hclvalue,
		"hclident":          hclident,
		"tfDoNotEditStamp":  func() string { return tfDoNotEditStamp },
		"alerting":          func() bool { return w.c.Bool(alertingFlag) },

		"TFAWSQueueModuleVersion":              func() string { return TFAWSQueueModuleVersion },
		"TFAWSQueueSubscriptionModuleVersion":  func() string { return TFAWSQueueSubscriptionModuleVersion },
		"TFAWSLambdaSubscriptionModuleVersion": func() string { return TFAWSLambdaSubscriptionModuleVersion },
		"TFAWSTopicModuleVersion":              func() string { return TFAWSTopicModuleVersion },
	})

	for _, name := range files {
		_, err := templates.New(name).Parse(string(MustAsset(filepath.Join("aws", name))))
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}
