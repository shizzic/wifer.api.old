package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

func SaveImage(dir string, c gin.Context) {
	file, _ := c.FormFile("file")

	if file.Size < 1000001 {
		id, _ := c.Cookie("id")
		path := "/var/www/html/" + id
		os.MkdirAll(path+"/public", os.ModePerm)
		os.MkdirAll(path+"/private", os.ModePerm)
		_, avatar := os.Stat(path + "/avatar.webp")

		if dir != "private" && avatar != nil {
			c.SaveUploadedFile(file, path+"/avatar.png")
			Replace(path + "/avatar.")
		} else {
			files, _ := ioutil.ReadDir(path + "/" + dir)
			root := path + "/" + dir + "/" + fmt.Sprint(len(files)+1) + "."
			c.SaveUploadedFile(file, root+"png")
			Replace(root)
		}
	}
}

func Replace(dir string) {
	src, _ := imaging.Open(dir + "png")
	dst := imaging.Resize(src, 750, 1000, imaging.Lanczos)
	imaging.Save(dst, dir+"png")

	buffer, _ := bimg.Read(dir + "png")
	newImage, _ := bimg.NewImage(buffer).Convert(bimg.WEBP)
	bimg.Write(dir+"webp", newImage)

	os.Remove(dir + "png")
}
