package g

import (
        "log"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() {
        var err error
        DB, err = sql.Open("mysql", Config().Database)
        if err != nil {
                log.Fatalln("open db fail:", err)
        }

        DB.SetMaxIdleConns(Config().MaxIdle)

        err = DB.Ping()
        if err != nil {
                log.Fatalln("ping db fail:", err)
        }
}