package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "string", // time.Time
	"datetime":           "string", // time.Time
	"timestamp":          "string", // time.Time
	"time":               "string", // time.Time
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

type column struct {
	ColumnName    string
	Type          string
	Nullable      string
	TableName     string
	ColumnComment interface{}
	Tag           string
}

type ModelTpl struct {
	Name    string
	Extends string
	Data    []MData
}

type MData struct {
	ColumnName string
	Type       string
	Tag        string
}

func main() {
	tableColumns := make(map[string][]column)
	db, err := sql.Open("mysql", "hdller:Hdlltest888@tcp(192.168.3.202:3306)/buyer")
	if err != nil {
		fmt.Println("tttccc")
		fmt.Println(err)
		return
	}
	var sqlStr = `SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()`
	// 是否指定了具体的table

	sqlStr += fmt.Sprintf(" AND TABLE_NAME = '%s'", "buyer")

	// sql排序
	sqlStr += " order by TABLE_NAME asc, ORDINAL_POSITION asc"

	fmt.Println(sqlStr)

	rows, err := db.Query(sqlStr)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	var modelTpl ModelTpl

	for rows.Next() {
		mData := MData{}
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.Type, &col.Nullable, &col.TableName, &col.ColumnComment)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if modelTpl.Name == "" {
			modelTpl.Name = col.TableName
		}

		//col.Json = strings.ToLower(col.ColumnName)
		mData.Type = typeForMysqlToGo[col.Type]
		mData.Tag = "`json:\"" + col.ColumnName + "\"`"
		mData.ColumnName = camelCase(col.ColumnName)
		modelTpl.Data = append(modelTpl.Data, mData)

		if _, ok := tableColumns[col.TableName]; !ok {
			tableColumns[col.TableName] = []column{}
		}
		tableColumns[col.TableName] = append(tableColumns[col.TableName], col)
	}

	fmt.Println(modelTpl)
	bs, err := ioutil.ReadFile("tpl/model.tpl")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	t := template.New("tpl/model.tpl")
	tmpl, err := t.Parse(string(bs))

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	newbytes := bytes.NewBufferString("")

	err = tmpl.Execute(newbytes, modelTpl)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	tplcontent, err := ioutil.ReadAll(newbytes)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", tplcontent)

}

func camelCase(str string) string {
	// 是否有表前缀, 设置了就先去除表前缀
	var text string
	//for _, p := range strings.Split(name, "_") {
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}
	return text
}
