package connection

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sajad-dev/go-framwork/Config/setting"
	"github.com/sajad-dev/go-framwork/Exception/exception"
)

var Database *sql.DB

func Connection() {
	db, err := sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(%s:%s)/",setting.USERNAME,setting.PASSWORD,setting.IP,setting.PORT))
	exception.Log(err)

	databasename := setting.DATABASE
	_, _ = db.Exec("CREATE DATABASE " + databasename)

	db, err = sql.Open("mysql", fmt.Sprintf( "%s:%s@tcp(%s:%s)/%s",setting.USERNAME,setting.PASSWORD,setting.IP,setting.PORT,databasename))
	exception.Log(err)

	if err := db.Ping(); err != nil {
		exception.Log(err)
	}

	Database = db
}
