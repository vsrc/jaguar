package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	conf   map[string]interface{}
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

	read, err := ioutil.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(read, &conf)

	r := gin.Default()
	r.GET("/", Get)
	r.GET("/:id", GetOne)
	r.POST("/", Create)
	r.POST("/bulk", BulkCreate)
	r.PATCH("/:id", Append)
	r.PUT("/:id", Update)
	r.DELETE("/bulk", BulkDelete)
	r.DELETE("/single/:id", Delete)

	r.Run(":" + conf["port"].(string))
}

// BulkDelete deletes record by id
func BulkDelete(c *gin.Context) {

	decoder := json.NewDecoder(c.Request.Body)

	var records []int

	err := decoder.Decode(&records)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
	}
	defer c.Request.Body.Close()

	counter := 0
	for i := 0; i < len(records); i++ {
		col.Delete(records[i])
		counter = counter + 1
	}

	c.JSON(200, strconv.Itoa(counter)+" records deleted")
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
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Println("docId", docID)

	c.JSON(200, record)
}

// BulkCreate adds multiple records to db
func BulkCreate(c *gin.Context) {

	type Record map[string]interface{}
	decoder := json.NewDecoder(c.Request.Body)

	var records []Record

	err := decoder.Decode(&records)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
	}
	defer c.Request.Body.Close()

	count := 0

	for i := 0; i < len(records); i++ {

		_, err := col.Insert(records[i])
		if err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}

		count = count + 1
	}

	c.JSON(200, strconv.Itoa(count)+" records saved")
}

// GetOne returns specific record by id
func GetOne(c *gin.Context) {
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

// Get returns all records
func Get(c *gin.Context) {

	limit := int(conf["limit"].(float64))
	offset := 0
	qp := c.Request.URL.Query()

	if len(qp["limit"]) > 0 {
		l, err := strconv.Atoi(qp["limit"][0])

		if err == nil {
			limit = l
		}
	}

	if len(qp["offset"]) > 0 {
		o, err := strconv.Atoi(qp["offset"][0])

		if err == nil {
			offset = o
		}
	}

	log.Println("limit is: " + strconv.Itoa(limit))
	log.Println("offset is: " + strconv.Itoa(offset))

	rec := make([]map[string]interface{}, 0)

	i := 0
	col.ForEachDoc(func(id int, doc []byte) bool {
		i++

		if i > offset && i <= (offset+limit) {

			var entry map[string]interface{}
			json.Unmarshal(doc, &entry)
			entry["id"] = string(strconv.AppendInt(nil, int64(id), 10))
			rec = append(rec, entry)
			return true
		}
		return true
	})

	c.JSON(200, rec)
}
