package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sourceDirectory string

	buildConfigFiles []string

	buildAggregators []string
	buildInputs      []string
	buildOutputs     []string
	buildProcessors  []string
)

// build Telegraf with only the specified plugins.
func build(cmd *cobra.Command, args []string) error {
	if len(buildConfigFiles) > 0 {
		if err := parseConfigFiles(); err != nil {
			return err
		}
	}

	pluginGroups := []*Plugins{
		{Type: "inputs", Plugins: buildInputs},
		{Type: "outputs", Plugins: buildOutputs},
		{Type: "processors", Plugins: buildProcessors},
		{Type: "aggregators", Plugins: buildAggregators},
	}

	for _, group := range pluginGroups {
		if err := group.WritePluginList(); err != nil {
			return err
		}

		defer group.RestorePluginList()
	}

	makeCmd := exec.Command("make")
	makeCmd.Dir = sourceDirectory

	out, err := makeCmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return errors.Wrap(err, "failed to build Telegraf")
	}

	return nil
}

// loop through any config files and parse them for plugins.
func parseConfigFiles() error {
	log.Println("config file(s):", buildConfigFiles)

	for _, filename := range buildConfigFiles {
		if err := readConfigFile(filename); err != nil {
			return err
		}
	}

	return nil
}

// open TOML config file, extract declared plugins, and append to plugin lists.
func readConfigFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "config file")
	}

	tomlTree, err := toml.Load(string(content))
	if err != nil {
		return errors.Wrap(err, "config file toml")
	}

	aggregators := tomlTree.Get("aggregators")
	if aggregators != nil {
		buildAggregators = append(buildAggregators, aggregators.(*toml.Tree).Keys()...)
	}

	inputs := tomlTree.Get("inputs")
	if inputs != nil {
		buildInputs = append(buildInputs, inputs.(*toml.Tree).Keys()...)
	}

	outputs := tomlTree.Get("outputs")
	if outputs != nil {
		buildOutputs = append(buildOutputs, outputs.(*toml.Tree).Keys()...)
	}

	processors := tomlTree.Get("processors")
	if processors != nil {
		buildProcessors = append(buildProcessors, processors.(*toml.Tree).Keys()...)
	}

	return nil
}
