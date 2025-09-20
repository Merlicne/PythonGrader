package main

// urfave/cli v3

import (
	"context"
	"python-runner/configuration"
	"python-runner/service"

	"github.com/urfave/cli/v3"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:    configuration.GetRequiredEnv("APP_NAME"),
		Usage:   "python grader",
		Version: configuration.GetRequiredEnv("APP_VERSION"),
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run python code",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "python file to run",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					file := cmd.String("file")
					return service.GradeFileByOldId(ctx, file)
				},
			},
			{
				Name:  "run-csv",
				Usage: "read csv file for ids to run",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "csvfile",
						Aliases: []string{"f"},
						Usage:   "csv file to read ids from",
					},
					&cli.StringFlag{
						Name:    "latestVersionDir",
						Aliases: []string{"l"},
						Usage:   "directory containing latest version files",
					},
					&cli.StringFlag{
						Name:    "olderVersionDir",
						Aliases: []string{"o"},
						Usage:   "directory containing older version files",
					},
					&cli.IntFlag{
						Name:    "workers",
						Aliases: []string{"w"},
						Usage:   "maximum number of concurrent workers (default: 4)",
						Value:   4,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					csvfile := cmd.String("csvfile")
					latestVersionDir := cmd.String("latestVersionDir")
					olderVersionDir := cmd.String("olderVersionDir")
					workers := cmd.Int("workers")
					return service.GradeFilesFromIdsCSVWithWorkers(csvfile, latestVersionDir, olderVersionDir, workers)
				},
			},
		},
	}
}
