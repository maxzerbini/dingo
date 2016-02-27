package model

type BaseField struct {
	FieldName string
	FieldType string
}

type ModelField struct {
	FieldName     string
	FieldType     string
	FieldMetadata string
	IsPK          bool
	IsAutoInc     bool
}

type ModelType struct {
	PackageName string
	TypeName    string
	Fields      []*ModelField
	PKFields    []*ModelField
	OtherFields []*ModelField
}

type DaoType struct {
	PackageName string
	TypeName    string
	Fields      []*BaseField
	Model       *ModelType
	Entity      *Table
	View        *View
}

type ModelPackage struct {
	BasePackage    string
	PackageName    string
	ImportPackages []string
	ModelTypes     []*ModelType
	ViewModelTypes []*ModelType
}

type DaoPackage struct {
	BasePackage    string
	PackageName    string
	ImportPackages []string
	DaoTypes       []*DaoType
	ViewDaoTypes   []*DaoType
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

func (pkg *DaoPackage) HasImports() bool {
	return len(pkg.ImportPackages) > 0
}

func (pkg *DaoPackage) AppendImport(pkgName string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == pkgName {
			return false
		}
	}
	pkg.ImportPackages = append(pkg.ImportPackages, pkgName)
	return true
}
