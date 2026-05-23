package main

import (
	"fmt"
	"os"

	"github.com/envlint/envlint/reporter"
	"github.com/envlint/envlint/schema"
	"github.com/envlint/envlint/validator"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: envlint <schema.yaml> <.env file>")
		os.Exit(1)
	}

	schemaPath := args[0]
	envPath := args[1]

	s, err := schema.Load(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading schema: %v\n", err)
		os.Exit(1)
	}

	envVars, err := godotenv.Read(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .env file: %v\n", err)
		os.Exit(1)
	}

	results := validator.Validate(s, envVars)
	reporter.Print(results, os.Stdout)

	for _, r := range results {
		if !r.Valid {
			os.Exit(1)
		}
	}
}
