package producers

import (
	"bytes"
	"log"
	"strings"

	"github.com/maxzerbini/dingo/model"
)

func ProduceModelPackage(config *model.Configuration, schema *model.DatabaseSchema) (pkg *model.ModelPackage) {
	pkg = &model.ModelPackage{PackageName: "model", BasePackage: config.BasePackage}
	for _, table := range schema.Tables {
		mt := &model.ModelType{TypeName: getModelTypeName(table.TableName), PackageName: "model"}
		pkg.ModelTypes = append(pkg.ModelTypes, mt)
		for _, column := range table.Columns {
			field := &model.ModelField{FieldName: getModelFieldName(column.ColumnName), FieldType: getModelFieldType(pkg, column), FieldMetadata: getFieldMetadata(pkg, column)}
			mt.Fields = append(mt.Fields, field)
		}
	}
	return pkg
}

func getModelTypeName(tablename string) string {
	name := strings.Replace(tablename, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)
	return name
}

func getModelFieldName(fieldname string) string {
	name := strings.Replace(fieldname, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)
	return name
}

func getModelFieldType(pkg *model.ModelPackage, column *model.Column) string {
	var ft string = ""
	switch column.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		if column.IsNullable {
			ft = "sql.NullString"
			pkg.AppendImport("database/sql")
		} else {
			ft = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		ft = "[]byte"
	case "date", "time", "datetime", "timestamp":
		ft = "time.Time"
		pkg.AppendImport("time")
	case "tinyint", "smallint":
		if column.IsNullable {
			ft = "sql.NullInt32"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int32"
		}
	case "int", "mediumint", "bigint":
		if column.IsNullable {
			ft = "sql.NullInt64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "int64"
		}
	case "float", "decimal", "double":
		if column.IsNullable {
			ft = "sql.NullFloat64"
			pkg.AppendImport("database/sql")
		} else {
			ft = "float64"
		}
	case "bit":
		ft = "bool"
	}
	if ft == "" {
		log.Printf("WARNING Incompatible Go type for column %s %s -> using string\r\n", column.ColumnName, column.ColumnType)
		ft = "string"
	}
	return ft
}

// Generate GORM like metadata
func getFieldMetadata(pkg *model.ModelPackage, column *model.Column) string {
	var buffer bytes.Buffer
	buffer.WriteString("sql:\"")
	buffer.WriteString("type:")
	buffer.WriteString(column.ColumnType)
	if !column.IsNullable {
		buffer.WriteString(";not null")
	}
	if column.IsUnique {
		buffer.WriteString(";unique")
	}
	if column.IsAutoIncrement {
		buffer.WriteString(";AUTO_INCREMENT")
	}
	buffer.WriteString("\"")
	return buffer.String()
}
