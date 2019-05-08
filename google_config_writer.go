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
	variables := map[string]string{
		"DataFlowTmpGCSLocation":  w.c.String(googleDataFlowTmpGCSLocationFlag),
		"DataFlowTemplateGCSPath": w.c.String(googleDataFlowTemplateGCSPathFlag),
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
		"variables":         func() map[string]string { return variables },
		"hclvalue":          hclvalue,
		"hclident":          hclident,
		"tfDoNotEditStamp":  func() string { return tfDoNotEditStamp },

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
