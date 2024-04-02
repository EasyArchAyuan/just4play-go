package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// gin 框架中采用的路由库是基于httprouter做的
// 地址为：https://github.com/julienschmidt/httprouter
// 原理：httprouter会根据规则构造前缀树
func main() {
	//1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	//routes group
	v1 := r.Group("/v1")
	{
		v1.GET("/login", login)
		v1.GET("submit", submit)
	}

	v2 := r.Group("v2")
	{
		v2.POST("/login", login)
		v2.POST("submit", submit)
	}

	//2.绑定路由规则
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello World!")
	})

	//可以通过Context的Param方法来获取API参数
	//localhost:8000/xxx/Krin
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		//截取/
		action = strings.Trim(action, "/")
		c.String(http.StatusOK, name+" is "+action)
	})

	//URL参数可以通过DefaultQuery()或Query()方法获取
	//DefaultQuery()若参数不存在，返回默认值，Query()若不存在，返回空串
	//API ? name=zs
	r.GET("/user", func(context *gin.Context) {
		name := context.DefaultQuery("name", "Ayuan")
		context.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	})

	//表单参数可以通过PostForm()方法获取，该方法默认解析的是x-www-form-urlencoded或from-data格式的参数
	//表单传输为post请求，http常见的传输格式为四种：
	//application/json
	//application/x-www-form-urlencoded
	//application/xml
	//multipart/form-data
	r.POST("/form", func(c *gin.Context) {
		types := c.DefaultPostForm("type", "post")
		username := c.PostForm("username")
		password := c.PostForm("userpassword")
		c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
	})

	//文件上传
	// 限制表单上传大小 8MB，默认为32MB
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err != nil {
			context.String(500, "上传图片出错哦")
		}
		_ = context.SaveUploadedFile(file, file.Filename)
		context.String(http.StatusOK, file.Filename)
	})

	//上传特定文件
	r.POST("/uploadPng", func(c *gin.Context) {
		_, headers, err := c.Request.FormFile("file")
		if err != nil {
			log.Printf("Error when try to get file: %v", err)
		}
		//headers.Size 获取文件大小
		if headers.Size > 1024*1024*2 {
			fmt.Println("文件太大了")
			return
		}
		//headers.Header.Get("Content-Type")获取上传文件的类型
		if headers.Header.Get("Content-Type") != "image/png" {
			fmt.Println("只允许上传png图片")
			return
		}
		c.SaveUploadedFile(headers, "./video/"+headers.Filename)
		c.String(http.StatusOK, headers.Filename)
	})

	//上传多个文件
	r.POST("/uploadFiles", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		//获取所有图片
		files := form.File["files"]
		//遍历所有图片
		for _, file := range files {
			if err := c.SaveUploadedFile(file, file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}
		}
		c.String(200, fmt.Sprintf("upload ok %d files", len(files)))
	})

	//3.监听端口
	r.Run(":8888")

}

func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}
