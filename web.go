package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

type Article struct {
	Id        int    `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"tontent"`
}

func initDb()

func main() {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	defer db.Close()
	if err != nil{
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
 
	if err := db.Ping(); err != nil{
		log.Fatalln(err)
	}


	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})
	
	router.GET("/articlelist", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, content FROM article")
		defer rows.Close()

		if err != nil {
			log.Fatalln(err)
		}

		articles := make([]Article, 0)

		for rows.Next() {
			var article Article
			rows.Scan(&article.Id, &article.Title, &article.Content)
			articles = append(articles, article)
		}
		if err = rows.Err(); err != nil {
			log.Fatalln(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"articles": articles,
		})
	})
	
	router.POST("/articleadd", func(c *gin.Context) {
		Title := c.Request.FormValue("title")
		Content := c.Request.FormValue("content")

		rs, err := db.Exec("INSERT INTO article(title, content) VALUES (?, ?)", Title, Content)
		if err != nil {
			log.Fatalln(err)
		}

		id, err := rs.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("insert article Id {}", id)
		msg := fmt.Sprintf("insert successful %d", id)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})
	
	router.GET("/article/:id", func(c *gin.Context) {
		id := c.Param("id")
		var article Article
		err := db.QueryRow("SELECT id, title, content FROM article WHERE id=?", id).Scan(
		&article.Id, &article.Title, &article.Content,
		)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"article": nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"article": article,
		})
	})
	
	router.POST("/articleupdate", func(c *gin.Context) {
		cId :=c.Request.FormValue("id")
		Title := c.Request.FormValue("title")
		Content := c.Request.FormValue("content")

		stmt, err := db.Prepare("UPDATE article SET title=?, content=? WHERE id=?")
		defer stmt.Close()
		if err != nil {
			log.Fatalln(err)
		}
		rs, err := stmt.Exec(Title, Content, cId)
		if err != nil {
			log.Fatalln(err)
		}
		ra, err := rs.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		msg := fmt.Sprintf("Update article %d successful %d", cId, ra)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})
	
	router.GET("/articledelete/:id", func(c *gin.Context) {
		cid := c.Param("id")
		
		rs, err := db.Exec("DELETE FROM article WHERE id=?", cid)
		if err != nil {
			log.Fatalln(err)
		}
		ra, err := rs.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		msg := fmt.Sprintf("Delete article %d successful %d", cid, ra)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})
	
	router.Run(":8000")
}
