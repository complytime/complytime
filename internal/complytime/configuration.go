// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/config"
	"github.com/oscal-compass/oscal-sdk-go/generators"
)

var (
	complyTimeDirectory        = filepath.Join(xdg.ConfigHome, "complytime")
	complyTimePluginsDirectory = filepath.Join(complyTimeDirectory, "plugins")
	complyTimeDirectoryBundles = filepath.Join(complyTimeDirectory, "bundles")
)

const compDefSuffix = "component-definition.json"

func FindComponentDefinitions() ([]oscalTypes.ComponentDefinition, error) {
	items, err := os.ReadDir(complyTimeDirectoryBundles)
	if err != nil {
		return nil, err
	}

	var compDefBundles []oscalTypes.ComponentDefinition
	for _, item := range items {

		if !strings.HasSuffix(item.Name(), compDefSuffix) {
			continue
		}
		compDefPath := filepath.Join(complyTimeDirectoryBundles, item.Name())
		compDefPath = filepath.Clean(compDefPath)
		file, err := os.Open(compDefPath)
		if err != nil {
			return nil, err
		}
		definition, err := generators.NewComponentDefinition(file)
		if err != nil {
			return nil, err
		}
		if definition == nil {
			return nil, errors.New("invalid component definitions")
		}
		compDefBundles = append(compDefBundles, *definition)
	}
	return compDefBundles, nil
}

func Config() (*config.C2PConfig, error) {
	cfg := config.DefaultConfig()
	cfg.PluginDir = complyTimePluginsDirectory
	compDefBundles, err := FindComponentDefinitions()
	if err != nil {
		return cfg, err
	}
	cfg.ComponentDefinitions = compDefBundles
	return cfg, nil
}
