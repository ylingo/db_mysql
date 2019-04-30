package db_mysql

import (
	"fmt"
	"strings"
	"sync"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var once sync.Once
var database *sqlx.DB
var logTofile bool = false
var errTofile bool = false

func InitDb(driver string, dsn string, maxopenconns, maxidleconns int, logfile bool, errfile bool) {

	onceInit := func() {
		logTofile = logfile
		errTofile = errfile
		//fmt.Printf("%p, %T\n", database, database)
		database, _ = sqlx.Open(driver, dsn)

		//fmt.Printf("%p, %T\n", database, database)
		if err := database.Ping(); err != nil {
			fmt.Println(err)
			return
		}
		database.SetMaxOpenConns(maxopenconns)
		database.SetMaxIdleConns(maxidleconns)

		//fmt.Println("数据库打开了")
	}
	once.Do(onceInit)
}

func SqlStringCheck(str string) string {
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "'", "")
	return str
}
