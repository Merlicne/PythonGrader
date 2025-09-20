package main

import (
	"context"
	"fmt"
	"log"
	"os"
	mysql "python-runner/MYSQL"
	"python-runner/executer"

	"github.com/urfave/cli/v3"
)

// CLI template for urfave/cli v3
func createCLIApp() *cli.Command {
	return &cli.Command{
		Name:    "python-runner",
		Usage:   "A tool for running Python code and managing test cases",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Execute Python code",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "Python file to execute",
					},
					&cli.StringFlag{
						Name:    "code",
						Aliases: []string{"c"},
						Usage:   "Python code string to execute",
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Enable verbose output",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					file := cmd.String("file")
					code := cmd.String("code")
					verbose := cmd.Bool("verbose")

					if file == "" && code == "" {
						return fmt.Errorf("either --file or --code must be provided")
					}

					if verbose {
						fmt.Println("Running Python code...")
					}

					// Initialize Python executor
					python := &executer.PythonExecutor{}

					if file != "" {
						if verbose {
							fmt.Printf("Executing file: %s\n", file)
						}
						// Add file execution logic here
						fmt.Printf("Would execute file: %s\n", file)
					} else {
						if verbose {
							fmt.Printf("Executing code: %s\n", code)
						}
						// Add code execution logic here
						fmt.Printf("Would execute code: %s\n", code)
					}

					_ = python // Use the executor as needed
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Run test cases",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "question-id",
						Aliases: []string{"q"},
						Usage:   "Question ID to run tests for",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Run all test cases",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					questionID := cmd.Int("question-id")
					all := cmd.Bool("all")

					if !all && questionID == 0 {
						return fmt.Errorf("either --question-id or --all must be provided")
					}

					// Initialize MySQL connection
					err := mysql.InitializeGlobalConnection()
					if err != nil {
						return fmt.Errorf("failed to initialize MySQL connection: %w", err)
					}

					mysqlExecuter := executer.NewMySQLExecuter()

					if all {
						fmt.Println("Running all test cases...")
						// Add logic to run all test cases
					} else {
						fmt.Printf("Running test cases for question ID: %d\n", questionID)
						// Add logic to run specific test cases
						_ = mysqlExecuter // Use the MySQL executor as needed
					}

					return nil
				},
			},
			{
				Name:  "db",
				Usage: "Database operations",
				Commands: []*cli.Command{
					{
						Name:  "connect",
						Usage: "Test database connection",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fmt.Println("Testing database connection...")
							err := mysql.InitializeGlobalConnection()
							if err != nil {
								return fmt.Errorf("failed to connect to database: %w", err)
							}
							fmt.Println("Database connection successful!")
							return nil
						},
					},
					{
						Name:  "query",
						Usage: "Execute a database query",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "sql",
								Aliases:  []string{"s"},
								Usage:    "SQL query to execute",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							sqlQuery := cmd.String("sql")
							fmt.Printf("Executing query: %s\n", sqlQuery)
							// Add query execution logic here
							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug mode",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Load configuration from `FILE`",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Global setup before any command runs
			if cmd.Bool("debug") {
				fmt.Println("Debug mode enabled")
			}

			if configFile := cmd.String("config"); configFile != "" {
				fmt.Printf("Loading config from: %s\n", configFile)
				// Add config loading logic here
			}

			return ctx, nil
		},
		After: func(ctx context.Context, cmd *cli.Command) error {
			// Global cleanup after any command runs
			fmt.Println("Command completed successfully")
			return nil
		},
	}
}

// Example main function using the CLI app
func mainWithCLI() {
	app := createCLIApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// Alternative: If you want to add commands dynamically
func createDynamicCLIApp() *cli.Command {
	app := &cli.Command{
		Name:    "python-runner",
		Usage:   "A tool for running Python code and managing test cases",
		Version: "1.0.0",
	}

	// Add commands dynamically
	app.Commands = append(app.Commands, createRunCommand())
	app.Commands = append(app.Commands, createTestCommand())
	app.Commands = append(app.Commands, createDatabaseCommands())

	return app
}

// Separate command creators for better organization
func createRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Execute Python code",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Python file to execute",
			},
			&cli.StringFlag{
				Name:    "code",
				Aliases: []string{"c"},
				Usage:   "Python code string to execute",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Implementation here
			return nil
		},
	}
}

func createTestCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Run test cases",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "question-id",
				Aliases: []string{"q"},
				Usage:   "Question ID to run tests for",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Implementation here
			return nil
		},
	}
}

func createDatabaseCommands() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Database operations",
		Commands: []*cli.Command{
			{
				Name:  "connect",
				Usage: "Test database connection",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// Implementation here
					return nil
				},
			},
		},
	}
}
