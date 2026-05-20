package fmtres

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OK(t *testing.T) {
	res := OK()

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "success", res.Message)
	assert.Equal(t, nil, res.Results)
}

func Test_OKWithMsg(t *testing.T) {
	res := OKWithMsg("操作成功")

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "操作成功", res.Message)
	assert.Nil(t, res.Results)
}

func Test_OKWithMsg_EmptyMessage(t *testing.T) {
	res := OKWithMsg("")

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, res.Results)
}

func Test_OKWithResults(t *testing.T) {
	data := map[string]any{"id": 1, "name": "test"}
	res := OKWithResults(data)

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "success", res.Message)
	assert.Equal(t, data, res.Results)
}

func Test_OKWithResults_Nil(t *testing.T) {
	res := OKWithResults(nil)

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "success", res.Message)
	assert.Nil(t, res.Results)
}

func Test_OKWithResults_String(t *testing.T) {
	res := OKWithResults("hello")

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "success", res.Message)
	assert.Equal(t, "hello", res.Results)
}

func Test_OKWithResults_Slice(t *testing.T) {
	items := []string{"a", "b", "c"}
	res := OKWithResults(items)

	assert.Equal(t, true, res.Success)
	assert.Equal(t, "success", res.Message)
	assert.Equal(t, items, res.Results)
}

func Test_Error(t *testing.T) {
	err := errors.New("something went wrong")
	res := Error(err)

	assert.Equal(t, false, res.Success)
	assert.Equal(t, "error", res.Message)
	assert.Equal(t, "something went wrong", res.Results)
}

func Test_Error_NilError(t *testing.T) {
	// Error(nil) will panic because it tries to call err.Error() on nil.
	// This is expected — callers should not pass nil to Error().
	assert.Panics(t, func() {
		Error(nil)
	})
}

func Test_ErrorStr(t *testing.T) {
	res := ErrorStr("参数错误")

	assert.Equal(t, false, res.Success)
	assert.Equal(t, "error", res.Message)
	assert.Equal(t, "参数错误", res.Results)
}

func Test_ErrorStr_Empty(t *testing.T) {
	res := ErrorStr("")

	assert.Equal(t, false, res.Success)
	assert.Equal(t, "error", res.Message)
	assert.Equal(t, "", res.Results)
}

func Test_ErrorFmt(t *testing.T) {
	err := errors.New("connection refused")
	res := ErrorFmt("数据库连接失败", err)

	assert.Equal(t, false, res.Success)
	assert.Equal(t, "error", res.Message)
	assert.EqualError(t, res.Results.(error), "数据库连接失败:connection refused")
}

func Test_ErrorFmt_NilError(t *testing.T) {
	res := ErrorFmt("某个错误", nil)

	assert.Equal(t, false, res.Success)
	assert.Equal(t, "error", res.Message)
	assert.EqualError(t, res.Results.(error), "某个错误:<nil>")
}

func Test_FormatResponse_AllFieldsComplete(t *testing.T) {
	res := FormatResponse{
		Success: true,
		Message: "自定义消息",
		Results: 42,
	}

	assert.True(t, res.Success)
	assert.Equal(t, "自定义消息", res.Message)
	assert.Equal(t, 42, res.Results)
}
