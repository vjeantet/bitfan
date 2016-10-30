// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// serviceRestartCmd represents the serviceRestart command
var serviceRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart a logfan service",

	Run: func(cmd *cobra.Command, args []string) {

		servicename, _ := cmd.Flags().GetString("name")
		svcConfig := getServiceConfig()
		svcConfig.Name = servicename
		svcConfig.DisplayName = servicename
		s := getService(svcConfig)

		log.Println("stopping logfan service...")
		err := s.Stop()
		if err != nil {
			log.Printf("sop service error : %s", err)
		} else {
			// log.Println("service logfan stopped")
		}

		log.Println("starting logfan service...")
		err = s.Start()
		if err != nil {
			log.Printf("start service error : %s", err)
		} else {
			// log.Println("service logfan started")
		}

	},
}

func init() {
	serviceCmd.AddCommand(serviceRestartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceRestartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceRestartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
