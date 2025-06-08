package gosqlgen

import (
	"fmt"
	"strings"
)

type Driver interface {
	Get(table *Table, keys []*Column) (string, error)
	Create(table *Table) (string, error)
	Update(table *Table, keys []*Column) (string, error)
	Delete(table *Table, keys []*Column) (string, error)
}

func SaveTemplate(t string, path string) error {
	fmt.Println(t)
	return nil
}

func CreateTemplates(d Driver, model *DBModel) error {
	res := make([]string, 0, len(model.Tables)*4)

	for _, table := range model.Tables {
		pk, bk, err := table.PkAndBk()
		if err != nil {
			return fmt.Errorf("Failed to fetch primary and business keys: %w", err)
		}
		temp, err := d.Get(table, pk)
		if err != nil {
			return fmt.Errorf("Failed to create GET template by primary keys for table %s: %w", table.Name, err)
		}
		res = append(res, temp)

		if bk != nil {
			temp, err := d.Get(table, bk)
			if err != nil {
				return fmt.Errorf("Failed to create GET template by business keys for table %s: %w", table.Name, err)
			}
			res = append(res, temp)
		}
	}

	fmt.Println(strings.Join(res, "\n"))
	return nil
}
