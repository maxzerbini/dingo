package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
import "errors"

var Connection *sql.DB

{{range .DaoTypes}} {{$oc := len .Entity.OtherColumns}}
// Data access object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}
{{if .HasAutoIncrementPK}}
// Insert a new {{.Model.TypeName}} entity and returns the last insert Id.
func (dao *{{.TypeName}}) Insert(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(lastInsertId int64, err error) {
	q := "INSERT INTO \"{{.Entity.TableName}}\"({{range $i, $e := .Entity.OtherColumns}}{{if $i}}, {{end}}\"{{.ColumnName}}\"{{end}})"
	q += " VALUES ({{range $i, $e := .Entity.OtherColumns}}{{if $i}}, {{end}}{{.GetPostgresParam $i}}{{end}})"
	q += " RETURNING {{range $i, $e := .Entity.PrimaryKeys}}\"{{.ColumnName}}\"{{end}}"
	err = conn.QueryRow(q, {{range $i, $e := .Model.OtherFields}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}}).Scan(&lastInsertId)
    if err != nil {
		return -1, err
	}
	return lastInsertId, err
}{{else}}
// Insert a new {{.Model.TypeName}} entity and returns the last insert Id.
func (dao *{{.TypeName}}) Insert(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(lastInsertId int64, err error) {
	q := "INSERT INTO \"{{.Entity.TableName}}\" VALUES ({{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}{{.GetPostgresParam $i}}{{end}})"
	_, err = conn.Exec(q, {{range $i, $e := .Model.Fields}}{{if $i}}, {{end}}{{if .IsAutoInc}}sql.NullInt64{}{{else}}dto.{{.FieldName}}{{end}}{{end}})
    if err != nil {
		return -1, err
	}
	return 0, err
}{{end}}
// Update a {{.Model.TypeName}} entity and returns the number of affected rows.
func (dao *{{.TypeName}}) Update(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(rowsAffected int64, err error) {
	q := "UPDATE \"{{.Entity.TableName}}\" SET {{range $i, $e := .Entity.OtherColumns}}{{if $i}}, {{end}}\"{{.ColumnName}}\"={{.GetPostgresParam $i}}{{end}}"
	q += " WHERE {{range $i, $e := .Entity.PrimaryKeys}}{{if $i}} AND {{end}}\"{{.ColumnName}}\" = {{.GetPostgresParamFrom $i $oc}}{{end}}"
	res, err := conn.Exec(q, {{range $i, $e := .Model.OtherFields}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}}, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}})
    if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	return rowsAffected, err
}
// Delete a {{.Model.TypeName}} entity and returns the number of affected rows.
func (dao *{{.TypeName}}) Delete(conn *sql.DB, dto *{{.Model.PackageName}}.{{.Model.TypeName}})(rowsAffected int64, err error) {
	q := "DELETE FROM \"{{.Entity.TableName}}\""
	q += " WHERE {{range $i, $e := .Entity.PrimaryKeys}}{{if $i}} AND {{end}}\"{{.ColumnName}}\" = {{.GetPostgresParam $i}}{{end}}"
	res, err := conn.Exec(q, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}dto.{{.FieldName}}{{end}})
    if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	return rowsAffected, err
}
// Find the {{.Model.TypeName}} entity by primary keys, returns nil if not found.
func (dao *{{.TypeName}}) FindByPrimaryKey(conn *sql.DB, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}{{.FieldName}} {{.FieldType}}{{end}}) (dto *{{.Model.PackageName}}.{{.Model.TypeName}}, err error) {
	q := "SELECT {{range $i, $e := .Entity.Columns}}{{if $i}}, {{end}}\"{{.ColumnName}}\"{{end}} FROM \"{{.Entity.TableName}}\" WHERE {{range $i, $e := .Entity.PrimaryKeys}}{{if $i}} AND {{end}}\"{{.ColumnName}}\" = {{.GetPostgresParam $i}}{{end}}"
	rows, err := conn.Query(q, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}{{.FieldName}}{{end}})
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
		return nil, errors.New("Not found.")
	}
}
// List the {{.Model.TypeName}} entities.
func (dao *{{.TypeName}}) List(conn *sql.DB, take int, skip int) (list []*{{.Model.PackageName}}.{{.Model.TypeName}}, err error) {
	q := "SELECT {{range $i, $e := .Entity.Columns}}{{if $i}}, {{end}}\"{{.ColumnName}}\"{{end}} FROM \"{{.Entity.TableName}}\" LIMIT $1 OFFSET $2"
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
// Count the {{.Model.TypeName}} entities.
func (dao *{{.TypeName}}) Count(conn *sql.DB) (count int64, err error){
	q := "SELECT count(*) FROM \"{{.Entity.TableName}}\""
	rows, err := conn.Query(q)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		return 0, nil
	}
}
{{end}}
{{range .ViewDaoTypes}}
// Data access object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}
// List the {{.Model.TypeName}} entities in the view.
func (dao *{{.TypeName}}) List(conn *sql.DB, take int, skip int) (list []*{{.Model.PackageName}}.{{.Model.TypeName}}, err error){
	q := "SELECT {{range $i, $e := .View.Columns}}{{if $i}}, {{end}}{{.ColumnName}}{{end}} FROM {{.View.ViewName}} LIMIT $1 OFFSET $2"
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
// Count the {{.Model.TypeName}} entities in the view.
func (dao *{{.TypeName}}) Count(conn *sql.DB) (count int64, err error){
	q := "SELECT count(*) FROM {{.View.ViewName}}"
	rows, err := conn.Query(q)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		return 0, nil
	}
}
{{end}}