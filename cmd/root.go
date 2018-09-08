// Copyright © 2018 Devin Pastoor <devin.pastoor@gmail.com>
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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/metrumresearchgroup/rsq/configlib"
	"github.com/spf13/cobra"
)

// VERSION is the current pkc version
const VERSION string = "0.0.0-1"

var log *logrus.Logger
var fs afero.Fs

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rsq",
	Short: "package check functionality",
	Long:  fmt.Sprintf("rsq cli version %s", VERSION),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, setGlobals)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/rsq.yml)")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	RootCmd.PersistentFlags().String("loglevel", "", "level for logging")
	viper.BindPFlag("loglevel", RootCmd.PersistentFlags().Lookup("loglevel"))

	RootCmd.PersistentFlags().String("libpaths", "", "library paths, colon separated list")
	viper.BindPFlag("libpaths", RootCmd.PersistentFlags().Lookup("libpaths"))


	RootCmd.PersistentFlags().Bool("debug", false, "use debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// globals

}

func setGlobals() {

	fs = afero.NewOsFs()

	log = logrus.New()

	switch logLevel := viper.GetString("loglevel"); logLevel {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	default:
		log.Level = logrus.InfoLevel
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" { // enable ability to specify config file via flag
	// 	viper.SetConfigFile(cfgFile)
	// }
	if viper.GetString("config") == "" {
		_ = configlib.LoadGlobalConfig("rsq")
	} else {
		_ = configlib.LoadConfigFromPath(viper.GetString("config"))
	}
	if viper.GetBool("debug") {
		viper.Debug()
	}
}

func logIfNonZero(lg *logrus.Logger, s []string, lvl string) error {
	lf := lg.Errorf
	switch logLevel := lvl; logLevel {
	case "debug":
		lf = lg.Debugf
	case "info":
		lf = lg.Infof
	case "warn":
		lf = lg.Warnf
	case "error":
		lf = lg.Errorf
	case "fatal":
		lf = lg.Fatalf
	case "panic":
		lf = lg.Panicf
	default:
	}
	if s == nil || len(s) == 0 {
		return nil
	}
	for _, v := range s {
		lf("%v", v)
	}
	return nil
}
