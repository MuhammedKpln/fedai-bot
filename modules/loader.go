package modules

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

			// Open the plugin
			plug, err := plugin.Open(path.Join("pl", file.Name()))
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			pluginSym, err := plug.Lookup("Plugin")
			if err != nil {
				fmt.Printf("Failed to find 'Run' in plugin %s: %v\n", file.Name(), err)
				continue
			}

			plugin, plOk := pluginSym.(*S.Plugin)

			// fmt.Printf("Type of symGreeter: %T\n", plugin)
			// fmt.Printf("Type of symGreeter: %T\n", command)
			// fmt.Printf("Type of symGreeter: %T\n", commandInfo)

			if !plOk {
				fmt.Printf("Unexpected plugin for '%s' in plugin %s\n", file.Name(), file.Name())
				continue
			}

			// Store the module with the greeter instance
			LoadedPlugins[fileNameWithoutExt] = *plugin
		}
	}
}
