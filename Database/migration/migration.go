package migration

import (
	"bytes"
	"database/sql"
	"fmt"
	"golang.org/x/exp/slices"
	"os/exec"
	"reflect"
	"regexp"

	// "regexp"
	"strings"

	// "database/sql"

	"github.com/sajad-dev/go-framwork/Config/setting"
	"github.com/sajad-dev/go-framwork/Database/connection"
	"github.com/sajad-dev/go-framwork/Exception/exception"
)

type Migrate struct {
	Table     string
	Funcation string
}

type Mig struct{}

var MigrateList []Migrate

var _MigrationInterface = (*Migrate)(nil)

func Handel() {
	MiggarionListAppend()

	HandelCheckTable()
}

func MiggarionListAppend() {
	ou := exec.Command("go", "doc", "migration", "MigrationInterface")

	output, err := ou.Output()
	exception.Log(err)

	slice := bytes.Split(output, []byte("\n"))
	for _, s := range slice {
		text := string(strings.ReplaceAll(string(s), "func", ""))
		index := strings.Index(text, "(")
		migration_index := strings.Index(text, "Migration")
		if index > 0 && migration_index > 0 {
			table_name := strings.ToLower(text[:index][:migration_index])
			mi := Migrate{Table: strings.TrimSpace(string(table_name)), Funcation: strings.TrimSpace(string(text[:index]))}
			MigrateList = append(MigrateList, mi)
		}
	}
}

func checkDeletedTable() {
	row, err := connection.Database.Query("SHOW TABLES")
	exception.Log(err)

	tableArr := []string{}
	for _, v := range MigrateList {
		tableArr = append(tableArr, v.Table)
	}
	for row.Next() {
		var table string
		row.Scan(&table)

		if !slices.Contains(tableArr, table) {
			connection.Database.Query(fmt.Sprintf("DROP TABLE %s", table))
		}

	}

}
func HandelCheckTable() {

	database := setting.DATABASE
	checkDeletedTable()


	sqlqu := fmt.Sprintf(` 
SELECT table_name
FROM information_schema.tables 
WHERE table_schema = '%s';
 `, database)

	row, err := connection.Database.Query(sqlqu)
	exception.Log(err)

	for row.Next() {
		var tb string

		err := row.Scan(&tb)
		exception.Log(err)

	}
	for _, m := range MigrateList {

		sql_qu := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", m.Table)
		qu, err := connection.Database.Query(sql_qu)
		exception.Log(err)

		x := 0

		for qu.Next() {

			var (
				field, fieldtype, null, key, extra, privileges, comment string
				collection, defult                                      sql.NullString
			)
			err = qu.Scan(&field, &fieldtype, &collection, &null, &key, &defult, &extra, &privileges, &comment)
			exception.Log(err)

			HandelUpdate(field, fieldtype, collection, null, defult, key, extra, privileges, comment, m.Funcation, m.Table, x)
			x++

		}
		strSlice := GetFromFunc(m.Funcation)
		if len(strSlice) > x {
			for i := x; i < len(strSlice); i++ {
				AddTable(strSlice[i], m.Table)
			}
		}

	}
}

func HandelUpdate(field string, fieldtype string, collection sql.NullString, null string, defult sql.NullString, key string, extera string, privileges string, comment string, function string, table string, x int) {
	strSlice := GetFromFunc(function)

	if x >= len(strSlice) {
		sql_qu := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, field)
		_, err := connection.Database.Query(sql_qu)
		exception.Log(err)
		return
	}

	if key == "UNI" && strings.Contains(strSlice[x], "UNIQUE") {

		sql_qu := fmt.Sprintf("SELECT DISTINCT COLUMN_NAME, INDEX_NAME   FROM information_schema.statistics  WHERE table_name = '%s'", table)
		qu, err := connection.Database.Query(sql_qu)
		exception.Log(err)

		unilist := []string{}
		for qu.Next() {
			var column, key string
			qu.Scan(&column, &key)
			if column == field {
				unilist = append(unilist, key)
			}
		}
		for _, v := range unilist {
			sql_qu := fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, v)
			connection.Database.Query(sql_qu)
			exception.Log(err)
		}

	}

	if key == "PRI" {
		strSlice[x] = strings.ReplaceAll(strSlice[x], "PRIMARY KEY", "")
	}
	UpdateTable(table, field, strSlice[x])

}

func UpdateTable(table string, column_name_old string, parametr string) {
	sql_qu := fmt.Sprintf("ALTER TABLE %s CHANGE %s %s ", table, column_name_old, parametr)
	if strings.Contains(sql_qu, "FOREIGN") {
		sql_du2 := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, table+"_"+column_name_old+"_fg")
		connection.Database.Query(sql_du2)
		sql_qu = fmt.Sprintf("ALTER TABLE %s CHANGE COLUMN %s %s ", table, column_name_old, parametr)
		sql_qu = strings.Replace(sql_qu, "CONSTRAINT", "ADD CONSTRAINT ", 1)
	}

	_, err := connection.Database.Query(sql_qu)
	exception.Log(err)
}

func AddTable(strSlice string, table string) {
	if strings.Contains(strSlice, "FOREIGN") {
		re := regexp.MustCompile(`\((.*?)\)`)
		matches := re.FindStringSubmatch(strSlice)
		sql_qu := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s INT", table, matches[1])
		_, err := connection.Database.Query(sql_qu)
		exception.Log(err)
		sql_qu = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", table, matches[1], strSlice)
		_, err = connection.Database.Query(sql_qu)
		exception.Log(err)
		return
	}
	sql_qu := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, strSlice)
	_, err := connection.Database.Query(sql_qu)
	exception.Log(err)
}

func CreateAll() {
	MiggarionListAppend()
	for _, v := range MigrateList {
		CreateTable(v.Funcation, v.Table)
	}
}

func CreateTable(function string, table string) {
	mig := Mig{}
	method := reflect.ValueOf(mig).MethodByName(function)
	args := []reflect.Value{}

	column := method.Call(args)

	strSlice := column[0].Interface().([]string)
	co := strings.Join(strSlice, " , ")

	str_sql := fmt.Sprintf("CREATE TABLE  IF NOT EXISTS %s (%s)", table, co)
	_, err := connection.Database.Query(str_sql)
	exception.Log(err)

}

func DropTable() {
	database := setting.DATABASE
	sql := fmt.Sprintf(` 
SELECT table_name
FROM information_schema.tables 
WHERE table_schema = '%s';
 `, database)
	row, err := connection.Database.Query(sql)
	exception.Log(err)

	for row.Next() {
		var tb string

		err := row.Scan(&tb)
		exception.Log(err)

		_, err = connection.Database.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", tb))
		exception.Log(err)

	}
}

func IntPrimary(name string, increment bool) string {
	getincrement := ""
	if increment {
		getincrement = "AUTO_INCREMENT"
	} else {
		getincrement = ""

	}
	return fmt.Sprintf("%s INT NOT NULL %s PRIMARY KEY", name, getincrement)
}

func StringVar(name string, max int) string {
	return fmt.Sprintf("%s VARCHAR(%d) ", name, max)
}



func GetFromFunc(function string) []string {
	mig := Mig{}
	method := reflect.ValueOf(mig).MethodByName(function)

	if !method.IsValid() {
		return []string{}
	}

	args := []reflect.Value{}
	column := method.Call(args)

	if len(column) == 0 {
		return []string{}
	}

	if result, ok := column[0].Interface().([]string); ok {
		return result
	} else {
		return []string{}
	}
}
