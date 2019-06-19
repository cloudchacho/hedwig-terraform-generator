package main

import (
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"text/template"
)

type googleConfigWriter struct {
	c      *cli.Context
	config *GoogleConfig
}

func (w *googleConfigWriter) shouldSkipFile(file string) bool {
	return false
}

func (w *googleConfigWriter) initTemplates() (*template.Template, error) {
	channels := map[string][]string{
		"QueueAlertNotificationChannels": w.c.StringSlice(queueAlertNotificationChannelsFlag),
		"DLQAlertNotificationChannels":   w.c.StringSlice(dlqAlertNotificationChannelsFlag),
	}
	variables := map[string]string{
		"DataflowTmpGCSLocation":                 w.c.String(dataflowTmpGCSLocationFlag),
		"DataflowPubSubToPubSubTemplateGCSPath":  w.c.String(dataflowPubSubToPubSubTemplateGCSPathFlag),
		"DataflowPubSubToStorageTemplateGCSPath": w.c.String(dataflowPubSubToStorageGCSPathFlag),
		"DataflowZone":                           w.c.String(googleDataflowZoneFlag),
		"DataflowRegion":                         w.c.String(googleDataflowRegionFlag),
		"DataflowOutputDirectory":                w.c.String(googleFirehoseDataflowOutputDirectoryFlag),
	}
	dataflowEnabled := w.c.Bool(enableFirehoseAllTopics)
	if !dataflowEnabled {
		for _, topic := range w.config.Topics {
			if topic.EnableFirehose {
				dataflowEnabled = true
				break
			}
		}
	}
	flags := map[string]bool{
		"EnableFirehoseAllTopics": w.c.Bool(enableFirehoseAllTopics),
		"DataflowEnabled":         dataflowEnabled,
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
		"channels":          func() map[string][]string { return channels },
		"variables":         func() map[string]string { return variables },
		"flags":             func() map[string]bool { return flags },
		"hclvalue":          hclvalue,
		"hclident":          hclident,
		"tfDoNotEditStamp":  func() string { return tfDoNotEditStamp },
		"alerting":          func() bool { return w.c.Bool(alertingFlag) },

		"TFGoogleSubscriptionModuleVersion": func() string { return TFGoogleSubscriptionModuleVersion },
		"TFGoogleTopicModuleVersion":        func() string { return TFGoogleTopicModuleVersion },
		"TFGoogleQueueModuleVersion":        func() string { return TFGoogleQueueModuleVersion },
	})

	for _, name := range files {
		_, err := templates.New(name).Parse(string(MustAsset(filepath.Join("google", name))))
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}
