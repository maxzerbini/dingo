package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .ModelTypes}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}

{{end}}