// Code generated by hertz generator.

package handler

import (
	"context"
	"crypto/md5"
	"fmt"
	"hzminioupload/biz/pkg/global"
	"hzminioupload/biz/uploadresponse"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// upload to minio .
func Upload(ctx context.Context, c *app.RequestContext) {
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(400, uploadresponse.ErrorResponse{
			Code:    400,
			Message: "file is empty",
			Reason:  "file is empty",
		})
		return
	}
	dm5filename := md5.Sum([]byte(file.Filename)) //转成加密编码
	filename := fmt.Sprintf("%x", dm5filename)
	//get file ext without using utils
	ext := file.Filename[strings.LastIndex(file.Filename, "."):]

	newfilename := filename + ext
	// get now time  YYYY/MM/DD string
	now := time.Now()
	nowstr := now.Format("2006/01/02")

	savepath := nowstr + "/" + newfilename

	//get upload file temp path
	tempfile, err := file.Open()
	if err != nil {
		c.JSON(400, uploadresponse.ErrorResponse{
			Code:    400,
			Message: "open file error",
			Reason:  err.Error(),
		})
		return
	}
	defer tempfile.Close()

	// info, err := global.S_MinioClient.Client.PutObject(ctx,
	// 	global.S_MinioClient.BucketName,
	// 	newfilename,
	// 	tempfile,
	// 	file.Size,
	// 	minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	// fmt.Println("upload object:", info)
	// fmt.Println("upload object error:", err)
	// return
	// save file to minio
	fileurl, err := global.S_MinioClient.UpLoadFile(ctx, savepath, file.Header.Get("content-type"), tempfile)
	if err != nil {
		c.JSON(400, uploadresponse.ErrorResponse{
			Code:    400,
			Message: "upload file error",
			Reason:  err.Error(),
		})
		return
	}

	c.JSON(200, uploadresponse.UploadResponse{
		Code:    200,
		Message: "success",
		Reason:  "success",
		Data: uploadresponse.File{
			FileName: file.Filename,
			FileSize: file.Size,
			FileType: file.Header.Get("Content-Type"),
			FileUrl:  fileurl,
		},
	})

}

// checkfile md5
func Checkfile(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, utils.H{
		"message": global.S_CONFIG.GetString("filedriver.storepath"),
	})
}
