package misakago

import (
//"gopkg.in/mgo.v2/bson"
)

//基础数据库模板
type BaseModel struct {
	Id_ string `bson:"_id"`
	//Timestamp   int
	//#Last_modify int
}
