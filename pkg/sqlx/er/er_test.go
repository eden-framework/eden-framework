package er_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/profzone/eden-framework/pkg/sqlx/er"
	"github.com/profzone/eden-framework/pkg/sqlx/generator/__examples__/database"
	"github.com/profzone/eden-framework/pkg/sqlx/postgresqlconnector"
)

func TestDatabaseERFromDB(t *testing.T) {
	er := er.DatabaseERFromDB(database.DBTest, &postgresqlconnector.PostgreSQLConnector{})
	data, _ := json.MarshalIndent(er, "", "  ")

	fmt.Println(string(data))
}
