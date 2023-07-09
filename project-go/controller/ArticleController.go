package controller

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/model"
	"MYBLOGWEB/response"
	"MYBLOGWEB/vo"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

type IArticleController interface { //定义一个接口
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Show(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAllArticle(ctx *gin.Context)
	ShowByType(ctx *gin.Context)
	ShowOtherUserArticle(ctx *gin.Context)
}

type ArticleController struct { //将数据库返回到接口的所有方法中
	DB *gorm.DB
}

func NewArticleController() IArticleController { //初始化该接口
	db := common.GetDB()
	return ArticleController{DB: db}
}

func (c ArticleController) Create(ctx *gin.Context) { //创建一个笔记
	userid, _ := strconv.Atoi(ctx.PostForm("userid"))

	article := model.Article{
		UserID:    uint(userid),
		Title:     ctx.PostForm("title"),
		Content:   ctx.PostForm("content"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      ctx.PostForm("type"),
		Favorites: 0,
		Likes:     0,
	}

	//在数据库中创建一个条目
	if err := c.DB.Table("articles").Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create!",
		})
		return
	}

	tagValues := ctx.PostFormArray("tag")
	fmt.Println(tagValues)
	tags := make([]model.Tag, len(tagValues))

	for i := 0; i < len(tagValues); i++ {
		tags[i].Tag = tagValues[i]
		tags[i].ArticleID = article.ArticleID

		//在数据库中创建一个条目
		if err := c.DB.Create(&tags[i]).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create!",
			})
			return
		}
	}
	// 处理上传的图片
	images := ctx.Request.MultipartForm.File["images"]
	if len(images) > 9 {
		response.Failed(ctx, nil, "最多只能上传9张图片")
		return
	}

	var imageUrls []string // 保存所有上传图片的 URL 地址

	for _, image := range images {
		// 调用 picture 函数将图片保存到本地，并返回图片的 URL 地址
		image_url := common.Picture(image)
		if image_url == "" {
			response.Failed(ctx, nil, "存储图片失败!")
			return
		}
		imageUrls = append(imageUrls, image_url)
	}

	// 将所有上传图片的 URL 地址保存到数据库中
	for _, imageUrl := range imageUrls {
		imageRecord := model.Image{
			Image:     imageUrl,
			ArticleID: article.ArticleID,
		}
		if err := c.DB.Create(&imageRecord).Error; err != nil {
			response.Failed(ctx, nil, "保存图片信息失败")
			return
		}
	}

	response.Success(ctx, nil, "")
}

func (c ArticleController) Update(ctx *gin.Context) { //更新笔记
	//绑定body的参数
	var requestArticle vo.CreateArticleRequest
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		response.Failed(ctx, nil, "数据验证错误!")
		return
	}

	//获取path中的参数
	articleId, _ := strconv.Atoi(ctx.Params.ByName("id"))

	var updateArticle model.Article
	if c.DB.Table("articles").Select("*").Where("article_id=?", articleId).First(&updateArticle).RecordNotFound() {
		response.Failed(ctx, nil, "文章不存在")
		return
	}

	//更新分类
	err := c.DB.Model(&updateArticle).Where("article_id=?", articleId).Update(requestArticle)
	if err.Error != nil {
		response.Failed(ctx, nil, "更新失败!")
		return
	}

	response.Success(ctx, gin.H{"articles": updateArticle}, "更新成功!")

}

