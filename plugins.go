package main

import (
	"bytes"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// stores information about a group of plugins (e.g. inputs, outputs).
type Plugins struct {
	Type        string
	Plugins     []string
	AllFileData []byte
}

// get a unique set of plugins to prevent duplicates.
func (p *Plugins) uniquePluginList() {
	occurred := map[string]bool{}
	uniq := []string{}

	for e := range p.Plugins {
		if !occurred[p.Plugins[e]] {
			occurred[p.Plugins[e]] = true
			uniq = append(uniq, p.Plugins[e])
		}
	}

	p.Plugins = uniq
}

// generate a string version of the all.go file for this plugin type.
func (p *Plugins) generateTemplate() error {
	t, err := template.ParseFiles("all.go.tmpl")
	if err != nil {
		return errors.Wrap(err, "failed to parse template file")
	}

	var tmpl bytes.Buffer
	err = t.Execute(&tmpl, map[string]interface{}{
		"Type":    p.Type,
		"Plugins": p.Plugins,
	})
	if err != nil {
		return errors.Wrap(err, "failed to generate template")
	}

	p.AllFileData = tmpl.Bytes()

	return nil
}

// backup all.go for this plugin type and write custom version.
func (p *Plugins) WritePluginList() error {
	p.uniquePluginList()
	log.Infof("%11s: %s", p.Type, p.Plugins)

	if err := p.generateTemplate(); err != nil {
		return err
	}

	originalFile := path.Join(sourceDirectory, "plugins", p.Type, "all/all.go")
	backupFile := path.Join(sourceDirectory, "plugins", p.Type, "all/all.go.bak")

	if err := os.Rename(originalFile, backupFile); err != nil {
		return errors.Wrap(err, "failed to backup all.go")
	}

	if err := os.WriteFile(originalFile, p.AllFileData, 0o644); err != nil {
		return errors.Wrap(err, "failed to write new all.go")
	}

	return nil
}

// restore the default all.go file.
func (p *Plugins) RestorePluginList() {
	originalFile := path.Join(sourceDirectory, "plugins", p.Type, "all/all.go")
	backupFile := path.Join(sourceDirectory, "plugins", p.Type, "all/all.go.bak")
	_ = os.Rename(backupFile, originalFile)
}
