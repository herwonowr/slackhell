package app

import (
	"github.com/herwonowr/slackhell/internal/bot"
	"github.com/herwonowr/slackhell/internal/config"
	"github.com/herwonowr/slackhell/internal/repository"
	"github.com/herwonowr/slackhell/internal/service"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run Slackhell C2",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd)
		},
	}
)

func run(cmd *cobra.Command) error {
	debug := false
	config, err := config.NewConfig(cmd)
	if err != nil {
		return err
	}

	if config.Log.Debug {
		debug = true
	}

	repo, err := repository.NewService(config.Database.Path)
	if err != nil {
		return err
	}

	err = repo.Init()
	if err != nil {
		return err
	}
	defer repo.Close()

	srv := service.NewService(repo)
	botClient := bot.NewService(srv, config.Slack.Token, debug)
	err = botClient.InitBot(databaseVersion, config.Account.ID, config.Account.RealName)
	if err != nil {
		return err
	}

	err = botClient.Listen()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringP("config", "c", "./data/config/slackhell.toml", "slackhell configuration file")
}
