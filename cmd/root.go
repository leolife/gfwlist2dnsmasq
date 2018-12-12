package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gfwlist2dnsmasq",
	Short: "Convert gfwlist to dnsmasq.conf with ipset supported",
	Long:  `Convert gfwlist to dnsmasq.conf with ipset supported inspired by mysqto/gfwlist2dnsmasq`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

// flags
var (
	inputFlag    string
	userRuleFlag string
	outputFlag   string
	serverFlag   string
	portFlag     string
	ipsetFlag    string
)

func init() {
	rootCmd.Flags().StringVarP(&inputFlag, "input", "i", "", "path to gfwlist file, raw url or local path, default will get from gfwlist github repo")
	rootCmd.Flags().StringVarP(&userRuleFlag, "user-rule", "u", "", "customized user rule, which will be append to gfwlist raw url or local path")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "path to output dnsmasq.conf, default will write to dnsmasq.gfwlist.conf in current directory")
	rootCmd.Flags().StringVarP(&serverFlag, "server", "s", "", "The upstream dns server address for the poisoned domain, default value is 127.0.0.1")
	rootCmd.Flags().StringVarP(&portFlag, "port", "p", "", "The upstream dns server port for the poisoned domain, default value is 5353")
	rootCmd.Flags().StringVarP(&ipsetFlag, "ipset", "e", "", "ipset name of the dnsmasq ipset, ipset not support if not presented")
}

func run() {
	content, err := GetInputContent(inputFlag)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if userRuleFlag != "" {
		user, err := GetInputContent(userRuleFlag)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		content = io.MultiReader(content, user)
	}

	list, err := ParseList(content)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	generate(list)
}

// Execute run the root command
func Execute() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
