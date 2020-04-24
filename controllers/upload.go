package controllers

import (
	"female/lib/result"
	"female/models/page"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func Upload(ctx *gin.Context) {
	response := new(result.JsonResponseString)
	response.ErrNum = result.RESULT_SUCCESS
	response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
	file, handler, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Fatal("ctx.Request.FormFile error,err=" + err.Error())
		response.ErrNum = result.RESULT_FILE_UPLOAD_ERROR
		response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
		response.EchoResult(ctx)
		return
	}
	fileName := handler.Filename
	loc, err := os.Create("./public/" + fileName)
	if err != nil {
		log.Fatal("os.Create error,err=" + err.Error())
		response.ErrNum = result.RESULT_FILE_UPLOAD_ERROR
		response.ErrMsg = result.ErrInfos[response.ErrNum].ErrMsg
		response.EchoResult(ctx)
		return
	}
	io.Copy(loc, file)

	page.Upload(ctx, fileName, response)
	response.EchoResult(ctx)
	return
}
