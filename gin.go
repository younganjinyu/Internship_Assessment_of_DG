package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/SelectAllUser", SelectAllUser)
	router.GET("/SelectAllMsg", SelectAllMsg)
	router.GET("/SelectUserByName", SelectUserByName)
	router.GET("/SelectMsgById", SelectMsgById)
	router.GET("SelectMsgByTitleAndDate", SelectMsgByTitleAndDate)
	router.GET("/DeleteMsgById", DeleteMsgById)
	router.GET("/SelectAllComment", SelectAllComment)
	router.POST("/InsertMsg", InsertMsg)
	router.POST("/UpdateMsg", UpdateMsg)
	router.POST("/InsertComment", InsertComment)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
