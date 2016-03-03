package producers

import (
	"log"

	"github.com/maxzerbini/dingo/model"
)

func ProduceViewModelPackage(config *model.Configuration, schema *model.DatabaseSchema) (pkg *model.ViewModelPackage) {
	pkg = &model.ViewModelPackage{PackageName: "viewmodel", BasePackage: config.BasePackage}
	for _, table := range schema.Tables {
		mt := &model.ViewModelType{TypeName: getModelTypeName(table.TableName), PackageName: "viewmodel"}
		pkg.ViewModelTypes = append(pkg.ViewModelTypes, mt)
		for _, column := range table.Columns {
			field := &model.ViewModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getViewModelFieldType(pkg, column)}
			if column.IsNullable {
				if field.FieldType != "time.Time" { // exclude time fields
					field.IsNullable = true
				}
			}
			mt.Fields = append(mt.Fields, field)
		}
	}
	for _, view := range schema.Views {
		mt := &model.ViewModelType{TypeName: getModelTypeName(view.ViewName), PackageName: "model"}
		pkg.ViewModelTypes = append(pkg.ViewModelTypes, mt)
		for _, column := range view.Columns {
			field := &model.ViewModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getViewModelFieldType(pkg, column)}
			if column.IsNullable {
				if field.FieldType != "time.Time" { // exclude time fields
					field.IsNullable = true
				}
			}
			mt.Fields = append(mt.Fields, field)
		}
	}
	return pkg
}

func getViewModelFieldType(pkg *model.ViewModelPackage, column *model.Column) string {
	var ft string = ""
	switch column.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		ft = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		ft = "[]byte"
	case "date", "time", "datetime", "timestamp":
		ft = "time.Time"
		pkg.AppendImport("time")
	case "tinyint", "smallint":
		ft = "int32"
	case "int", "mediumint", "bigint":
		ft = "int64"
	case "float", "decimal", "double":
		ft = "float64"
	case "bit":
		ft = "bool"
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		ft = "string"
	}
	return ft
}
