package handlers

import (
	"fiber_curd_api/app"
	"fiber_curd_api/models"
	"fiber_curd_api/tools/fmtres"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// @Tags    task
// @Summary 创建任务
// @Accept  json
// @Produce json
// @Param   data body models.Task true "任务信息"
// @Success 200 {object} fmtres.FormatResponse{results=models.Task}
// @Router /api/tasks [post]
func CreateTasks(c fiber.Ctx) error {
	task := new(models.Task)

	if err := c.Bind().JSON(task); err != nil {
		return c.JSON(fmtres.Error(err))
	}

	result := app.ConnectDB().Create(task)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	//// 添加任务到异步任务队列
	// asyncTask := asynq.NewTask("default", fmt.Appendf(nil, "%d", task.ID))
	// info, err := app.InitAsynqClient().Enqueue(asyncTask)
	// if err != nil {
	// 	return c.JSON(fmtres.Error(err))
	// }
	// return c.JSON(fmtres.OKWithResults(info))

	return c.JSON(fmtres.OK())
}

// @Tags    task
// @Summary 搜索任务列表
// @Accept  json
// @Produce json
// @Param   page query string false "当前页数(默认:1)"
// @Param   limit query string false "每页多少条数据(默认:10)"
// @Param   search query string false "搜索参数(默认:{},json字符串{f:title,v:Golang}搜索 title 字段带有 Golang 的数据)"
// @Success 200 {object} fmtres.FormatResponse{results=[]models.Task}
// @Router /api/tasks [get]
func SearchTasks(c fiber.Ctx) error {
	pagination := new(PaginationQuery)

	// TODO: 支持搜索参数
	if err := c.Bind().Query(pagination); err != nil {
		return c.JSON(fmtres.Error(err))
	}

	offset := (pagination.Page - 1) * pagination.Limit

	var items []models.Task
	var count int64

	db := app.ConnectDB().Model(&models.Task{})
	db.Count(&count)
	db.Omit("MetaData").Find(&items).Offset(offset).Limit(pagination.Limit)

	return c.JSON(fmtres.OKWithResults(NewPaginationResult(items, count)))
}

// @Tags    task
// @Summary 获取任务详情
// @Accept  json
// @Produce json
// @Param   id query string true "任务ID:1"
// @Success 200 {object} fmtres.FormatResponse{results=models.Task}
// @Router /api/tasks/:id [get]
func TaskDetail(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}
	var task models.Task
	app.ConnectDB().First(&task, id)
	return c.JSON(fmtres.OKWithResults(task))
}

// @Tags    task
// @Summary 删除任务
// @Accept  json
// @Produce json
// @Param   id query string true "任务ID"
// @Success 200 {object} fmtres.FormatResponse{}
// @Router /api/tasks/:id [delete]
func DeleteTasks(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(fmtres.Error(err))
	}

	var task models.Task
	app.ConnectDB().First(&task, idInt)
	if task.Status != models.TaskStatusFailed {
		return c.JSON(fmtres.ErrorStr("task status is not failed"))
	}

	result := app.ConnectDB().Delete(&task)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	return c.JSON(fmtres.OKWithResults(result.RowsAffected))
}

// @Tags    task
// @Summary 重试失败的任务
// @Accept  json
// @Produce json
// @Param   id query string true "任务ID"
// @Success 200 {object} fmtres.FormatResponse{}
// @Router /api/tasks/:id/retry [get]
func RetryTasks(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(fmtres.ErrorStr("not found article id"))
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(fmtres.Error(err))
	}

	var task models.Task
	app.ConnectDB().First(&task, idInt)
	if task.Status != models.TaskStatusFailed {
		return c.JSON(fmtres.ErrorStr("task status is not failed"))
	}

	task.Status = models.TaskStatusPending
	app.ConnectDB().Save(&task)

	return c.JSON(fmtres.OK())

	// 重新入列
	// asyncTask := asynq.NewTask("default", fmt.Appendf(nil, "%d", task.ID))
	// info, err := app.InitAsynqClient().Enqueue(asyncTask)
	// if err != nil {
	// 	return c.JSON(fmtres.Error(err))
	// }
	// return c.JSON(fmtres.OKWithResults(info))
}
