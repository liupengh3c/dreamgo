package page

import (
	"female/lib/result"
	"female/lib/tools"
	"female/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

/**
* [Upload 上传文件分析数据]
* @Author   liupeng17
* @DateTime 2019-05-26 08:42
* @param    [*gin.Context]     ctx      [上下文]
* @param    [string]   		   fileName [文件名称]
* @param    [response]         *result.JsonResponseString  [数据返回]
* @return
 */
func Upload(ctx *gin.Context, fileName string, response *result.JsonResponseString) {
	curPath := tools.GetCurrentDirectory()
	// 获取图片名称
	tmp := strings.Split(fileName, ".")
	log.Notice("log debug")
	fmt.Println(curPath + tmp[0])
	return
}
