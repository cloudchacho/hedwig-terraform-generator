package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"text/template"
)

const (
	queuesFile        = "queues.tf"
	subscriptionsFile = "subscriptions.tf"
	topicsFile        = "topics.tf"
	variablesFile     = "variables.tf"
)

const (
	queuesTmplFile        = "queues.tf.tmpl"
	subscriptionsTmplFile = "subscriptions.tf.tmpl"
	topicsTmplFile        = "topics.tf.tmpl"
	variablesTmplFile     = "variables.tf.tmpl"
)

type configWriterImpl interface {
	initTemplates() (*template.Template, error)
	shouldSkipFile(file string) bool
}

type configWriter struct {
	c                *cli.Context
	configWriterImpl configWriterImpl
	config           interface{}
}

func newConfigWriter(c *cli.Context, config interface{}) *configWriter {
	var writerImpl configWriterImpl
	if c.GlobalString(cloudProviderFlag) == cloudProviderGoogle {
		writerImpl = &googleConfigWriter{c, config.(*GoogleConfig)}
	} else if c.GlobalString(cloudProviderFlag) == cloudProviderAWS {
		writerImpl = &awsConfigWriter{c, config.(*AWSConfig)}
	} else {
		return nil
	}
	return &configWriter{
		c,
		writerImpl,
		config,
	}
}

func (w *configWriter) writeFiles(module string, templates *template.Template) error {
	files := []string{queuesFile, subscriptionsFile, topicsFile, variablesFile}

	for _, file := range files {
		if w.configWriterImpl.shouldSkipFile(file) {
			continue
		}
		path := filepath.Join(module, file)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		templateName := fmt.Sprintf("%s.tmpl", file)
		err = templates.ExecuteTemplate(f, templateName, w.config)
		if err != nil {
			return err
		}
		errClose := f.Close()
		if errClose != nil {
			return errClose
		}
	}
	return nil
}

func (w *configWriter) writeTerraform() error {
	module := w.c.String(moduleFlag)

	if len(module) == 0 {
		return fmt.Errorf("invalid module")
	}

	if err := os.Mkdir(module, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("module directory already exists: %q", module)
		}
		return fmt.Errorf("unable to create module dir: %q, error: %q", module, err)
	}

	templates, err := w.configWriterImpl.initTemplates()
	if err != nil {
		return errors.Wrap(err, "unable to initialize templates")
	}

	if err := w.writeFiles(module, templates); err != nil {
		return err
	}

	if err := hclFmtDir(module); err != nil {
		return err
	}

	return nil
}
