// Copyright © 2016 Valere JEANTET <valere.jeantet@gmail.com>
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
	"os"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/veino/runtime/metrics"
)

var (
	configPath, configString                            string
	logPath                                             string
	prometheusListen, prometheusPath                    string
	webhookListen                                       string
	filterworkers                                       int
	verbose, debug                                      bool
	version, configtest, prometheus                     bool
	installWindowsServiceFlag, removeWindowsServiceFlag bool
	stats                                               metrics.IStats
)

func parseFlags() {
	pflag.StringVarP(&configPath, "config", "f", "", "Load the Logstash config from a file or directory")
	pflag.StringVarP(&configString, "eval", "e", "", "Use the given string as the configuration data.")
	pflag.IntVarP(&filterworkers, "filterworkers", "w", runtime.NumCPU(), "number of workers")
	pflag.StringVarP(&logPath, "log", "l", "", "Log to a given path. Default is to log to stdout.")
	pflag.BoolVarP(&verbose, "verbose", "", false, "Increase verbosity to the first level (info), less verbose.")
	pflag.BoolVarP(&debug, "debug", "", false, "Increase verbosity to the last level (trace), more verbose.")
	pflag.BoolVarP(&prometheus, "prometheus", "", false, "Export stats using prometheus output")
	pflag.StringVarP(&prometheusListen, "prometheus.listen", "", "0.0.0.0:24232", "Address and port to bind Prometheus metrics")
	pflag.StringVarP(&prometheusPath, "prometheus.path", "", "/metrics", "Expose Prometheus metrics at specified path.")
	pflag.StringVarP(&webhookListen, "webhook.listen", "", "0.0.0.0:19090", "Address and port to bind webhooks")
	pflag.BoolVarP(&version, "version", "V", false, "Display the version of Logstash.")
	pflag.BoolVarP(&configtest, "configtest", "t", false, "Test config file or directory")
	if runtime.GOOS == "windows" {
		pflag.BoolVar(&installWindowsServiceFlag, "install-windows-service", false, "Install logfan as Windows service")
		pflag.BoolVar(&removeWindowsServiceFlag, "remove-windows-service", false, "Remove logfan Windows service")
	}

	pflag.Parse()
	if version {
		printVersion()
		os.Exit(0)
	}
	if configtest {
		testConfig(configPath, configString)
		os.Exit(0)
	}
	if prometheus {
		stats = metrics.NewPrometheus(prometheusListen, prometheusPath)
	} else {
		stats = &metrics.StatsVoid{}
	}
	
}
