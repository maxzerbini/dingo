package producers

import (
	"strings"

	"github.com/maxzerbini/dingo/model"
	"github.com/maxzerbini/dingo/util"
)

func ProduceServicePackage(config *model.Configuration, viewpkg *model.ViewModelPackage, bizpkg *model.BizPackage) (pkg *model.ServicePackage) {
	pkg = &model.ServicePackage{PackageName: "service", BasePackage: config.BasePackage}
	pkg.AppendImport("github.com/gin-gonic/gin")
	pkg.AppendImport(viewpkg.BasePackage + "/" + viewpkg.PackageName)
	pkg.AppendImport(bizpkg.BasePackage + "/" + bizpkg.PackageName)
	for _, biz := range bizpkg.BizTypes {
		srv := &model.ServiceType{TypeName: biz.ViewModel.TypeName + "Service", PackageName: "service"}
		srv.Biz = biz
		srv.ViewModel = biz.ViewModel
		srv.IsSimplePK = biz.ViewModel.IsSimplePK
		if config.ForcePluralResourceName {
			srv.ResourceName = util.Pluralize(strings.ToLower(biz.ViewModel.TypeName))
		} else {
			srv.ResourceName = strings.ToLower(biz.ViewModel.TypeName)
		}
		srv.Fields = append(srv.Fields, &model.BaseField{FieldName: "Biz", FieldType: "*" + biz.PackageName + "." + biz.TypeName})
		pkg.ServiceTypes = append(pkg.ServiceTypes, srv)
	}
	return pkg
}
