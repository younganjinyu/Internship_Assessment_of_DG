package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5430
	user     = "postgres"
	password = "postgres"
	dbName   = "assess"
)

func getDBEngine() *xorm.Engine {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	//格式
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	//打印生成的SQL语句
	engine.ShowSQL()

	err = engine.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("connect postgresql success")
	return engine
}

type SqlResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var sqlResponse SqlResponse

type Msg struct {
	Id      int
	U_name  string
	Title   string
	Date    string
	Content string
}

type User struct {
	Name   string
	Passwd int
}
type Comment struct {
	Id_comment int
	Content    string
	Id_msg     int
	Date       string
}

//查询所有用户
func SelectAllUser(c *gin.Context) {
	var u []User
	engine := getDBEngine()
	err := engine.Find(&u)
	if err != nil {
		fmt.Println("查询错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = u
	c.JSON(http.StatusOK, sqlResponse)
}

//根据用户名查询
func SelectUserByName(c *gin.Context) {
	var user User
	engine := getDBEngine()
	name := c.Query("name")
	_, err := engine.Where("name=?", name).Get(&user)
	if err != nil {
		fmt.Println("查询数据失败")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询数据失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	fmt.Println(user, "查询数据")
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = user
	c.JSON(http.StatusOK, sqlResponse)
}

//查询所有帖子
func SelectAllMsg(c *gin.Context) {
	var msg []Msg
	engine := getDBEngine()
	err := engine.Find(&msg)
	if err != nil {
		fmt.Println("查询错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = msg
	c.JSON(http.StatusOK, sqlResponse)
}

//按帖子id查询
func SelectMsgById(c *gin.Context) {
	var msg Msg
	id := c.Query("id")
	engine := getDBEngine()
	_, err := engine.Where("id=?", id).Get(&msg)
	if err != nil {
		fmt.Println("查询错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = msg
	c.JSON(http.StatusOK, sqlResponse)
}

//按标题和日期查询
func SelectMsgByTitleAndDate(c *gin.Context) {
	var msg Msg
	title := c.Query("title")
	date := c.Query("date")
	engine := getDBEngine()
	_, err := engine.Where("title=?", title).And("date=?", date).Get(&msg)
	if err != nil {
		fmt.Println("查询错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = msg
	c.JSON(http.StatusOK, sqlResponse)
}

//查询一个帖子所有回复
func SelectAllComment(c *gin.Context) {
	var comment []Comment
	engine := getDBEngine()
	id := c.Query("id")
	err := engine.Join("INNER", "msg", "comment.id_msg = msg.id").Where("id_msg=?", id).Find(&comment)
	if err != nil {
		fmt.Println("查询错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "查询失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = comment
	c.JSON(http.StatusOK, sqlResponse)
}

//发帖
func InsertMsg(c *gin.Context) {
	var msg *Msg
	engine := getDBEngine()
	err := c.BindJSON(&msg)
	if err != nil {
		fmt.Println("解析数据错误", err)
		sqlResponse.Code = 400
		sqlResponse.Message = "传递数据错误"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05")
	msg.Date = timeString
	rows, err := engine.Insert(msg)
	if err != nil || rows <= 0 {
		fmt.Println("插入数据错误", err)
		sqlResponse.Code = 400
		sqlResponse.Message = "插入数据错误"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Code = 200
	sqlResponse.Message = "插入数据成功"
	sqlResponse.Data = msg
	c.JSON(http.StatusOK, sqlResponse)
}

//回复帖子
func InsertComment(c *gin.Context) {
	var comment *Comment
	engine := getDBEngine()
	err := c.BindJSON(&comment)
	if err != nil {
		fmt.Println("解析数据错误", err)
		sqlResponse.Code = 400
		sqlResponse.Message = "传递数据错误"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05")
	comment.Date = timeString
	rows, err := engine.Insert(comment)
	if err != nil || rows <= 0 {
		fmt.Println("插入数据错误", err)
		sqlResponse.Code = 400
		sqlResponse.Message = "插入数据错误"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Code = 200
	sqlResponse.Message = "插入数据成功"
	sqlResponse.Data = comment
	c.JSON(http.StatusOK, sqlResponse)

}

//修改帖子
func UpdateMsg(c *gin.Context) {
	var msg *Msg
	id := c.Query("id")
	engine := getDBEngine()
	err := c.BindJSON(&msg)
	if err != nil {
		fmt.Println("解析数据错误", err)
		sqlResponse.Code = 400
		sqlResponse.Message = "传递数据错误"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05")
	msg.Date = timeString
	_, err_update := engine.Where("id=?", id).Update(msg)
	if err_update != nil {
		log.Println(err_update)
		return
	}
	sqlResponse.Code = 200
	sqlResponse.Message = "查询成功"
	sqlResponse.Data = msg
	c.JSON(http.StatusOK, sqlResponse)
}

//根据id删除msg
func DeleteMsgById(c *gin.Context) {
	engine := getDBEngine()
	id := c.Query("id")
	var msg Msg
	_, err := engine.Where("id=?", id).Get(&msg)
	fmt.Println(err)
	if err != nil {
		fmt.Println("删除错误")
		sqlResponse.Code = 400
		sqlResponse.Message = "删除失败,没有找到相应帖子"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	result, err_delete := engine.Exec("delete from msg where id=?", id)
	if err_delete != nil {
		log.Println(err_delete)
		sqlResponse.Code = 400
		sqlResponse.Message = "删除失败"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	rows, err_affected := result.RowsAffected()
	if err_affected == nil && rows > 0 {
		sqlResponse.Code = 1
		sqlResponse.Message = "success"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
}
