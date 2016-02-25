package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .DaoTypes}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}
func (dao *{{.TypeName}}) Insert(dto {{.Model.PackageName}}.{{.Model.TypeName}}){
	
}
{{end}}