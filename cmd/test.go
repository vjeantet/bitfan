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
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veino/logfan/lib"
	config "github.com/veino/veino/config"
	runtime "github.com/veino/veino/runtime"
)

type location struct {
	path        string
	kind        string
	workingpath string
}

type locations struct {
	items []*location
}

func (l *location) expand() ([]string, error) {
	locs := []string{}
	if fi, err := os.Stat(l.path); err == nil {

		if false == fi.IsDir() {
			locs = append(locs, l.path)
			return locs, nil
		}

		files, err := filepath.Glob(filepath.Join(l.path, "*.*"))
		if err != nil {
			return locs, fmt.Errorf("error %s", err.Error())

		}

		//use each file
		for _, file := range files {
			switch strings.ToLower(filepath.Ext(file)) {
			case ".conf":
				locs = append(locs, file)
				continue
			default:

			}
		}
	} else {
		return locs, fmt.Errorf("%s not found", l.path)
	}
	return locs, nil
}

func (l *locations) expand() []string {
	locs := []string{}

	for _, loc := range l.items {
		if loc.kind == "url" {
			locs = append(locs, loc.path)
			continue
		}

		if fi, err := os.Stat(loc.path); err == nil {

			if false == fi.IsDir() {
				locs = append(locs, loc.path)
				continue
			}

			files, err := filepath.Glob(filepath.Join(loc.path, "*.*"))
			if err != nil {
				log.Printf("error %s", err.Error())
				continue
			}

			//use each file
			for _, file := range files {
				switch strings.ToLower(filepath.Ext(file)) {
				case ".conf":
					locs = append(locs, file)
					continue
				default:
					log.Printf("ignored file %s", file)
				}
			}
		} else {
			log.Println(loc.path, " not found")
		}

	}

	return locs
}

func (l *locations) add(ref string, cwl string) error {
	loc := &location{}
	if v, _ := url.Parse(ref); v.Scheme == "http" || v.Scheme == "https" {
		loc.kind = "url"
		loc.path = ref
	} else if _, err := os.Stat(ref); err == nil {
		loc.kind = "file"
		loc.path, err = filepath.Abs(ref)
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(filepath.Join(cwl, ref)); err == nil {
		loc.kind = "file"
		loc.path = filepath.Join(cwl, ref)
	} else if v, _ := url.Parse(cwl); v.Scheme == "http" || v.Scheme == "https" {
		loc.kind = "url"
		loc.path = cwl + ref
	} else {
		loc.kind = "inline"
		loc.path = ref

		// return fmt.Errorf("unknow location %s -- current working location is %s", ref, cwl)
	}

	loc.workingpath = cwl

	l.items = append(l.items, loc)
	return nil
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test configurations (files, url, directories)",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		var content []byte
		var cwd string
		var ncwl string
		cwd, _ = os.Getwd()

		var locs locations

		// Si il n'y a pas argument alors utilise viper.config
		if len(args) == 0 {
			cwd = filepath.Dir(viper.ConfigFileUsed())
			for _, conf := range viper.GetStringSlice("config") {
				if filepath.IsAbs(conf) {
					locs.add(conf, cwd)
				} else {
					locs.add(filepath.Join(cwd, conf), cwd)
				}
			}
		} else {
			for _, v := range args {
				locs.add(v, cwd)
			}
		}

		// Si location est un dossier
		//   calcul les autres locations
		var cko int
		var ctot int
		for _, loc := range locs.expand() {
			ctot++
			content, ncwl, err = lib.GetContentFromLocation(loc, cwd)
			if err != nil {
				fmt.Printf("error %s\n", err)
			}

			err = testConfigContent(content, ncwl)
			if err != nil {
				fmt.Printf("%s\n -> %s\n\n", loc, err)
				cko++
			}
		}

		if ctot == 0 && len(args) == 1 {
			ctot++
			err = testConfigContent([]byte(args[0]), cwd)
			if err != nil {
				fmt.Printf("-> %s\n\n", err)
				cko++
			}
		}

		if ctot == 0 {
			fmt.Println("No configuration available to test")
		} else if cko == 0 {
			fmt.Printf("Everything is ok, %d configurations checked\n", ctot)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}

func testConfigContent(content []byte, cwl string) error {
	// logError := log.New(os.Stderr, "ERROR: ", 0)
	// logInfo := log.New(os.Stdout, "", 0)
	// logWarning := log.New(os.Stdout, "WARNING: ", 0)

	configAgents, err := lib.ParseConfig("test", content, cwl)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	configAgentsOrdered := config.Sort(configAgents, config.SortOutputsFirst)
	for _, configAgent := range configAgentsOrdered {
		_, err := runtime.NewAgent(configAgent, 0)
		if err != nil {
			// logError.Fatalf("plugin '%s' may not start : %s", configAgent.Type, err.Error())
			return fmt.Errorf("%s", err.Error())
		}
	}

	return nil
}
