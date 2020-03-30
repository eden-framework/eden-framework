package builder

import (
	"fmt"
)

func DB(name string) *Database {
	return &Database{
		Name:   name,
		Tables: Tables{},
	}
}

type Database struct {
	Name string
	Tables
}

func (d *Database) Register(table *Table) *Database {
	table.DB = d
	d.Tables.Add(table)
	return d
}

func (d *Database) Table(name string) *Table {
	if table, ok := d.Tables[name]; ok {
		return table
	}
	return nil
}

func (d *Database) String() string {
	return quote(d.Name)
}

func (d *Database) Drop() *Stmt {
	return (*Stmt)(Expr(fmt.Sprintf("DROP DATABASE %s", d.String())))
}

func (d *Database) Create(ifNotExists bool) *Stmt {
	if ifNotExists {
		return (*Stmt)(Expr(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", d.String())))
	}
	return (*Stmt)(Expr(fmt.Sprintf("CREATE DATABASE %s", d.String())))
}
