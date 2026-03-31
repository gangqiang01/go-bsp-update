package cmd

import (
	"github.com/edgehook/ithings/webserver"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:     "updateBsp",
	Long:    `iot update bsp manager.. `,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: To help debugging, immediately log version
		klog.Infof("###########  Start the updateBsp...! ###########")
		registerModules()
		// start all modules
		core.Run()
	},
}

func init() {
}

// register all module into beehive.
func registerModules() {
	webserver.Register()
}
