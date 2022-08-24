/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/cage1016/alfred-paletter/sqlite"
)

const updateJobName = "checkForUpdate"

var (
	repo = "cage1016/alfred-paletter"
	wf   *aw.Workflow
	cfg  Config
	db   *gorm.DB
)

type Config struct {
	DbConfig sqlite.Config
}

func CheckForUpdate() {
	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")
		cmd := exec.Command(os.Args[0], "update")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	if wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		wf.NewItem("An update is available!").
			Subtitle("⇥ or ↩ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(&aw.Icon{Value: "update-available.png"})
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Paletter",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		CheckForUpdate()
		wf.SendFeedback()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	wf.Run(func() {
		err := envconfig.Process("qs", &cfg)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		db, err = sqlite.Connect(&cfg.DbConfig)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if err := rootCmd.Execute(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	})
}

func init() {
	wf = aw.New(update.GitHub(repo), aw.HelpURL(repo+"/issues"))
	if _, err := os.Stat(fmt.Sprintf("%s/cmd-paletter", wf.DataDir())); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(fmt.Sprintf("%s/cmd-paletter", wf.DataDir()), 0755)
	}
	if _, err := os.Stat(fmt.Sprintf("%s/cmd-copy", wf.DataDir())); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(fmt.Sprintf("%s/cmd-copy", wf.DataDir()), 0755)
	}
	if _, err := os.Stat(fmt.Sprintf("%s/cmd-export", wf.DataDir())); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(fmt.Sprintf("%s/cmd-export", wf.DataDir()), 0755)
	}
	wf.Args() // magic for "workflow:update"
}
