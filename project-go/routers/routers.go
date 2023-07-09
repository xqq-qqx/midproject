package routers

import (
	"MYBLOGWEB/controller"

	"github.com/gin-gonic/gin"
)

func CollectRouter(router *gin.Engine) *gin.Engine {
	api := router.Group("api")

	PersonController := controller.NewPersonController()
	//用户注册
	api.POST("/register", PersonController.Register)
	//用户登录
	api.POST("/login", PersonController.Login)
	//用户访问个人中心
	api.GET("/getuserinfo", PersonController.Home)
	//用户编辑个人数据
	api.POST("/edituserinfo", PersonController.EditPersonInfo)
	//用户编辑个人头像
	api.POST("/editAvatar", PersonController.EditAvatar)
	comment := router.Group("comment")
	CommentController := controller.NewCommentController()
	//用户添加评论
	comment.POST("/postcomment", CommentController.AddComment)
	//用户删除评论
	comment.DELETE("/deletecomment", CommentController.DeleteComment)

	//操作文章
	articleRoutes := router.Group("/articles")
	articleController := controller.NewArticleController()

	//文章的四个函数的路由
	articleRoutes.POST("/create", articleController.Create) //创建文章
	//articleRoutes.PUT("/", articleController.Update)            //通过article_id来更新文章
	articleRoutes.GET("/show/details", articleController.Show)     //通过article_id来查看文章
	articleRoutes.DELETE("/delete", articleController.Delete)      //通过article_id来删除文章
	articleRoutes.GET("/all", articleController.GetAllArticle)     //返回所有的文章信息
	articleRoutes.GET("/showbytype", articleController.ShowByType) //通过type来查看文章信息

	//通过标签查看文章
	tagRoutes := router.Group("articles")
	tagController := controller.NewTagController()
	tagRoutes.GET("/showbytag", tagController.ShowByTag)

	//收藏与取消收藏
	router.POST("/collectarticle", controller.AddToFavorites)
	router.DELETE("/cancelcollect", controller.DeleteFromFavorites)

	//获取我的收藏
	router.GET("/api/find/favorite/article", controller.GetPersonalFavorites)

	//点赞和取消点赞
	router.POST("/likearticle", controller.AddToLikes)
	router.DELETE("/cancellike", controller.DeleteFromLikes)

	//关注和取消关注
	router.POST("/followuser", controller.Follow)
	router.DELETE("/unfollowuser", controller.UnFollow)

	//获取个人发布的所有文章
	router.GET("/api/find/article", controller.GetPersonalArticles)

	//获取所关注的用户
	router.GET("/api/get/follow", PersonController.GetFollow)

	//获取其它用户的详细信息
	router.GET("/api/get/follow/info", PersonController.ShowOtherUser)

	//获取其它用户发布的所有文章
	router.GET("/api/get/follow/articles", articleController.ShowOtherUserArticle)
	return router
}
