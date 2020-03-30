package sqlx

import (
	"fmt"
	"os"
	"reflect"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"

	"github.com/sirupsen/logrus"
)

func NewFeatureDatabase(name string) *Database {
	if projectFeature, exists := os.LookupEnv("PROJECT_FEATURE"); exists && projectFeature != "" {
		name = name + "__" + projectFeature
	}
	return NewDatabase(name)
}

func NewDatabase(name string) *Database {
	return &Database{
		Database: builder.DB(name),
	}
}

type Database struct {
	*builder.Database
}

func (database *Database) Register(model Model) *builder.Table {
	database.mustStructType(model)
	rv := reflect.Indirect(reflect.ValueOf(model))
	table := builder.T(database.Database, model.TableName())
	ScanDefToTable(rv, table)
	database.Database.Register(table)
	return table
}

func (database Database) T(model Model) *builder.Table {
	database.mustStructType(model)
	return database.Database.Table(model.TableName())
}

func (database Database) mustStructType(model Model) {
	tpe := reflect.TypeOf(model)
	if tpe.Kind() != reflect.Ptr {
		panic(fmt.Errorf("model %s must be a pointer", tpe.Name()))
	}
	tpe = tpe.Elem()
	if tpe.Kind() != reflect.Struct {
		panic(fmt.Errorf("model %s must be a struct", tpe.Name()))
	}
}

func (database *Database) Insert(model Model) *builder.StmtInsert {
	table := database.T(model)

	fieldValues := FieldValuesFromStructByNonZero(model)

	if autoIncrementCol := table.AutoIncrement(); autoIncrementCol != nil {
		delete(fieldValues, autoIncrementCol.FieldName)
	}

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	return table.Insert().Columns(cols).Values(vals...)
}

func (database *Database) Update(model Model, zeroFields ...string) *builder.StmtUpdate {
	table := database.T(model)

	fieldValues := FieldValuesFromStructByNonZero(model, zeroFields...)

	if autoIncrementCol := table.AutoIncrement(); autoIncrementCol != nil {
		delete(fieldValues, autoIncrementCol.FieldName)
	}

	return table.Update().Set(table.AssignsByFieldValues(fieldValues)...)
}

func (database *Database) MustMigrateTo(db *DB, dryRun bool) {
	if err := database.MigrateTo(db, dryRun); err != nil {
		logrus.Panic(err)
	}
}

func (database *Database) MigrateTo(db *DB, dryRun bool) error {
	database.Register(&SqlMetaEnum{})

	currentDatabase := DBFromInformationSchema(db, database.Name, database.Tables.TableNames()...)

	if !dryRun {
		logrus.Debugf("=================== migrating database `%s` ====================", database.Name)
		defer logrus.Debugf("=================== migrated database `%s` ====================", database.Name)

		if currentDatabase == nil {
			currentDatabase = &Database{
				Database: builder.DB(database.Name),
			}
			if err := db.Do(currentDatabase.Create(true)).Err(); err != nil {
				return err
			}
		}

		for name, table := range database.Tables {
			currentTable := currentDatabase.Table(name)
			if currentTable == nil {
				if err := db.Do(table.Create(true)).Err(); err != nil {
					return err
				}
				continue
			}

			stmt := currentTable.Diff(table)
			if stmt != nil {
				if err := db.Do(stmt).Err(); err != nil {
					return err
				}
				continue
			}
		}

		if err := database.SyncEnum(db); err != nil {
			return err
		}

		return nil
	}

	if currentDatabase == nil {
		currentDatabase = &Database{
			Database: builder.DB(database.Name),
		}

		fmt.Printf("=================== need to migrate database `%s` ====================\n", database.Name)
		fmt.Println(currentDatabase.Create(true).Query)
		fmt.Printf("=================== need to migrate database `%s` ====================\n", database.Name)
	}

	for name, table := range database.Tables {
		currentTable := currentDatabase.Table(name)
		if currentTable == nil {
			fmt.Println(table.Create(true).Query)
			continue
		}

		stmt := currentTable.Diff(table)
		if stmt != nil {
			fmt.Println(stmt.Query)
			continue
		}
	}

	if err := database.SyncEnum(db); err != nil {
		return err
	}

	return nil
}
