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
			field := &model.ViewModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getViewModelFieldType(config.DatabaseType, pkg, column)}
			if column.IsNullable {
				field.IsNullable = true
			}
			mt.Fields = append(mt.Fields, field)
			if column.IsPrimaryKey {
				mt.PKFields = append(mt.PKFields, field)
			}
		}
		mt.IsSimplePK = checkSimplePK(pkg, mt)
	}
	for _, view := range schema.Views {
		mt := &model.ViewModelType{TypeName: getModelTypeName(view.ViewName), PackageName: "viewmodel"}
		pkg.ViewModelTypes = append(pkg.ViewModelTypes, mt)
		for _, column := range view.Columns {
			field := &model.ViewModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getViewModelFieldType(config.DatabaseType, pkg, column)}
			if column.IsNullable {
				field.IsNullable = true
			}
			mt.Fields = append(mt.Fields, field)
		}
	}
	return pkg
}

func getViewModelFieldType(databaseType string, pkg *model.ViewModelPackage, column *model.Column) string {
	switch databaseType {
	case "mysql":
		return getMySQLViewModelFieldType(pkg, column)
	case "postgres":
		return getPostgresViewModelFieldType(pkg, column)
	default:
		return getMySQLViewModelFieldType(pkg, column)
	}
}

func getMySQLViewModelFieldType(pkg *model.ViewModelPackage, column *model.Column) string {
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
		ft = "[]byte" // sql/driver/Value does not supports bool
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for MySQL column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		ft = "string"
	}
	return ft
}

func getPostgresViewModelFieldType(pkg *model.ViewModelPackage, column *model.Column) string {
	var ft string = ""
	switch column.ColumnType {
	case "char", "varchar", "text", "character":
		ft = "string"
	case "bytea":
		ft = "[]byte"
	case "date", "time", "timetz", "timestamptz", "timestamp", "interval":
		ft = "time.Time"
		pkg.AppendImport("time")
	case "int2", "int4":
		ft = "int32"
	case "int8":
		ft = "int64"
	case "float4", "float8", "numeric":
		ft = "float64"
	case "bit", "bool":
		ft = "bool" // sql/driver/Value does not supports bool
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for Postgres column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		ft = "string"
	}
	return ft
}

func checkSimplePK(pkg *model.ViewModelPackage, mt *model.ViewModelType) bool {
	if len(mt.PKFields) == 1 {
		switch mt.PKFields[0].FieldType {
		case "string":
			mt.PKStringConv = "ret := value"
			mt.PKType = "string"
			return true
		case "int32":
			mt.PKStringConv = "ret1, _ := strconv.ParseInt(value, 10, 32); ret := int32(ret1)"
			pkg.AppendImport("strconv")
			mt.PKType = "int32"
			return true
		case "int64":
			mt.PKStringConv = "ret, _ := strconv.ParseInt(value, 10, 64)"
			pkg.AppendImport("strconv")
			mt.PKType = "int64"
			return true
		case "float64":
			mt.PKStringConv = "ret, _ := strconv.ParseFloat(value, 64)"
			pkg.AppendImport("strconv")
			mt.PKType = "float64"
			return true
		}
	}
	return false
}
