package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .DaoTypes}}
// Data access object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}
// Insert a new {{.Model.TypeName}} entity and returns the last insert Id.
func (dao *{{.TypeName}}) Insert(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(lastInsertId int64, err error){
	q := "INSERT INTO {{.Entity.TableName}} VALUES ({{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}?{{end}})"
	res, err := conn.Exec(q, {{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}{{if .IsAutoInc}}sql.NullInt64{}{{else}}dto.{{.FieldName}}{{end}}{{end}})
    if err != nil {
		return -1, err
	}
	lastInsertId, err = res.LastInsertId()
	return lastInsertId, err
}
// Update a {{.Model.TypeName}} entity and returns the number of affected rows.
func (dao *{{.TypeName}}) Update(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(rowsAffected int64, err error){
	q := "UPDATE {{.Entity.TableName}} SET {{range $i, $e := .Entity.Columns}}{{if $i}}, {{end}}{{.ColumnName}}=?{{end}}"
	q += " WHERE {{range $i, $e := .Entity.Columns}}{{if .IsPrimaryKey}}{{if $i}}, {{end}}{{.ColumnName}} = ?{{end}}{{end}}"
	res, err := conn.Exec(q, {{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}},{{range $i, $e := .Model.Fields}}{{if .IsPK}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}}{{end}})
    if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	return rowsAffected, err
}
// Find the {{.Model.TypeName}} entity by primary keys, returns nil if not found.
func (dao *{{.TypeName}}) FindByPrimaryKey(conn *sql.DB, {{range $i, $e := .Model.Fields}}{{if .IsPK}}{{if $i}}, {{end}}{{.FieldName}} {{.FieldType}}{{end}}{{end}}) (dto *{{.Model.PackageName}}.{{.Model.TypeName}}, err error){
	q := "SELECT {{range $i, $e := .Entity.Columns}}{{if $i}}, {{end}}{{.ColumnName}}{{end}} FROM {{.Entity.TableName}} WHERE {{range $i, $e := .Entity.Columns}}{{if .IsPrimaryKey}}{{if $i}}, {{end}}{{.ColumnName}} = ?{{end}}{{end}}"
	rows, err := conn.Query(q, {{range $i, $e := .Model.Fields}}{{if .IsPK}}{{if $i}}, {{end}}{{.FieldName}}{{end}}{{end}})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		dto = &{{.Model.PackageName}}.{{.Model.TypeName}}{}
		err := rows.Scan({{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}&dto.{{.FieldName}}{{end}})
		if err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, nil
	}
	
}
// List the {{.Model.TypeName}} entities.
func (dao *{{.TypeName}}) List(conn *sql.DB, take int32, skip int32) (list []*{{.Model.PackageName}}.{{.Model.TypeName}}, err error){
	q := "SELECT {{range $i, $e := .Entity.Columns}}{{if $i}}, {{end}}{{.ColumnName}}{{end}} FROM {{.Entity.TableName}} LIMIT ? OFFSET ?"
	rows, err := conn.Query(q, take, skip)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		dto := &{{.Model.PackageName}}.{{.Model.TypeName}}{}
		err := rows.Scan({{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}&dto.{{.FieldName}}{{end}})
		if err != nil {
			return nil, err
		}
		list = append(list, dto)
	}
	return list, nil
}
{{end}}