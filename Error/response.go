package response

import "errors"

type ServerRes struct {
	Status   int         `json:"status" bson:"status" binding:"required"`
	Success  bool        `json:"success" bson:"success" binding:"required"`
	Response interface{} `json:"response" bson:"response" binding:"required"`
	Error    error       `json:"error" bson:"error" binding:"required"`
}

type BindingErrRes struct {
	Status   int    `json:"status" bson:"status" binding:"required"`
	Response string `json:"response" bson:"response" binding:"required"`
	Error    string `json:"error" bson:"error" binding:"required"`
}

var BindingErr = BindingErrRes{
	Status:   405,
	Error:    "Error",
	Response: "Binding Error",
}

type ServerErrRes struct {
	Status   int    `json:"status" bson:"status" binding:"required"`
	Response string `json:"response" bson:"response" binding:"required"`
}

var SeverErr = ServerErrRes{
	Status:   404,
	Response: "Sever Failed",
}

var StringToObjevctIdError = ServerRes{
	Status:   400,
	Success:  false,
	Response: "Invalid ObjectID",
	Error:    errors.New("Invalid ObjectID"),
}
