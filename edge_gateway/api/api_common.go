package api

var DefaultParentId = int64(-1)

func Resp(data interface{}) map[string]interface{} {
	resp := map[string]interface{}{}
	resp["data"] = data
	resp["code"] = "200"
	return resp
}

func RespError(msg interface{}) map[string]interface{} {
	resp := map[string]interface{}{}
	resp["code"] = "500"
	resp["msg"] = msg
	return resp
}

type DeleteReq struct {
	ID int64 `json:"id" binding:"required"`
}
