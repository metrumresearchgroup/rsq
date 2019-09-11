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
	"os"
	"path/filepath"

	"github.com/metrumresearchgroup/rsq/server/db"
	"github.com/metrumresearchgroup/rsq/server/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start an rsq instance",
	Long: `
	start
 `,
	RunE: start,
}

func start(cmd *cobra.Command, args []string) error {
	client := db.NewClient()
	dbpath, _ := os.Getwd()
	if viper.GetString("dbpath") != "" {
		dbpath = viper.GetString("dbpath")
	}
	badgerPath := filepath.Join(dbpath, "badger")
	if _, err := os.Stat(badgerPath); os.IsNotExist(err) {
		err := fs.MkdirAll(badgerPath, 0755)
		if err != nil {
			log.Fatalf("could not create folder for db %v", err)
			os.Exit(1)
		}
	}
	client.Path = badgerPath
	err := client.Open()
	defer client.Close()
	if err != nil {
		log.Fatalf("could not open db %v", err)
		os.Exit(1)
	}

	js := client.JobService()
	// fmt.Println("about to set job")
	// fmt.Println(testJob)
	httpserver.NewHTTPServer(js, VERSION, viper.GetString("port"), viper.GetInt("workers"), log)
	return nil
}

func init() {
	startCmd.PersistentFlags().Int("workers", 0, "number of workers to execute with")
	viper.BindPFlag("workers", startCmd.PersistentFlags().Lookup("workers"))
	startCmd.PersistentFlags().Int("memory", 0, "total memory available for calculating job queuing")
	viper.BindPFlag("memory", startCmd.PersistentFlags().Lookup("memory"))
	startCmd.PersistentFlags().String("port", "", "port number")
	viper.BindPFlag("port", startCmd.PersistentFlags().Lookup("port"))
	startCmd.PersistentFlags().String("dbpath", "", "database path")
	viper.BindPFlag("dbpath", startCmd.PersistentFlags().Lookup("dbpath"))
	RootCmd.AddCommand(startCmd)

}
