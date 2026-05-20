package handlers

import (
	"fiber_curd_api/app"
	"fiber_curd_api/models"
	"fiber_curd_api/tools/fmtres"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// @Tags    article
// @Summary 搜索文章列表
// @Accept  json
// @Produce json
// @Param   page query string false "当前页数(默认:1)"
// @Param   limit query string false "每页多少条数据(默认:10)"
// @Param   search query string false "搜索参数(默认:{},json字符串{f:title,v:Golang}搜索 title 字段带有 Golang 的数据)"
// @Success 200 {object} fmtres.FormatResponse{results=[]models.Article}
// @Router /api/articles [get]
func SearchArticles(c fiber.Ctx) error {
	pagination := new(PaginationQuery)

	// TODO: 支持搜索参数
	if err := c.Bind().Query(pagination); err != nil {
		return c.JSON(fmtres.ErrorFmt("failed to parse query params", err))
	}

	offset := (pagination.Page - 1) * pagination.Limit

	var items []models.Article
	var count int64

	db := app.ConnectDB().Model(&models.Article{})
	db.Count(&count)
	db.Preload("Author").Omit("contents").Find(&items).Offset(offset).Limit(pagination.Limit)

	return c.JSON(fmtres.OKWithResults(NewPaginationResult(items, count)))
}

// @Tags    article
// @Summary 获取文章详情
// @Accept  json
// @Produce json
// @Param   id query string true "文章ID:1"
// @Success 200 {object} fmtres.FormatResponse{results=models.Article}
// @Router /api/articles/:id [get]
func ArticleDetail(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}
	var article models.Article
	app.ConnectDB().First(&article, id)
	return c.JSON(fmtres.OKWithResults(article))
}

// @Tags    article
// @Summary 创建文章
// @Accept  json
// @Produce json
// @Param   data body models.Article true "文章信息"
// @Success 200 {object} fmtres.FormatResponse{results=models.Article}
// @Router /api/articles [post]
func CreateArtilce(c fiber.Ctx) error {
	article := new(models.Article)

	if err := c.Bind().JSON(article); err != nil {
		return c.JSON(fmtres.Error(err))
	}

	result := app.ConnectDB().Create(article)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	return c.JSON(fmtres.OK())
}

// @Tags    article
// @Summary 修改文章信息
// @Accept  json
// @Produce json
// @Param   data body models.Article true "文章信息"
// @Success 200 {object} fmtres.FormatResponse{results=models.Article}
// @Router /api/articles/:id [patch]
func UpdateArticle(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}

	article := new(models.Article)

	if err := c.Bind().JSON(article); err != nil {
		return c.JSON(fmtres.Error(err))
	}

	result := app.ConnectDB().Where("id = ?", id).Updates(article)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	return c.JSON(fmtres.OKWithResults(article))
}

// @Tags    article
// @Summary 删除文章
// @Accept  json
// @Produce json
// @Param   id query string true "文章ID"
// @Success 200 {object} fmtres.FormatResponse{}
// @Router /api/articles/:id [delete]
func DeleteArticle(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(fmtres.Error(err))
	}

	result := app.ConnectDB().Delete(&models.Article{
		ID: uint(idInt),
	})
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	return c.JSON(fmtres.OKWithResults(result.RowsAffected))
}
