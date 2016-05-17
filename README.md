# DinGo
Data access in Go (DinGo). From database schema to RESTful API: all the code is generated for you in few seconds. 

![Dingo](doc/img/dingo_small.png)

## Build status
[![Build Status](https://drone.io/github.com/maxzerbini/dingo/status.png)](https://drone.io/github.com/maxzerbini/dingo/latest) [![Build Status](https://travis-ci.org/maxzerbini/dingo.svg?branch=master)](https://travis-ci.org/maxzerbini/dingo)

## Main features

DinGo creates a Microservice application starting from your database schema. 
Supported databases are MySQL and PostgreSQL (currently in beta).

These are the main steps followed by Dingo:
- Data Model generation
- Data Access Object (DAO) generation
- Business Object (Biz) generation
- View-Model generation
- Service Object generation
- Host Server creation
- JSON configuration file creation

The result you get is a Web API application that you just compile and run.

## Model Generation
DinGo generates Model Data Transfer Object (DTO) reading the database schema.
All the generated strutcs are written in the *model.go* file inside the *model* package's directory.

This is an example of generated DTO for MySql:
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
DinGo generates DAO structs reading the database schema and the Model.
Every DAO defines these methods to perform CRUD operations on entities:
- Insert(conn *sql.DB, dto *model.ModelStruct)(lastInsertId int64, err error)
- Update(conn *sql.DB, dto *model.ModelStruct)(rowsAffected int64, err error)
- Delete(conn *sql.DB, dto *model.ModelStruct)(rowsAffected int64, err error)
- FindByPrimaryKey(conn *sql.DB, pk1 pk1Type, pk2 pk2Type, ...) (dto *model.ModelStruct, err error)
- List(conn *sql.DB, take int32, skip int32) (list []*model.ModelStruct, err error)
- Count(conn *sql.DB)

Generated DAO supports table's primary keys and auto-increment columns.
All the DAO generated structs are written in the *dao.go* file inside the *dao* package's directory.

This is an example of generated DAO struct for MySql:
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
	q := "UPDATE customer SET Name=?, State=?, CreationDate=?, UpdateDate=?"
	q += " WHERE Id = ?"
	res, err := conn.Exec(q, dto.Name, dto.State, dto.CreationDate, dto.UpdateDate, dto.Id)
    if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	return rowsAffected, err
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
// Count the Customer entities.
func (dao *CustomerDao) Count(conn *sql.DB) (count int64, err error){
	...
}
```
The last two methods are available also for Views of the database.

## View-Model Generation
The View-Model objects are produced in a similar manner to those of the Model. 
The difference between these structures is that the former do not depend on the sql package and are designed to be serialized in JSON.

## Business Object Generation
These objects are wrapper around DAO objects. Their interface differ from the DAO interface because of View-Model types present on the method signatures.
This is an example of generated Biz object:
```Go
// Business object for Customer entities.
type CustomerBiz struct {
	Dao *dao.CustomerDao	
}
// Create a CustomerBiz
func NewCustomerBiz() *CustomerBiz {
	return &CustomerBiz{ Dao:&dao.CustomerDao{} }
}
// Convert an model entity in a view-model
func (b *CustomerBiz) ToViewModel(m *model.Customer) *viewmodel.Customer{
	v := &viewmodel.Customer{}
	v.Id = m.Id
	...
	return v
}
// Convert a view-model in a model entity
func (b *CustomerBiz) ToModel(v *viewmodel.Customer) *model.Customer{
	m := &model.Customer{}
	m.Id = v.Id
	...
	return m
}
// Insert a new Customer entity and returns the last insert Id.
func (b *CustomerBiz) Insert(v *viewmodel.Customer) (lastInsertId int64, err error) {
	return b.Dao.Insert(dao.Connection, b.ToModel(v))
}
// Update a Customer entity and returns the number of affected rows.
func (b *CustomerBiz) Update(v *viewmodel.Customer) (rowsAffected int64, err error) {
	return b.Dao.Update(dao.Connection, b.ToModel(v))
}
// Delete a Customer entity and returns the number of affected rows.
func (b *CustomerBiz) Delete(v *viewmodel.Customer) (rowsAffected int64, err error) {
	return b.Dao.Delete(dao.Connection, b.ToModel(v))
}
// Find the Customer entity by primary keys, returns nil if not found.
func (b *CustomerBiz) Find(mv *viewmodel.Customer) (v *viewmodel.Customer, err error){
	...
}
// List the Customer entities.
func (b *CustomerBiz) List(take int32, skip int32) (list []*viewmodel.Customer, err error) {
	...
}
// Count the Customer entities.
func (b *CustomerBiz) Count() (count int64, err error){
	return b.Dao.Count(dao.Connection)
}
```
## Service Object Generation
These objects offer a set of methods used to construct and expose RESTful API.

## REST API Generation
DinGo can generate the set of RESTful API endpoints needed to perform CRUD operations on entities. Each entity corresponds to a resource, and each resource has the necessary endpoints to be managed. 
If the resource has a simple identifier (the corresponding entity has one-column PK) then Dingo generates these endpoints:
- *GET [basehost]/resourcename?skip=[value]&take=[value]&count=all* lists the elements
- *POST [basehost]/resourcename* creates a new element
- *PUT [basehost]/resourcename/:id* updates the element has PK = id
- *DELETE [basehost]/resourcename/:id* deletes an element
- *GET [basehost]/resourcename/:id* finds an element by primary key

If the resource has a complex identifier and the corresponding entity has a PK made of more columns, the generated endpoints are these:
- *GET [basehost]/resourcename?skip=[value]&take=[value]* lists the elements
- *POST [basehost]/resourcename* creates a new element
- *PUT [basehost]/resourcename* updates an element
- *POST [basehost]/resourcename/delete* deletes an element
- *POST [basehost]/resourcename/find* finds an element by primary keys
- *GET [basehost]/resourcename/count* counts the etities


## Building DinGo
```bash
$ go get github.com/maxzerbini/dingo
$ go build -i github.com/maxzerbini/dingo
```

## Running DinGo
Make sure to properly set the *config.json* file with your connection parameters and run
```bash
$ dingo 
```
If you rename or move the configuration file then run
```bash
$ dingo -conf=/mypath/myconfig.json
```

## DinGo Configuration
The MySQL connection and other configuration parameters are defined in the *config.json* file.
Here is a configuration example:
```JSON
{
	"Hostname": "localhost", 
	"Port": "3306", 
	"DatabaseType": "MySQL",
	"DatabaseName": "Customers", 
	"Username": "zerbo", 
	"Password": "Mysql.2016",
	"BasePackage": "github.com/maxzerbini/prjtest",
	"OutputPath": "$GOPATH/src/github.com/maxzerbini/prjtest",
	"ExcludedEntities": [],
	"Entities": [],
	"SkipDaoGeneration": false,
	"SkipBizGeneration": false,
	"SkipServiceGeneration": false,
	"ForcePluralResourceName": true,
	"PostgresSchema": "public"
}
```
The _DatabaseType_ can assume one of these values
- "MySQL"
- "Postgres"

Optional configuration parameters
- _ExcludedEntities_ is an optional list of enity names that will be exluded
- _Entities_ is a list of included entities, if it's void all the entities are considered
- _SkipDaoGeneration_ skip DAO generation step and the following steps
- _SkipBizGeneration_ skip Biz generation step and the following steps
- _SkipServiceGeneration_ skip Service Object generation step and the following steps
- _ForcePluralResourceName_ force english plural names for resource endpoints
- _PostgresSchema_ is the PostgreSQL Schema name (usually "public")

## Using generated DAO and Biz code
It's very easy using generated code. Here an example:
```Go
// open the connection
conn, err := sql.Open("mysql","myuser:password@tcp(localhost:3306)/myDatabase?parseTime=true")
if err != nil {
	panic(err)
} else {
	// set the connection
	dao.Connection = conn
}
b := biz.NewCustomerBiz()
cust := &viewmodel.Customer{Name: "Max", State: "PENDING", CreationDate: time.Now(), UpdateDate: time.Now()}
// insert a new customer
id, err := b.Insert(cust)
// get the list of customers
result, err := b.List(100, 0)
// find entity by primary key
cust2 := &viewmodel.Customer{Id: 7}
cust2, err = b.Find(cust2)
```

## Extension points
All generated structs can be extended creating custom methods. 
For example if you want add a custom method to a *CustomerBiz*, you just create a file *bitext.go* in the *biz* directory and put in the new method
```Go
package biz

import "github.com/myuser/myproject/model"
import "github.com/myuser/myproject/dao"
import "github.com/myuser/myproject/viewmodel"

func (b *CustomerBiz) InsertCustomerAndProduct(v *viewmodel.Customer, idProduct int64)(lastInsertId int64, err error) {
	lastInsertId, err = b.Dao.Insert(dao.Connection, b.ToModel(v))
	cpdao := &dao.CustomerprodutcDao{}
	cp := &model.Customerproduct{IdCustomer:lastInsertId, IdProduct:idProduct, UpdateDate:time.Now()}
	id, err := cpdao.Insert(dao.Connection, cp)
	return lastInsertId, err
}
```
You can do this with all the objects generated by DinGo such as Model, Dao, Biz, ViewModel and Service Objects.
If you need to add endpoints to the API then you can register new endpoints in the file *customresources.go* inside the method *registerCustomResources*
```Go
package main

import "github.com/gin-gonic/gin"

// Register custom resource endpoints here.
func registerCustomResources(conf Configuration, router *gin.Engine) {
	// TODO add your custom endpoints
	// so := service.NewMyService()
	// router.GET(conf.WebBaseHost+"/myresources", so.MyMethod)
}
```
The file *customresources.go* is generated only once then is no longer rewritten so that it can be modified without the risk of losing changes.

## Known issues
- The DAO components are produced correctly if the tables have a PK
- The Update method is not implemented if the entity does not have at least one column not PK
- Some columns types that are not recognized (such as JSON) are mapped to string fields
- PostgreSQL array types are not supported (these columns are all mapped as string)
- DinGo maps DATE, TIME, DATETIME and TIMESTAMP column types to *time.Time* assuming that the connection has opened using the DSN parameter *parseTime=true*
- If you have a lot of entities in your database, you could produce a *"SOA Monolith"*, but using the configuration parameters _ExcludedEntities_ or _Entities_ and changing the _BasePackage_ you can limit the number of endpoints and you can produce many small applications, obtaining a set of indipendent Microservices
- Some HTTP verbs such as DELETE are not used defining the service endpoints of the resources that have complex primary keys

## Warning
It's recommended to test the generated code before using it in production. 
If you find a problem feel free to submit an issue.

## Roadmap
### Supported Databases
- [x] MySQL 
- [x] PostgreSQL (beta)
- [ ] SQLite
- [ ] MS SQL Server
