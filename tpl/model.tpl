package model

type {{.Name}} struct {
    {{range .Data}}
        {{.ColumnName}}  {{.Type}}  {{.Tag}} 
    {{end}}

}