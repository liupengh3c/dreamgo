package data

import (
	"database/sql"
	"female/lib/tools"
	"female/log"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var db = &sql.DB{}

type dbConfig struct {
	IP        string `toml:"ip"`
	Port      string `toml:"port"`
	User      string `toml:"user"`
	Password  string `toml:"password"`
	DefaultDb string `toml:"default_db"`
}

func init() {
	dbConfigInfo := new(dbConfig)
	toml.DecodeFile(tools.GetCurrentDirectory()+"/conf/db.toml", dbConfigInfo)
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbConfigInfo.User, dbConfigInfo.Password, dbConfigInfo.IP, dbConfigInfo.Port, dbConfigInfo.DefaultDb)
	db, _ = sql.Open("mysql", str)

	return
}
func SelectRecord(table string, filed []string, conds map[string]string) ([]map[string]string, error) {
	strConds := ""
	fileds := strings.ToLower(strings.Join(filed, ","))
	data := []map[string]string{}
	for key, value := range conds {
		strConds += key + "\"" + value + "\"" + " and "
	}
	strConds = strings.ToLower(strConds)
	// strings.TrimRight(strConds, "and")
	//去除最后的and字符串
	strConds = strings.TrimSuffix(strConds, " and ")
	strQuery := "select " + fileds + " from " + table + " where " + strConds
	log.Debug(strQuery)
	rows, err := db.Query(strQuery)
	if err != nil {
		log.Fatal("db query fail, db.Query return err, errnum=" + err.Error())
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("rows Columns fail, rows.Columns return err, errnum=" + err.Error())
		return nil, err
	}
	//vals用来存放从数据库中读出的数据
	// vals := make([][]byte, len(cols)) //也可以这样定义
	vals := make([]string, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}
	// fmt.Println(scans)
	// var results []map[string]string
	for rows.Next() {
		rows.Scan(scans...) //scan函数的参数是一个空接口的切片
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		data = append(data, row)
	}
	// fmt.Println(queryData)
	rows.Close()
	return data, nil
}
func Query(query string) ([]map[string]string, error) {
	queryData := []map[string]string{}
	query = strings.ToLower(query)
	rows, err := db.Query(query)
	log.Debug(query)
	if err != nil {
		log.Fatal("db query fail, db.Query return err, errnum=" + err.Error() + " sql=" + query)
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("rows Columns fail, rows.Columns return err, errnum=" + err.Error())
		return nil, err
	}
	//vals用来存放从数据库中读出的数据
	// vals := make([][]byte, len(cols)) //也可以这样定义
	vals := make([]string, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}
	for rows.Next() {
		rows.Scan(scans...) //scan函数的参数是一个空接口的切片
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		queryData = append(queryData, row)
	}
	rows.Close()
	return queryData, nil
}
func InsertRecord(table string, data map[string]string) error {
	strKeys := "("
	strValue := "("
	length := len(data)
	vals := make([]string, 0, length)
	scans := make([]interface{}, length)
	for key, value := range data {
		strKeys += key + ","
		strValue += "?,"
		vals = append(vals, value)
	}
	for i := range vals {
		scans[i] = vals[i]
	}
	strKeys = strings.TrimSuffix(strKeys, ",") + ")"
	strValue = strings.TrimSuffix(strValue, ",") + ")"
	strKeys = "insert into " + table + strKeys + " values" + strValue
	log.Debug(strKeys)
	// fmt.Println(strKeys)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("db begin(InsertRecord) err, err=" + err.Error() + "sql=" + strKeys)
		return err
	}
	_, err = tx.Exec(strKeys, scans...)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	tx.Commit()
	return nil
}
func UpdateRecord(table string, data map[string]string, conds map[string]string) error {
	strKeys := ""
	length := len(data)
	// tx.Exec("UPdate user set age=?,id=? where uid=? and name="?",i,i)
	sliceVals := make([]string, 0, length)
	sliceData := make([]interface{}, length+len(conds))
	for key, value := range data {
		strKeys += key + "=?,"
		sliceVals = append(sliceVals, value)
	}
	strKeys = strings.TrimSuffix(strKeys, ",")
	strKeys += " where "
	for key, value := range conds {
		strKeys += key + "? and "
		sliceVals = append(sliceVals, value)
	}
	strKeys = strings.TrimSuffix(strKeys, " and ")
	for i := range sliceVals {
		sliceData[i] = sliceVals[i]
	}
	strKeys = "UPdate " + table + " set " + strKeys
	log.Debug(strKeys)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("db begin(UpdateRecord) err, err=" + err.Error())
		return err
	}
	tx.Exec(strKeys, sliceData...)
	tx.Commit()
	return nil
}
func DeleteRecord(table string, conds map[string]string) error {
	strKeys := ""
	sliceVals := make([]string, 0, len(conds))
	sliceData := make([]interface{}, len(conds))
	strKeys += " where "
	for key, value := range conds {
		strKeys += key + "? and "
		sliceVals = append(sliceVals, value)
	}
	strKeys = strings.TrimSuffix(strKeys, " and ")
	for i := range sliceVals {
		sliceData[i] = sliceVals[i]
	}
	strKeys = "Delete from " + table + strKeys
	log.Debug(strKeys)
	// fmt.Printf(strKeys)
	// fmt.Println(sliceData)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("db begin(DeleteRecord) err, err=" + err.Error())
		return err
	}
	tx.Exec(strKeys, sliceData...)
	tx.Commit()
	return nil
}
