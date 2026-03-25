package cmd

import (
	"github.com/edgehook/ithings/webserver"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:     "update process",
	Long:    `iot update process manager.. `,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: To help debugging, immediately log version
		klog.Infof("###########  Start the update process...! ###########")
		registerModules()
		// start all modules
		core.Run()
	},
}

// register all module into beehive.
func registerModules() {
	webserver.Register()
}
