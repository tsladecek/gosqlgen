//go:build ignore

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/tsladecek/gosqlgen"
	gosqldrivermysql "github.com/tsladecek/gosqlgen/drivers/gosqldriver_mysql"
)

func main() {
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

	for _, table := range dbModel.Tables {
		fmt.Printf("%+v\n", table)
		for _, c := range table.Columns {
			fmt.Printf("%+v\n", c)
		}
	}

	d := gosqldrivermysql.NewDriver()
	gosqlgen.CreateTemplates(d, dbModel)
}
