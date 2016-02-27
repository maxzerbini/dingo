package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .ModelTypes}}
// Data transfer object for {{.TypeName}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}
{{end}}
{{range .ViewModelTypes}}
// Data transfer object for view {{.TypeName}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}} `{{.FieldMetadata}}`
	{{end}}
}

{{end}}