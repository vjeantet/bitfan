// +build !windows

package cmd

func Execute() {
	parseFlags()

	if configPath != "" || configString != "" {
		startLogfanAndWait(configPath, configString, stats)
	}

	pflag.Usage()
}

func installWindowsService(serviceName string) error {
	return nil
}

func removeWindowsService(serviceName string) error {
	return nil
}