func (c ArticleController) Show(ctx *gin.Context) { //展示笔记
	articleID := ctx.Query("articleid")
	//获取登录人的ID
	UserID, _ := strconv.Atoi(ctx.Query("userid"))

	var articles_db model.Article
	if err := c.DB.Table("articles").Where("article_id=?", articleID).First(&articles_db).Error; err != nil {
		response.Failed(ctx, nil, "没有找到相关的数据！")
	}

	var article = model.ArticlesShow{}

	var create_article_user model.User
	c.DB.Table("users").Select("*").Where("user_id=?", articles_db.UserID).First(&create_article_user)

	article.ArticleID = articles_db.ArticleID
	article.UserID = articles_db.UserID
	article.Title = articles_db.Title
	article.UserName = create_article_user.Username
	article.Content = articles_db.Content
	article.CreatedAt = articles_db.CreatedAt.Format("2006-01-02 15:04:05")
	article.Type = articles_db.Type
	article.Favorites = articles_db.Favorites
	article.Likes = articles_db.Likes

	var collect model.Collect

	if err := c.DB.Table("collects").Where("user_id=? AND article_id=?", UserID, articleID).First(&collect).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询结果为空
			article.IsFavorites = false
		} else {
			// 查询出错
			panic(err)
		}
	} else {
		// 查询结果非空
		article.IsFavorites = true
	}

	var like model.Like
	if err := c.DB.Table("likes").Where("user_id=? AND article_id=?", UserID, articleID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询结果为空
			article.IsLike = false
		} else {
			// 查询出错
			panic(err)
		}
	} else {
		// 查询结果非空
		article.IsLike = true
	}

	var comment_db []model.Comment
	c.DB.Table("comments").Where("article_id=?", article.ArticleID).Order("created_at desc").Find(&comment_db)

	var comments_show = make([]model.ArticleCommentReturn, len(comment_db))
	for j := 0; j < len(comment_db); j++ {
		comments_show[j].UserId = comment_db[j].UserID

		Comment_User_ID := comment_db[j].UserID

		var comment_user model.User
		c.DB.Table("users").Select("*").Where("user_id=?", Comment_User_ID).First(&comment_user)

		comments_show[j].UserName = comment_user.Username
		comments_show[j].Message = comment_db[j].Message
		comments_show[j].CommentId = comment_db[j].CommentID
		comments_show[j].CreatedAt = comment_db[j].CreatedAt
	}

	var tag_db []model.Tag
	c.DB.Table("tags").Where("article_id=?", article.ArticleID).Find(&tag_db)
	var tag_show = make([]model.TagReturn, len(tag_db))
	for k := 0; k < len(tag_show); k++ {
		tag_show[k].Tag = tag_db[k].Tag
	}

	article.Comment = comments_show
	article.Tag = tag_show

	var image_db []model.Image
	c.DB.Table("images").Where("article_id=?", article.ArticleID).Find(&image_db)
	var image_show = make([]model.Image_return, len(image_db))
	for m := 0; m < len(image_db); m++ {
		image_show[m].ImageID = image_db[m].ImageID
		image_show[m].Image = image_db[m].Image
	}

	article.Image = image_show
	article.Avatar = create_article_user.Avatar

	var follow model.Follow
	if err := c.DB.Table("follows").Where("user_id=? AND followed_user_id=?", UserID, create_article_user.UserID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询结果为空
			article.IsFollowed = false
		} else {
			// 查询出错
			panic(err)
		}
	} else {
		// 查询结果非空
		article.IsFollowed = true
	}

	response.Success(ctx, gin.H{"data": article}, "")
}

func (c ArticleController) Delete(ctx *gin.Context) { //删除笔记
	articleId := ctx.Query("articleid")
	var images []model.Image

	c.DB.Table("images").Select("*").Where("article_id=?", articleId).Find(&images)

	for i := 0; i < len(images); i++ {
		path := images[i].Image
		os.Remove(path)
	}

	if err := c.DB.Table("articles").Select("*").Where("article_id=?", articleId).Delete(&model.Article{}).Error; err != nil {
		response.Failed(ctx, nil, "删除失败，请重试！")
		return
	}

	response.Success(ctx, nil, "删除成功")

}

