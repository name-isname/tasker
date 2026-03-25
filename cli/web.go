package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"taskctl/web"
)

var webPort string

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting web server on port %s...\n", webPort)
		return web.StartServer(webPort)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVar(&webPort, "port", "8080", "Port for the web server")
}
