package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
import "strconv"
{{range .ServiceTypes}}
// Service object for {{.TypeName}}
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}
	{{end}}
}
// Endpoint GET [basehost]/{{.ResourceName}}?take=[value]&skip=[value]
func (s *{{.TypeName}}) List(c *gin.Context) {
    take,_ := strconv.Atoi(c.DefaultQuery("take", "10"))
    skip,_ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	var list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}
    list, err := s.Biz.List(take, skip)
	if err != nil {
		c.JSON(400, err)
	} else {
		c.JSON(200, list)
	}
}
{{end}}