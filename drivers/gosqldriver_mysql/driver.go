package gosqldrivermysql

import (
	"fmt"

	"github.com/tsladecek/gosqlgen"
)

type driver struct{}

func NewDriver() driver {
	return driver{}
}

func (d driver) Get(table *gosqlgen.Table, keys []*gosqlgen.Column) (string, error) {
	return fmt.Sprintf("Get %s where %s = ?", table.Name, keys[0].Name), nil
}
func (d driver) Create(table *gosqlgen.Table) (string, error) {
	return "", nil
}
func (d driver) Update(table *gosqlgen.Table, keys []*gosqlgen.Column) (string, error) {
	return "", nil
}
func (d driver) Delete(table *gosqlgen.Table, keys []*gosqlgen.Column) (string, error) {
	return "", nil
}
