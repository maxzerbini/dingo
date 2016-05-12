package {{.PackageName}}
{{if .HasImports}}
{{range .ImportPackages}}import "{{.}}"
{{end}}{{end}}
{{range .BizTypes}}
// {{.TypeName}} is a business object for {{.Model.TypeName}} entities.
type {{.TypeName}} struct {
	{{range .Fields}}{{.FieldName}} {{.FieldType}}{{end}}
}
// Create a new {{.TypeName}}
func New{{.TypeName}}() *{{.TypeName}} {
	return &{{.TypeName}}{ Dao:&{{.Dao.PackageName}}.{{.Dao.TypeName}}{} }
}
// Convert a model entity in a view-model
func (b *{{.TypeName}}) ToViewModel(m *{{.Model.PackageName}}.{{.Model.TypeName}}) *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{
	v := &{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{}
	{{range .Model.Fields}}{{if .IsNullable}}if m.{{.FieldName}}.Valid { v.{{.FieldName}} = m.{{.FieldName}}.{{.NullableFieldType}} }{{else}}v.{{.FieldName}} = m.{{.FieldName}}{{end}}
	{{end}}
	return v
}
// Convert a view-model in a model entity
func (b *{{.TypeName}}) ToModel(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) *{{.Model.PackageName}}.{{.Model.TypeName}}{
	m := &{{.Model.PackageName}}.{{.Model.TypeName}}{}
	{{range .Model.Fields}}{{if .IsNullable}}m.{{.FieldName}} =  {{.FieldType}} {Valid:true, {{.NullableFieldType}}:v.{{.FieldName}} }{{else}}m.{{.FieldName}} = v.{{.FieldName}}{{end}}
	{{end}}
	return m
}
{{if .IsReadOnly}}{{else}}// Insert a new {{.Model.TypeName}} entity and returns the last insert Id.
func (b *{{.TypeName}}) Insert(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (lastInsertId int64, err error) {
	return b.Dao.Insert(dao.Connection, b.ToModel(v))
}
// Update a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) Update(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (rowsAffected int64, err error) {
	return b.Dao.Update(dao.Connection, b.ToModel(v))
}
// Delete a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) Delete(v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (rowsAffected int64, err error) {
	return b.Dao.Delete(dao.Connection, b.ToModel(v))
}
// Find the {{.Model.TypeName}} entity by primary keys, returns nil if not found.
func (b *{{.TypeName}}) Find(vm *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error){
	m := b.ToModel(vm)
	m, err = b.Dao.FindByPrimaryKey(dao.Connection, {{range $i, $e := .Model.PKFields}}{{if $i}}, {{end}}m.{{.FieldName}}{{end}})
	if err != nil {
		return nil, err
	} else {
		return b.ToViewModel(m), nil
	}
}
{{if .ViewModel.IsSimplePK}}// Find the {{.Model.TypeName}} entity by primary key (converting it to the correct type), returns nil if not found.
func (b *{{.TypeName}}) FindById(id string) (v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error){
	vm := &{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{}
	m, err := b.Dao.FindByPrimaryKey(dao.Connection, vm.ConvertPK(id))
	if err != nil {
		return nil, err
	} else {
		return b.ToViewModel(m), nil
	}
}
// Update a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) UpdateById(id string, v *{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}) (rowsAffected int64, err error) {
	{{range .ViewModel.PKFields}}v.{{.FieldName}} = v.ConvertPK(id){{end}}
	return b.Dao.Update(dao.Connection, b.ToModel(v))
}
// Delete a {{.Model.TypeName}} entity and returns the number of affected rows.
func (b *{{.TypeName}}) DeleteById(id string) (rowsAffected int64, err error) {
	v := &{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}{}
	{{range .ViewModel.PKFields}}v.{{.FieldName}} = v.ConvertPK(id){{end}}
	return b.Dao.Delete(dao.Connection, b.ToModel(v))
}{{end}}{{end}}
// List the {{.Model.TypeName}} entities.
func (b *{{.TypeName}}) List(take int, skip int) (list []*{{.ViewModel.PackageName}}.{{.ViewModel.TypeName}}, err error) {
	mlist, err := b.Dao.List(dao.Connection, take, skip)
	if err != nil {
		return nil, err
	} else {
		for _, m := range mlist {
			list = append(list, b.ToViewModel(m))
		}
		return list, nil
	}
}
// Count the {{.Model.TypeName}} entities.
func (b *{{.TypeName}}) Count() (count int64, err error){
	return b.Dao.Count(dao.Connection)
}
{{end}}