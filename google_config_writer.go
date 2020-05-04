package main

import (
	"path/filepath"
	"text/template"

	"gopkg.in/urfave/cli.v1"
)

type googleConfigWriter struct {
	c      *cli.Context
	config *GoogleConfig
}

func (w *googleConfigWriter) shouldSkipFile(file string) bool {
	switch file {
	case alertsFile:
		return !w.c.Bool(alertingFlag)
	default:
		return false
	}
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
		"GoogleProjectAlerting":                  w.c.String(googleProjectAlerting),
	}
	enableDataflow := w.c.Bool(enableFirehoseAllTopics)
	if !enableDataflow {
		for _, topic := range w.config.Topics {
			if topic.EnableFirehose {
				enableDataflow = true
				break
			}
		}
	}
	flags := map[string]bool{
		"EnableFirehoseAllTopics": w.c.Bool(enableFirehoseAllTopics),
		"EnableDataflow":          enableDataflow,
		"EnableAlerting":          w.c.Bool(alertingFlag),
	}
	files := []string{
		alertsTmplFile,
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
		"hclvalue":          hclvalueV2,
		"hclident":          hclident,
		"tfDoNotEditStamp":  func() string { return tfDoNotEditStamp },

		"TFGoogleSubscriptionModuleVersion": func() string { return TFGoogleSubscriptionModuleVersion },
		"TFGoogleTopicModuleVersion":        func() string { return TFGoogleTopicModuleVersion },
		"TFGoogleAlertsModuleVersion":       func() string { return TFGoogleAlertsModuleVersion },
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
