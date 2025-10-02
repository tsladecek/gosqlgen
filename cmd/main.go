//go:build ignore

package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"

	"github.com/tsladecek/gosqlgen"
	gosqldrivermysql "github.com/tsladecek/gosqlgen/drivers/gosqldriver_mysql"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(strings.TrimPrefix(err.Error(), "\n"))
		os.Exit(1)
	}
}

func run() error {
	supportedDrivers := []string{"gosqldriver_mysql"}

	debug := flag.Bool("debug", false, "debug")
	driver := flag.String("driver", "", "Driver to use. Supported: "+strings.Join(supportedDrivers, ", "))

	output := flag.String("out", "generatedMethods.go", "Path to output")
	outputTest := flag.String("outTest", "generatedMethods_test.go", "Path to output of test code")

	flag.Parse()

	if !slices.Contains(supportedDrivers, *driver) {
		return gosqlgen.Errorf("unsupported driver %s; supported are: %s", *driver, strings.Join(supportedDrivers, ", "))
	}

	filename := os.Getenv("GOFILE")
	if filename == "" {
		return gosqlgen.Errorf("GOFILE environment variable not set.")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return gosqlgen.Errorf("when parsing file: %w", err)
	}

	dbModel, err := gosqlgen.NewDBModel(fset, f)
	if err != nil {
		return gosqlgen.Errorf("when constructing db model: %w", err)
	}

	if *debug {
		dbModel.Debug()
	}

	d, err := gosqldrivermysql.New()
	if err != nil {
		return gosqlgen.Errorf("when initializing driver: %w", err)
	}
	err = gosqlgen.CreateTemplates(d, dbModel, *output, *outputTest)
	if err != nil {
		return gosqlgen.Errorf("when generating code from templates: %w", err)
	}

	return nil
}
