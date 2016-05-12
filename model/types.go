package model

import (
	"strconv"
)

type BaseField struct {
	FieldName string
	FieldType string
}

type ModelField struct {
	FieldName         string
	FieldType         string
	FieldMetadata     string
	IsPK              bool
	IsAutoInc         bool
	IsNullable        bool
	NullableFieldType string
}

type ViewModelField struct {
	FieldName  string
	FieldType  string
	IsNullable bool
}

type ModelType struct {
	PackageName string
	TypeName    string
	Fields      []*ModelField
	PKFields    []*ModelField
	OtherFields []*ModelField
}

type DaoType struct {
	PackageName        string
	TypeName           string
	Fields             []*BaseField
	Model              *ModelType
	Entity             *Table
	View               *View
	HasAutoIncrementPK bool
}

type ViewModelType struct {
	PackageName  string
	TypeName     string
	Fields       []*ViewModelField
	PKFields     []*ViewModelField
	IsSimplePK   bool
	PKType       string
	PKStringConv string
}

type BizType struct {
	PackageName string
	TypeName    string
	Fields      []*BaseField
	Model       *ModelType
	ViewModel   *ViewModelType
	Dao         *DaoType
	IsReadOnly  bool
}

type ServiceType struct {
	PackageName  string
	TypeName     string
	ResourceName string
	Fields       []*BaseField
	ViewModel    *ViewModelType
	Biz          *BizType
	IsSimplePK   bool
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

type ViewModelPackage struct {
	BasePackage    string
	PackageName    string
	ImportPackages []string
	ViewModelTypes []*ViewModelType
}

type BizPackage struct {
	BasePackage    string
	PackageName    string
	ImportPackages []string
	BizTypes       []*BizType
}

type ServicePackage struct {
	BasePackage    string
	PackageName    string
	ImportPackages []string
	ServiceTypes   []*ServiceType
}

func (pkg *ModelPackage) HasImport(impPkg string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == impPkg {
			return true
		}
	}
	return false
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

func (pkg *ViewModelPackage) HasImports() bool {
	return len(pkg.ImportPackages) > 0
}

func (pkg *ViewModelPackage) AppendImport(pkgName string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == pkgName {
			return false
		}
	}
	pkg.ImportPackages = append(pkg.ImportPackages, pkgName)
	return true
}

func (pkg *BizPackage) HasImports() bool {
	return len(pkg.ImportPackages) > 0
}

func (pkg *BizPackage) AppendImport(pkgName string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == pkgName {
			return false
		}
	}
	pkg.ImportPackages = append(pkg.ImportPackages, pkgName)
	return true
}

func (pkg *ServicePackage) HasImports() bool {
	return len(pkg.ImportPackages) > 0
}

func (pkg *ServicePackage) AppendImport(pkgName string) bool {
	for _, imp := range pkg.ImportPackages {
		if imp == pkgName {
			return false
		}
	}
	pkg.ImportPackages = append(pkg.ImportPackages, pkgName)
	return true
}

func (mf ModelField) GetPostgresParam(i int) string {
	return "$" + strconv.Itoa(i+1)
}

func (mf ModelField) GetPostgresParamFrom(i int, s int) string {
	return "$" + strconv.Itoa(s+i+1)
}
