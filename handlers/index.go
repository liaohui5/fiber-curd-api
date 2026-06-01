package handlers

import (
	"fiber_curd_api/tools/fmtres"

	"github.com/gofiber/fiber/v3"
)

// PaginationQuery 分页参数(公共参数多个接口可能用到) ?page=1&limit=10
type PaginationQuery struct {
	Page  int `query:"page,default:1"   validate:"min=1"`
	Limit int `query:"limit,default:10" validate:"min=10,max=50"`
}

// PaginationResult 分页响应结果(公共响应需要分页返回数据的接口都会用到)
type PaginationResult struct {
	Items any   `json:"items"`
	Count int64 `json:"count"`
}

func NewPaginationResult(items any, count int64) PaginationResult {
	return PaginationResult{
		Items: items,
		Count: count,
	}
}

// CheckHealth 检查服务健康状态
func CheckHealth(c fiber.Ctx) error {
	return c.JSON(fmtres.OKWithMsg("server is running"))
}
