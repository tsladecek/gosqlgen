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
	supportedDrivers := []string{"gosqldriver_mysql"}

	debug := flag.Bool("debug", false, "debug")
	driver := flag.String("driver", "", "Driver to use. Supported: "+strings.Join(supportedDrivers, ", "))

	output := flag.String("out", "generatedMethods.go", "Path to output")
	outputTest := flag.String("outTest", "generatedMethods_test.go", "Path to output of test code")

	flag.Parse()

	if !slices.Contains(supportedDrivers, *driver) {
		fmt.Println("Error: Unsupported driver. Supported: " + strings.Join(supportedDrivers, ", "))
		os.Exit(1)
	}

	filename := os.Getenv("GOFILE")
	if filename == "" {
		fmt.Println("Error: GOFILE environment variable not set.")
		os.Exit(1)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	dbModel, err := gosqlgen.NewDBModel(f)
	if err != nil {
		panic(err)
	}

	if *debug {
		dbModel.Debug()
	}

	d, err := gosqldrivermysql.New()
	if err != nil {
		panic(err)
	}
	err = gosqlgen.CreateTemplates(d, dbModel, *output, *outputTest)
	if err != nil {
		panic(err)
	}
}
