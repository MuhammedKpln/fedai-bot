package core

import (
	"fmt"
	"log"
	S "muhammedkpln/fedai/shared"
	"os"
	"path"
	"plugin"
	"strings"
)

var LoadedPlugins = map[string]S.Plugin{}

func LoadModules() {
	files, err := os.ReadDir("./pl")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			var fileNameWithoutExt = strings.Split(file.Name(), ".")[0]
			var filePath = path.Join("pl", file.Name())
			// Open the plugin
			plug, err := plugin.Open(filePath)
			if err != nil {
				AppLog().Errorf("Loader: Error loading plugin, deleting ", err.Error())

				DeletePlugin(filePath, file.Name())
				continue
			}

			pluginSym, err := plug.Lookup("Plugin")
			if err != nil {
				AppLog().Errorf("Loader: Error loading plugin, deleting ", err.Error())
				continue
			}

			plugin, plOk := pluginSym.(*S.Plugin)

			if !plOk {
				fmt.Printf("Unexpected plugin for '%s' in plugin %s\n", file.Name(), file.Name())
				continue
			}

			// Store the module with the greeter instance
			LoadedPlugins[fileNameWithoutExt] = *plugin
		}
	}
}
