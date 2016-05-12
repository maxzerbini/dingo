package producers

import (
	"github.com/maxzerbini/dingo/model"
)

func ProduceBizPackage(config *model.Configuration, mpkg *model.ModelPackage, daopkg *model.DaoPackage, viewpkg *model.ViewModelPackage) (pkg *model.BizPackage) {
	pkg = &model.BizPackage{PackageName: "biz", BasePackage: config.BasePackage}
	pkg.AppendImport(mpkg.BasePackage + "/" + mpkg.PackageName)
	pkg.AppendImport(daopkg.BasePackage + "/" + daopkg.PackageName)
	pkg.AppendImport(viewpkg.BasePackage + "/" + viewpkg.PackageName)
	if mpkg.HasImport("database/sql") {
		pkg.AppendImport("database/sql")
	}
	if mpkg.HasImport("github.com/go-sql-driver/mysql") {
		pkg.AppendImport("github.com/go-sql-driver/mysql")
	}
	if mpkg.HasImport("github.com/lib/pq") {
		pkg.AppendImport("github.com/lib/pq")
	}
	i := 0
	for _, table := range mpkg.ModelTypes {
		biz := &model.BizType{TypeName: mpkg.ModelTypes[i].TypeName + "Biz", PackageName: "biz"}
		biz.Model = table
		biz.Dao = daopkg.DaoTypes[i]
		biz.ViewModel = viewpkg.ViewModelTypes[i]
		biz.Fields = append(biz.Fields, &model.BaseField{FieldName: "Dao", FieldType: "*" + daopkg.DaoTypes[i].PackageName + "." + daopkg.DaoTypes[i].TypeName})
		pkg.BizTypes = append(pkg.BizTypes, biz)
		i++
	}
	j := 0
	for _, view := range mpkg.ViewModelTypes {
		biz := &model.BizType{TypeName: mpkg.ViewModelTypes[j].TypeName + "Biz", PackageName: "biz"}
		biz.Model = view
		biz.Dao = daopkg.ViewDaoTypes[j]
		biz.ViewModel = viewpkg.ViewModelTypes[i]
		biz.Fields = append(biz.Fields, &model.BaseField{FieldName: "Dao", FieldType: "*" + daopkg.ViewDaoTypes[j].PackageName + "." + daopkg.ViewDaoTypes[j].TypeName})
		biz.IsReadOnly = true
		pkg.BizTypes = append(pkg.BizTypes, biz)
		i++
		j++
	}
	return pkg
}
