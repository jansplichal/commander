package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*Person type*/
type Person struct {
	Name  string
	Phone string
}

/*Connection is mongo db connection*/
type Connection struct {
	session *mgo.Session
	coll    *mgo.Collection
}

/*Connect to mongo DB*/
func Connect(server string) (*Connection, error) {
	session, err := mgo.Dial(server)
	if err != nil {
		return nil, err
	}
	conn := Connection{session: session}
	session.SetMode(mgo.Monotonic, true)

	return &conn, err
}

/*Close mongo DB connection*/
func (conn *Connection) Close() {
	conn.session.Close()
}

/*Use selects the database and collection
to use*/
func (conn *Connection) Use(database, collection string) {
	conn.coll = conn.session.DB(database).C(collection)
}

/*Insert a object into DB*/
func (conn *Connection) Insert(object ...interface{}) error {
	err := conn.coll.Insert(object...)
	return err
}

/*FindOne an object by bson*/
func (conn *Connection) FindOne(crit bson.M, result interface{}) (interface{}, error) {
	err := conn.coll.Find(crit).One(result)
	if err != nil {
		return nil, err
	}
	return result, err
}
