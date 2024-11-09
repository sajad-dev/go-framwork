package main_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/exp/slices"

	"github.com/sajad-dev/go-framwork/Config/setting"
	"github.com/sajad-dev/go-framwork/Database/migration"
	testutils "github.com/sajad-dev/go-framwork/Test-Utils"
)

func TestMigrationTables(t *testing.T) {
	var db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", setting.USERNAME, setting.PASSWORD, setting.IP, setting.PORT, setting.DATABASE))
	if err == nil {
		qu, err := db.Query("SHOW TABLES")
		if err != nil {
			return
		}
		x := 0
		table, _ := testutils.MiggarionListAppend()
		for qu.Next() {
			x++
			var name = ""
			qu.Scan(&name)
			if !slices.Contains(table, name) {
				t.Fatalf("Database %s not deleted", name)
			}
		}

		if len(table) != x {
			t.Fatal("You have problem in tables")
		}

	}

}

func TestMigrationTablesParams(t *testing.T) {
	var db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", setting.USERNAME, setting.PASSWORD, setting.IP, setting.PORT, setting.DATABASE))
	if err != nil {
		table, function := testutils.MiggarionListAppend()

		for i, v := range function {
			rqfunc := migration.GetFromFunc(v)
			if len(rqfunc) == 0 {
				t.Fatal("Your migration not be empty")
			}

			rqdb, err := db.Query(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table[i]))

			if err != nil {
				t.Fatal(err.Error())

			}
			arr := []string{}
			for rqdb.Next() {
				var name, tp, null, extra, key, privileges, comment string
				var collation, df sql.NullString

				if err := rqdb.Scan(&name, &tp, &collation, &null, &key, &df, &extra, &privileges, &comment); err != nil {
					t.Fatal("Error scanning row:", err)
				}
				if null == "YES" {
					null = "NULL"
				} else {
					null = "NOT NULL"
				}
				dfStr := "DEFAULT ''"
				if df.Valid {
					dfStr = fmt.Sprintf("DEFAULT '%s'", df.String)

				}
				if key == "PRI" {
					key = "PRIMARY"
				}
				if key == "UNI" {
					key = "UNIQUE"
				}
				if extra != "" {
					extra = fmt.Sprintf(" %s", extra)
				}

				str := fmt.Sprintf("%s %s %s %s%s %s", name, strings.ToUpper(tp), null, dfStr, strings.ToUpper(extra), key)
				if name == "id" {
					str = str + " KEY"
					str = strings.ReplaceAll(str, "INT(11)", "INT")
					str = strings.ReplaceAll(str, "DEFAULT '' ", "")
				}
				arr = append(arr, str)

			}
			for ind, val := range arr {
				if !strings.Contains(rqfunc[ind], val) {
					t.Fatal("Migration create not like migration exsist")
				}
			}

		}
	}

}
