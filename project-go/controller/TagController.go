package controller

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"MYBLOGWEB/response"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ITagController interface { //定义一个接口
	ShowByTag(ctx *gin.Context)
}

func (T TagController) ShowByTag(ctx *gin.Context) { //通过标签查看笔记
	tags := ctx.Query("tag")
	userID, _ := strconv.Atoi(ctx.Query("userid"))

	//获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "3"))

	//开始分页
	var tags_db []model.Tag
	T.DB.Table("tags").Where("tag LIKE ?", "%"+tags+"%").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&tags_db)

	//创建返回的文章对象
	var articles = make([]model.ArticlesShow, len(tags_db))

	for i := 0; i < len(tags_db); i++ {
		article_id := tags_db[i].ArticleID

		var articles_db model.Article
		T.DB.Table("articles").Where("article_id=?", article_id).First(&articles_db)

		articles[i].ArticleID = articles_db.ArticleID

		create_user_id := articles_db.UserID
		var user model.User
		T.DB.Table("users").Where("user_id=?", create_user_id).First(&user)

		articles[i].UserName = user.Username
		articles[i].UserID = create_user_id
		articles[i].Title = articles_db.Title
		articles[i].UserName = user.Username
		articles[i].Content = articles_db.Content
		articles[i].CreatedAt = articles_db.CreatedAt.Format("2006-01-02 15:04:05")
		articles[i].Type = articles_db.Type
		articles[i].Favorites = articles_db.Favorites
		articles[i].Likes = articles_db.Likes

		var collect model.Collect

		if err := T.DB.Table("collects").Where("user_id=? AND article_id=?", userID, articles[i].ArticleID).First(&collect).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 查询结果为空
				articles[i].IsFavorites = false
			} else {
				// 查询出错
				panic(err)
			}
		} else {
			// 查询结果非空
			articles[i].IsFavorites = true
		}

		var like model.Like
		if err := T.DB.Table("likes").Where("user_id=? AND article_id=?", userID, articles[i].ArticleID).First(&like).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 查询结果为空
				articles[i].IsLike = false
			} else {
				// 查询出错
				panic(err)
			}
		} else {
			// 查询结果非空
			articles[i].IsLike = true
		}

		var comment_db []model.Comment
		T.DB.Table("comments").Where("article_id=?", articles[i].ArticleID).Order("created_at desc").Find(&comment_db)

		var comments_show = make([]model.ArticleCommentReturn, len(comment_db))
		for j := 0; j < len(comment_db); j++ {
			comments_show[j].UserId = comment_db[j].UserID

			Comment_User_ID := comment_db[j].UserID

			var comment_user model.User
			T.DB.Table("users").Select("*").Where("user_id=?", Comment_User_ID).First(&comment_user)

			comments_show[j].UserName = comment_user.Username
			comments_show[j].Message = comment_db[j].Message
			comments_show[j].CommentId = comment_db[j].CommentID
			comments_show[j].CreatedAt = comment_db[j].CreatedAt
		}

		var tag_db []model.Tag
		T.DB.Table("tags").Where("article_id=?", articles[i].ArticleID).Find(&tag_db)
		var tag_show = make([]model.TagReturn, len(tag_db))
		for k := 0; k < len(tag_show); k++ {
			tag_show[k].Tag = tag_db[k].Tag
		}

		articles[i].Comment = comments_show
		articles[i].Tag = tag_show

		var image_db []model.Image
		T.DB.Table("images").Where("article_id=?", articles[i].ArticleID).Find(&image_db)
		var image_show = make([]model.Image_return, len(image_db))
		for m := 0; m < len(image_db); m++ {
			image_show[m].ImageID = image_db[m].ImageID
			image_show[m].Image = image_db[m].Image
		}

		articles[i].Image = image_show
		articles[i].Avatar = user.Avatar
	}

	response.Success(ctx, gin.H{"data": articles}, "查看成功")
}

type TagController struct { //将打开的数据库DB返回到接口的每一个函数中
	DB *gorm.DB
}

func NewTagController() ITagController { //初始化该接口
	db := common.GetDB()
	db.AutoMigrate(&model.Tag{})
	return TagController{DB: db}
}
