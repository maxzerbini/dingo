package producers

import (
	"github.com/maxzerbini/dingo/model"
)

func ProduceDaoPackage(config *model.Configuration, schema *model.DatabaseSchema, mpkg *model.ModelPackage) (pkg *model.DaoPackage) {
	pkg = &model.DaoPackage{PackageName: "dao", BasePackage: config.BasePackage}
	pkg.AppendImport(mpkg.BasePackage + "/" + mpkg.PackageName)
	pkg.AppendImport("database/sql")
	i := 0
	for _, table := range schema.Tables {
		dao := &model.DaoType{TypeName: mpkg.ModelTypes[i].TypeName + "Dao", PackageName: "dao"}
		dao.Model = mpkg.ModelTypes[i]
		dao.Entity = table
		pkg.DaoTypes = append(pkg.DaoTypes, dao)
		i++
	}
	i = 0
	for _, view := range schema.Views {
		dao := &model.DaoType{TypeName: mpkg.ViewModelTypes[i].TypeName + "Dao", PackageName: "dao"}
		dao.Model = mpkg.ModelTypes[i]
		dao.View = view
		pkg.ViewDaoTypes = append(pkg.ViewDaoTypes, dao)
		i++
	}
	return pkg
}
