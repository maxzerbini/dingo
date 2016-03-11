package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .ViewModelTypes}}
// View-Model object for {{.TypeName}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}
{{if .IsSimplePK}}
func (vm *{{.TypeName}}) ConvertPK(value string) {{.PKType}} {
	{{.PKStringConv}}
	return ret
}
{{end}}
{{end}}