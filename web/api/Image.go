package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iahfdoa/crawlsForBeauty/web/dao"
	"github.com/iahfdoa/crawlsForBeauty/web/util"
	"golang.org/x/exp/rand"
	"net/http"
	"os"
	"time"
)

func ImageController(c *gin.Context) {
	config, err := util.LoadImagePathsFromConfigFile("config/config.json")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "系统错误",
			"code":    "5000",
		})
		return
	}
	t := c.DefaultQuery("type", fmt.Sprintf("%d", dao.GetConfig().Type))
	imagesPath := config[t]
	if len(imagesPath) == 0 {
		c.JSON(200, gin.H{
			"message": "存储库没有照片",
			"code":    "4000",
		})
		return
	}
	rand.Seed(uint64(time.Now().Unix()))      // 使用当前时间设置随机数种子
	randomIndex := rand.Intn(len(imagesPath)) // 生成随机索引
	imageFilename := imagesPath[randomIndex]
	imageData, err := os.ReadFile(imageFilename)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"code":    "5000",
		})
		return
	}
	// 设置响应头，指定图片的Content-Type
	//c.Header("Content-Type", "image/jpeg")

	// 将图片内容作为响应返回给客户端
	c.Data(http.StatusOK, "image/jpeg", imageData)
	return
}
