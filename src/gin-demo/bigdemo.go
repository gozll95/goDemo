package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/gin-gonic/gin.v1"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
	Age      int    `form:"age" json:"age"`
}

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before middleware")
		c.Set("request", "clinet_request")
		c.Next()
		fmt.Println("after middleware")
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "123" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}

func main() {
	router := gin.Default()

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"message": message,
			"nick":    nick,
		})
	})

	router.PUT("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
		})
	})

	router.LoadHTMLGlob("templates/*")
	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{"title": "index"})
	})

	//  curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/ghost/Desktop/pic.jpg" -H "Content-Type: multipart/form-data"
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
			//log.Fatal(err)
		}
		filename := header.Filename

		fmt.Println(file, err, filename)

		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusCreated, "upload successful")
	})
	// curl -X POST http://127.0.0.1:8000/multi/upload -F "upload=@/Users/ghost/Desktop/pic.jpg" -F "upload=@/Users/ghost/Desktop/journey.png" -H "Content-Type: multipart/form-data"
	router.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(200000)
		if err != nil {
			log.Fatal(err)
		}

		formdata := c.Request.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileheaders
		files := formdata.File["upload"] // grab the filenames

		for i, _ := range files { // loop through the files one by one
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			out, err := os.Create(files[i].Filename)

			defer out.Close()

			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				log.Fatal(err)
			}

			c.String(http.StatusCreated, "upload successful")

		}

	})

	v1 := router.Group("/v1")

	v1.Use(MiddleWare())

	v1.GET("/login", func(c *gin.Context) {
		fmt.Println(c.MustGet("request").(string))
		c.String(http.StatusOK, "v1 login")
	})

	v2 := router.Group("/v2")

	v2.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "v2 login")
	})

	router.POST("/login", func(c *gin.Context) {
		var user User
		var err error
		contentType := c.Request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err = c.BindJSON(&user)
		case "application/x-www-form-urlencoded":
			err = c.BindWith(&user, binding.Form)
		}
		//err = c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
			"passwd":   user.Passwd,
			"age":      user.Age,
		})

	})

	router.GET("/render", func(c *gin.Context) {
		contentType := c.DefaultQuery("content_type", "json")
		if contentType == "json" {
			c.JSON(http.StatusOK, gin.H{
				"user":   "rsj217",
				"passwd": "123",
			})
		} else if contentType == "xml" {
			c.XML(http.StatusOK, gin.H{
				"user":   "rsj217",
				"passwd": "123",
			})
		} else {
			c.YAML(http.StatusOK, gin.H{
				"user":   "rsj217",
				"passwd": "123",
			})
		}

	})

	router.GET("/redict/google", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://google.com")
	})

	router.GET("/before", MiddleWare(), func(c *gin.Context) {
		request := c.MustGet("request").(string)
		fmt.Println("before handler")
		c.JSON(http.StatusOK, gin.H{
			"middile_request": request,
		})
	})

	router.Use(MiddleWare())
	{
		router.GET("/middleware", func(c *gin.Context) {
			request := c.MustGet("request").(string)
			req, _ := c.Get("request")
			fmt.Println(req)
			c.JSON(http.StatusOK, gin.H{
				"middile_request": request,
				"request":         req,
			})
		})
	}

	router.GET("/after", func(c *gin.Context) {
		request := c.MustGet("request").(string)
		c.JSON(http.StatusOK, gin.H{
			"middile_request": request,
		})
	})

	router.GET("/auth/signin", func(c *gin.Context) {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "123",
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		c.String(http.StatusOK, "Login successful")
	})

	router.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "home"})
	})

	router.GET("/sync", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		log.Println("Done! in path" + c.Request.URL.Path)
	})

	router.GET("/async", func(c *gin.Context) {
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path" + cCp.Request.URL.Path)
		}()
	})

	router.Run(":8000")
}
