package function

//无用代码
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var blogs []Blog

func main01() {
	router := gin.Default()

	// 前端从后端获取所有blog（可以拓展实现标签筛选功能）
	router.GET("/blogs", func(c *gin.Context) {
		c.JSON(http.StatusOK, blogs)
	})

	// 前端从后端获取id的blog
	router.GET("/blogs/:id", func(c *gin.Context) {
		id := getBlogID(c)
		if id == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid blog ID"})
			return
		}

		for _, blog := range blogs {
			if blog.ID == id {
				c.HTML(http.StatusOK, "blog.html", gin.H{
					"title":   blog.Title,
					"content": blog.Content,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"})
	})

	// 前端post新博客给后端
	router.POST("/blogs", func(c *gin.Context) {
		var blog Blog
		if err := c.ShouldBindJSON(&blog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
			return
		}
		//存储blog
		blog.ID = len(blogs) + 1
		blogs = append(blogs, blog)

		c.JSON(http.StatusOK, blog)
	})

	// Update a blog
	router.PUT("/blogs/:id", func(c *gin.Context) {
		id := getBlogID(c)
		if id == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid blog ID"})
			return
		}

		var updatedBlog Blog
		if err := c.ShouldBindJSON(&updatedBlog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
			return
		}

		for i, blog := range blogs {
			if blog.ID == id {
				updatedBlog.ID = id
				blogs[i] = updatedBlog
				c.JSON(http.StatusOK, updatedBlog)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"})
	})

	// Delete a blog
	router.DELETE("/blogs/:id", func(c *gin.Context) {
		id := getBlogID(c)
		if id == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid blog ID"})
			return
		}

		for i, blog := range blogs {
			if blog.ID == id {
				blogs = append(blogs[:i], blogs[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"})
	})

	router.LoadHTMLFiles("blog.html")
	router.Run(":8080")
}

func getBlogID(c *gin.Context) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return -1
	}
	return id
}
