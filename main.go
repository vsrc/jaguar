package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gin-gonic/gin"
)

var (
	jsondb *db.DB
	col    *db.Col
)

const (
	dbname = "db"
	clname = "jaguar"
)

func main() {

	jsondb, err := db.OpenDB(dbname)
	if err != nil {
		fmt.Println(err)
	}
	defer jsondb.Close()

	if _, err = os.Stat(filepath.Join(dbname, clname)); os.IsNotExist(err) {
		if err = jsondb.Create(clname); err != nil {
			panic(err)
		}
	}
	col = jsondb.Use(clname)

	r := gin.Default()
	r.GET("/", GetPeople)
	r.GET("/:id", Get)
	r.POST("/", Create)
	r.PATCH("/:id", Append)
	r.PUT("/:id", Update)
	r.DELETE("/:id", Delete)

	r.Run(":8080")
}

// Delete deletes record by id
func Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	uid, _ := strconv.Atoi(id)
	col.Delete(uid)

	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

// Update updates record by id
func Update(c *gin.Context) {

	var record map[string]interface{}
	id := c.Params.ByName("id")

	uid, _ := strconv.Atoi(id)

	var err error
	if _, err = col.Read(uid); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&record)

	err = col.Update(uid, record)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println("Unable to update data")
	}

	c.JSON(200, record)

}

// Append updates record by id
func Append(c *gin.Context) {

	id := c.Params.ByName("id")

	uid, _ := strconv.Atoi(id)

	var readBack map[string]interface{}
	var err error
	if readBack, err = col.Read(uid); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&readBack)

	err = col.Update(uid, readBack)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println("Unable to update data")
	}

	c.JSON(200, readBack)

}

// Create adds record to db
func Create(c *gin.Context) {

	var record map[string]interface{}
	c.BindJSON(&record)

	docID, err := col.Insert(record)
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	fmt.Println("docId", docID)

	c.JSON(200, docID)
}

// Get returns specific record by id
func Get(c *gin.Context) {
	id := c.Params.ByName("id")

	uid, _ := strconv.Atoi(id)

	var readBack map[string]interface{}
	var err error
	if readBack, err = col.Read(uid); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	fmt.Println("read data", readBack)

	c.JSON(200, readBack)

}

// GetPeople returns all records
func GetPeople(c *gin.Context) {

	rec := make([]map[string]interface{}, 0)

	col.ForEachDoc(func(id int, doc []byte) bool {
		var entry map[string]interface{}
		json.Unmarshal(doc, &entry)
		entry["id"] = string(strconv.AppendInt(nil, int64(id), 10))
		rec = append(rec, entry)
		return true
	})

	c.JSON(200, rec)
}
