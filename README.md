# DinGo
Data access in Go (DinGo). Generate Data Access Object (DAO) code from MySQL database schema.

![Dingo](doc/img/dingo_small.png)

The project is under development.

## Model Generation
DinGo generates Model Data Transfer Object (DTO) reading the MySql database schema.
All the generated strutcs are written in the *model.go* file inside the *model* package's directory.

This is an example of generated DTO:
```Go
// Data transfer object for Customer
type Customer struct {
	Id int64 `sql:"type:int(10) unsigned;not null;AUTO_INCREMENT"`
	Name string `sql:"type:varchar(60);not null"`
	State string `sql:"type:varchar(12);not null"`
	CreationDate time.Time `sql:"type:datetime;not null"`
	UpdateDate time.Time `sql:"type:datetime;not null"`
	
}
```
Every field has a *GORM style* metadata.

## DAO Generation
DinGo generates DAO structs reading the MySql database schema and the Model.
Every DAO defines these methods to perform CRUD operations on entities:
- Insert(conn *sql.DB, dto *model.ModelStruct)(lastInsertId int64, err error)
- Update(conn *sql.DB, dto *model.ModelStruct)(rowsAffected int64, err error)
- Delete(conn *sql.DB, dto *model.ModelStruct)(rowsAffected int64, err error)
- FindByPrimaryKey(conn *sql.DB, pk1 pk1Type, pk2 pk2Type, ...) (dto *model.ModelStruct, err error)
- List(conn *sql.DB, take int32, skip int32) (list []*model.ModelStruct, err error)

Generated DAO supports table's primary keys and auto-increment columns.
All the DAO generated structs are written in the *dao.go* file inside the *dao* package's directory.

This is an example of generated DAO struct:
```Go
// Data access object for Customer entities.
type CustomerDao struct {
	
}
// Insert a new Customer entity and returns the last insert Id.
func (dao *CustomerDao) Insert(conn *sql.DB, dto *model.Customer)(lastInsertId int64, err error){
	q := "INSERT INTO customer VALUES (?, ?, ?, ?, ?)"
	res, err := conn.Exec(q, sql.NullInt64{}, dto.Name, dto.State, dto.CreationDate, dto.UpdateDate)
    if err != nil {
		return -1, err
	}
	lastInsertId, err = res.LastInsertId()
	return lastInsertId, err
}
// Update a Customer entity and returns the number of affected rows.
func (dao *CustomerDao) Update(conn *sql.DB, dto *model.Customer)(rowsAffected int64, err error){
	...
}
// Delete a Customer entity and returns the number of affected rows.
func (dao *CustomerDao) Delete(conn *sql.DB, dto *model.Customer)(rowsAffected int64, err error){
	...
}
// Find the Customer entity by primary keys, returns nil if not found.
func (dao *CustomerDao) FindByPrimaryKey(conn *sql.DB, Id int64) (dto *model.Customer, err error){
	...
}
// List the Customer entities.
func (dao *CustomerDao) List(conn *sql.DB, take int32, skip int32) (list []*model.Customer, err error){
	...
}
```

## Configuration
The MySQL connection and other configuration parameters are defined in the *config.json* file.
Here is a configuration example:
```JSON
{
	"Hostname": "localhost", 
	"Port": "3306", 
	"DatabaseName": "Customers", 
	"Username": "zerbo", 
	"Password": "Mysql.2016",
	"BasePackage": "github.com/maxzerbini/prjtest",
	"OutputPath": "$GOPATH/src/github.com/maxzerbini/prjtest"
}
```

## Building DinGo
```bash
$ go get github.com/maxzerbini/dingo
$ go build -i github.com/maxzerbini/dingo
```

## Running DinGo
Make sure to properly set the config.json file with your connection parameters and run
```bash
$ dingo 
```
If you rename or move the configuration file then run
```bash
$ dingo -conf=/mypath/myconfig.json
```

## Known issues
- The DAO components are produced correctly if the tables have a PK
- Some columns types that are not recognized (such as JSON) are mapped to string fields
- DinGo maps DATE, TIME, DATETIME and TIMESTAMP column types to *time.Time* assuming that the connection has opened using the DSN parameter *parseTime=true*

## Warning
It's recommended to test the generated code before use in production.