func (c ArticleController) GetAllArticle(ctx *gin.Context) {
	//获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "3"))

	//获取userid
	UserID, _ := strconv.Atoi(ctx.Query("userid"))

	//开始分页
	var articles_db []model.Article
	c.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles_db)

	var articles = make([]model.ArticlesShow, len(articles_db))

	for i := 0; i < len(articles_db); i++ {
		create_user_id := articles_db[i].UserID

		var user model.User
		c.DB.Table("users").Select("*").Where("user_id=?", create_user_id).First(&user)

		articles[i].ArticleID = articles_db[i].ArticleID
		articles[i].UserID = create_user_id
		articles[i].Title = articles_db[i].Title
		articles[i].UserName = user.Username
		articles[i].Content = articles_db[i].Content
		articles[i].CreatedAt = articles_db[i].CreatedAt.Format("2006-01-02 15:04:05")
		articles[i].Type = articles_db[i].Type
		articles[i].Favorites = articles_db[i].Favorites
		articles[i].Likes = articles_db[i].Likes

		var collect model.Collect

		if err := c.DB.Table("collects").Where("user_id=? AND article_id=?", UserID, articles[i].ArticleID).First(&collect).Error; err != nil {
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
		if err := c.DB.Table("likes").Where("user_id=? AND article_id=?", UserID, articles[i].ArticleID).First(&like).Error; err != nil {
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
		c.DB.Table("comments").Where("article_id=?", articles[i].ArticleID).Order("created_at desc").Find(&comment_db)

		var comments_show = make([]model.ArticleCommentReturn, len(comment_db))
		for j := 0; j < len(comment_db); j++ {
			comments_show[j].UserId = comment_db[j].UserID

			Comment_User_ID := comment_db[j].UserID

			var comment_user model.User
			c.DB.Table("users").Select("*").Where("user_id=?", Comment_User_ID).First(&comment_user)

			comments_show[j].UserName = comment_user.Username
			comments_show[j].Message = comment_db[j].Message
			comments_show[j].CommentId = comment_db[j].CommentID
			comments_show[j].CreatedAt = comment_db[j].CreatedAt
		}

		var tag_db []model.Tag
		c.DB.Table("tags").Where("article_id=?", articles[i].ArticleID).Find(&tag_db)
		var tag_show = make([]model.TagReturn, len(tag_db))
		for k := 0; k < len(tag_show); k++ {
			tag_show[k].Tag = tag_db[k].Tag
		}

		articles[i].Comment = comments_show
		articles[i].Tag = tag_show

		var image_db []model.Image
		c.DB.Table("images").Where("article_id=?", articles[i].ArticleID).Find(&image_db)
		var image_show = make([]model.Image_return, len(image_db))
		for m := 0; m < len(image_db); m++ {
			image_show[m].ImageID = image_db[m].ImageID
			image_show[m].Image = image_db[m].Image
		}

		articles[i].Image = image_show
		articles[i].Avatar = user.Avatar
	}

	//前端渲染需要知道分页的总条数
	var total int
	c.DB.Model(model.Article{}).Count(&total)

	response.Success(ctx, gin.H{"data": articles, "total": total}, "获取所有文章成功!")

}

func (c ArticleController) ShowByType(ctx *gin.Context) {
	article_type := ctx.Query("type")
	UserID, _ := strconv.Atoi(ctx.Query("userid"))

	//获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "3"))

	//开始分页
	var articles_db []model.Article
	c.DB.Table("articles").Where("type=?", article_type).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles_db)

	var articles = make([]model.ArticlesShow, len(articles_db))

	for i := 0; i < len(articles_db); i++ {
		create_user_id := articles_db[i].UserID

		var user model.User
		c.DB.Table("users").Select("*").Where("user_id=?", create_user_id).First(&user)

		articles[i].ArticleID = articles_db[i].ArticleID
		articles[i].UserID = create_user_id
		articles[i].Title = articles_db[i].Title
		articles[i].UserName = user.Username
		articles[i].Content = articles_db[i].Content
		articles[i].CreatedAt = articles_db[i].CreatedAt.Format("2006-01-02 15:04:05")
		articles[i].Type = articles_db[i].Type
		articles[i].Favorites = articles_db[i].Favorites
		articles[i].Likes = articles_db[i].Likes

		var collect model.Collect

		if err := c.DB.Table("collects").Where("user_id=? AND article_id=?", UserID, articles[i].ArticleID).First(&collect).Error; err != nil {
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
		if err := c.DB.Table("likes").Where("user_id=? AND article_id=?", UserID, articles[i].ArticleID).First(&like).Error; err != nil {
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
		c.DB.Table("comments").Where("article_id=?", articles[i].ArticleID).Order("created_at desc").Find(&comment_db)

		var comments_show = make([]model.ArticleCommentReturn, len(comment_db))
		for j := 0; j < len(comment_db); j++ {
			comments_show[j].UserId = comment_db[j].UserID

			Comment_User_ID := comment_db[j].UserID

			var comment_user model.User
			c.DB.Table("users").Select("*").Where("user_id=?", Comment_User_ID).First(&comment_user)

			comments_show[j].UserName = comment_user.Username
			comments_show[j].Message = comment_db[j].Message
			comments_show[j].CommentId = comment_db[j].CommentID
			comments_show[j].CreatedAt = comment_db[j].CreatedAt
		}

		var tag_db []model.Tag
		c.DB.Table("tags").Where("article_id=?", articles[i].ArticleID).Find(&tag_db)
		var tag_show = make([]model.TagReturn, len(tag_db))
		for k := 0; k < len(tag_show); k++ {
			tag_show[k].Tag = tag_db[k].Tag
		}

		articles[i].Comment = comments_show
		articles[i].Tag = tag_show

		var image_db []model.Image
		c.DB.Table("images").Where("article_id=?", articles[i].ArticleID).Find(&image_db)
		var image_show = make([]model.Image_return, len(image_db))
		for m := 0; m < len(image_db); m++ {
			image_show[m].ImageID = image_db[m].ImageID
			image_show[m].Image = image_db[m].Image
		}

		articles[i].Image = image_show
		articles[i].Avatar = user.Avatar

	}

	response.Success(ctx, gin.H{"data": articles}, "")
}

func (c ArticleController) ShowOtherUserArticle(ctx *gin.Context) {
	UserID := ctx.Query("userid")
	MyID := ctx.Query("myid")

	//获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "3"))

	//开始分页
	var articles_db []model.Article
	c.DB.Table("articles").Where("user_id=?", UserID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles_db)

	var articles = make([]model.ArticlesShow, len(articles_db))

	for i := 0; i < len(articles_db); i++ {
		var user model.User
		c.DB.Table("users").Select("*").Where("user_id=?", UserID).First(&user)

		articles[i].ArticleID = articles_db[i].ArticleID
		articles[i].Title = articles_db[i].Title
		articles[i].UserName = user.Username
		articles[i].Content = articles_db[i].Content
		articles[i].CreatedAt = articles_db[i].CreatedAt.Format("2006-01-02 15:04:05")
		articles[i].Type = articles_db[i].Type
		articles[i].Favorites = articles_db[i].Favorites
		articles[i].Likes = articles_db[i].Likes

		var collect model.Collect

		if err := c.DB.Table("collects").Where("user_id=? AND article_id=?", MyID, articles[i].ArticleID).First(&collect).Error; err != nil {
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
		if err := c.DB.Table("likes").Where("user_id=? AND article_id=?", MyID, articles[i].ArticleID).First(&like).Error; err != nil {
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
		c.DB.Table("comments").Where("article_id=?", articles[i].ArticleID).Order("created_at desc").Find(&comment_db)

		var comments_show = make([]model.ArticleCommentReturn, len(comment_db))
		for j := 0; j < len(comment_db); j++ {
			comments_show[j].UserId = comment_db[j].UserID

			Comment_User_ID := comment_db[j].UserID

			var comment_user model.User
			c.DB.Table("users").Select("*").Where("user_id=?", Comment_User_ID).First(&comment_user)

			comments_show[j].UserName = comment_user.Username
			comments_show[j].Message = comment_db[j].Message
			comments_show[j].CommentId = comment_db[j].CommentID
			comments_show[j].CreatedAt = comment_db[j].CreatedAt
		}

		var tag_db []model.Tag
		c.DB.Table("tags").Where("article_id=?", articles[i].ArticleID).Find(&tag_db)
		var tag_show = make([]model.TagReturn, len(tag_db))
		for k := 0; k < len(tag_show); k++ {
			tag_show[k].Tag = tag_db[k].Tag
		}

		articles[i].Comment = comments_show
		articles[i].Tag = tag_show

		var image_db []model.Image
		c.DB.Table("images").Where("article_id=?", articles[i].ArticleID).Find(&image_db)
		var image_show = make([]model.Image_return, len(image_db))
		for m := 0; m < len(image_db); m++ {
			image_show[m].ImageID = image_db[m].ImageID
			image_show[m].Image = image_db[m].Image
		}

		articles[i].Image = image_show
		articles[i].Avatar = user.Avatar

	}

	response.Success(ctx, gin.H{"data": articles}, "获取其它用户的所有文章成功!")
}

func GetPersonalArticles(ctx *gin.Context) { //获取个人发布的所有文章
	var user, temp_user model.User
	var article []model.Article

	userid, er := strconv.Atoi(ctx.Query("userid"))
	user.UserID = uint(userid)
	if er != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
		return
	}
	if err := common.DB.Where("user_id = ?", user.UserID).Order("created_at desc").Find(&article).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400})
		return
	}
	for i := range article {
		if err := common.DB.First(&model.Collect{UserID: article[i].UserID, ArticleID: article[i].ArticleID}).Error; err == nil {
			article[i].IsFavored = true
		}
		if err := common.DB.First(&model.Like{UserID: article[i].UserID, ArticleID: article[i].ArticleID}).Error; err == nil {
			article[i].IsLiked = true
		}
		//取tag
		common.DB.Where("article_id = ?", article[i].ArticleID).Find(&article[i].Tags)
		article[i].TagReturns = make([]model.TagReturn, len(article[i].Tags))
		for j := range article[i].Tags {
			article[i].TagReturns[j].Tag = article[i].Tags[j].Tag
		}
		//取image
		common.DB.Where("article_id = ?", article[i].ArticleID).Find(&article[i].Images)
		article[i].ImageReturn = make([]model.Image_return, len(article[i].Images))
		for j := range article[i].Images {
			article[i].ImageReturn[j].Image = article[i].Images[j].Image
			article[i].ImageReturn[j].ImageID = article[i].Images[j].ImageID
		}
		//取avatar
		if err := common.DB.Where("user_id = ?", article[i].UserID).First(&temp_user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400})
			return
		}
		article[i].Avatar = temp_user.Avatar
		article[i].AuthorName = temp_user.Username
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": article})
}
