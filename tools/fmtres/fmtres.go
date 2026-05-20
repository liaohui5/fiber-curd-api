package fmtres

import "fmt"

///////////////////////////////////////
// NOTE: 格式化响应的 JSON 格式
///////////////////////////////////////

type FormatResponse struct {
	Success bool   `json:"success"` // 操作是否成功
	Message string `json:"message"` // 操作附带消息
	Results any    `json:"results"` // 操作附带数据
}

func OK() FormatResponse {
	return FormatResponse{
		Success: true,
		Message: "success",
		Results: nil,
	}
}

func OKWithMsg(msg string) FormatResponse {
	return FormatResponse{
		Success: true,
		Message: msg,
		Results: nil,
	}
}

func OKWithResults(results any) FormatResponse {
	return FormatResponse{
		Success: true,
		Message: "success",
		Results: results,
	}
}

func Error(err error) FormatResponse {
	return FormatResponse{
		Success: false,
		Message: "error",
		Results: err.Error(),
	}
}

func ErrorStr(message string) FormatResponse {
	return FormatResponse{
		Success: false,
		Message: "error",
		Results: message,
	}
}

func ErrorFmt(message string, err error) FormatResponse {
	return FormatResponse{
		Success: false,
		Message: "error",
		Results: fmt.Errorf("%s:%v", message, err).Error(),
	}
}
