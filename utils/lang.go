package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

// GetNowTime 获取当前时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Print 输出常规绿色字符
func Print(msg string) {
	fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, GetNowTime()+"[INFO]: "+msg, 0x1B)
}

// PrintErr 输出红色错误字符
func PrintErr(msg string) {
	fmt.Printf("%c[1;0;31m%s%c[0m\n", 0x1B, msg, 0x1B)
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("CheckError, error:%s\n", err.Error())
		panic(err.Error())
	}
}

// CheckDbConnect 校验数据库是否可以成功连接
func CheckDbConnect(dbType int, host string, port int, username string, password string, dbName string) (*sql.DB, bool) {
	var db *sql.DB
	if dbType == 1 {
		db, _ = sql.Open("mysql", username+":"+password+"@tcp("+host+":"+strconv.Itoa(port)+")/"+dbName)
		err := db.Ping()
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
	}
	return db, true
}
