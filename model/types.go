package model

type ModelField struct {
	FieldName     string
	FieldType     string
	FieldMetadata string
}

type ModelType struct {
	TypeName string
	Fields   []*ModelField
}

type ModelPackage struct {
	PackageName    string
	ImportPackages []string
	ModelTypes     []*ModelType
}

func (pkg *ModelPackage) HasImports() bool {
	return len(pkg.ImportPackages) > 0
}

func (pkg *ModelPackage) AppendImport(pkgName string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == pkgName {
			return false
		}
	}
	pkg.ImportPackages = append(pkg.ImportPackages, pkgName)
	return true
}
