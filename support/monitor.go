package support

import (
	"container/list"
	"data-sync/model"
	"data-sync/utils"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"time"
)

var (
	configTableList = list.New()
	dbMap           = map[int]*sql.DB{}
	db, _           = sql.Open("sqlite3", "./Data.db")
)

// CheckStart 程序启动校验&数据库连接建立
func CheckStart() {
	utils.Print("数据同步程序启动...")

	rows, err := db.Query("select * from db_list")
	utils.CheckError(err)
	utils.Print("开始校验数据库连接合法性")

	for rows.Next() {
		var id int
		var name string
		var dbType int
		var ip string
		var port int
		var username string
		var password string
		var dbName string
		err = rows.Scan(&id, &name, &dbType, &ip, &port, &username, &password, &dbName)
		utils.CheckError(err)
		utils.Print("数据库【" + name + "】开始测试连接")
		conn, tag := utils.CheckDbConnect(dbType, ip, port, username, password, dbName)
		if !tag {
			panic("数据库" + name + "异常，无法正常连接")
		} else {
			dbMap[id] = conn
			utils.Print("数据库【" + name + "】正常连接")
		}
	}

	// 读取配置项准备进行同步数据
	rows, err = db.Query("select * from config_table")
	utils.CheckError(err)
	for rows.Next() {

		configTable := new(model.ConfigTable)
		err = rows.Scan(&configTable.Id, &configTable.Name, &configTable.SourceId, &configTable.TargetId, &configTable.SourceTable, &configTable.TargetTable, &configTable.Interval, &configTable.LastTime, &configTable.WhereSql)
		utils.CheckError(err)
		fieldList := list.New()
		fieldRow, err := db.Query("select * from config_table_field where config_id = " + strconv.Itoa(configTable.Id))
		utils.CheckError(err)
		for fieldRow.Next() {
			field := new(model.ConfigTableField)
			err = fieldRow.Scan(&field.Id, &field.ConfigId, &field.SourceField, &field.TargetField)
			fieldList.PushBack(field)
		}

		configTable.Field = fieldList
		// 获取字段同步数据
		configTableList.PushBack(configTable)
		utils.CheckError(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		utils.CheckError(err)
	}(rows)
}

// Monitor 实时监控需要处理的处理
func Monitor() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			unixTime := time.Now().Unix()
			timeStr := strconv.FormatInt(unixTime, 10)
			// 遍历配置开始
			for i := configTableList.Front(); i != nil; i = i.Next() {
				configTable := i.Value.(*model.ConfigTable)
				oldLastTime := configTable.LastTime
				// 判断是否到达时间
				if configTable.LastTime == "" {
					// 没有值
					configTable.LastTime = timeStr
				} else {
					lastTime, err := strconv.ParseInt(configTable.LastTime, 10, 64)
					utils.CheckError(err)
					if lastTime+int64(configTable.Interval) < unixTime {
						configTable.LastTime = timeStr
					}
				}
				if configTable.LastTime == timeStr {
					_, err := db.Exec("update config_table set last_time = '" + configTable.LastTime + "' where id = " + strconv.Itoa(configTable.Id))
					utils.CheckError(err)
					utils.Print("配置：" + configTable.Name + "，上次执行时间：" + oldLastTime + " 准备执行同步")
					// 连接源数据库
					sourceDb := dbMap[configTable.SourceId]

					sourceField := ""
					targetField := ""

					// 获取要转移的字段
					for fieldItem := configTable.Field.Front(); fieldItem != nil; fieldItem = fieldItem.Next() {
						tableField := fieldItem.Value.(*model.ConfigTableField)
						sourceField += tableField.SourceField + ","
						targetField += tableField.TargetField + ","
					}

					sourceField = strings.TrimRight(sourceField, ",")
					targetField = strings.TrimRight(targetField, ",")
					rows, err := sourceDb.Query("select " + sourceField + " from " + configTable.SourceTable + " " + configTable.WhereSql)
					utils.CheckError(err)
					for rows.Next() {
						cols, _ := rows.Columns()
						row := make([]interface{}, len(cols))
						rowData := make([]interface{}, len(row))

						for i := range row {
							rowData[i] = &row[i]
						}

						val := ""

						if err := rows.Scan(rowData...); err == nil {
							for i := range row {
								if v, ok := row[i].([]byte); ok {
									val += "'" + string(v) + "',"
								}
							}
						} else {
							utils.CheckError(err)
						}

						val = strings.TrimRight(val, ",")
						// 构建新的sql语句
						insertSql := "insert into " + configTable.TargetTable + " (" + targetField + ") values(" + val + ")"
						// 插入目标的数据库
						_, insErr := dbMap[configTable.TargetId].Exec(insertSql)
						utils.CheckError(insErr)
						utils.Print(configTable.Name + "->数据同步中..." + insertSql)
					}
					utils.Print("配置：" + configTable.Name + " 同步完成")
				}
			}
		}
	}()
}
