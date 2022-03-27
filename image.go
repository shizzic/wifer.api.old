package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
			c.SaveUploadedFile(file, path+"/avatar.webp")
			ConvertResize(path + "/avatar.webp")
		} else {
			files, _ := ioutil.ReadDir(path + "/" + dir)
			root := path + "/" + dir + "/" + fmt.Sprint(len(files)+1) + ".webp"
			c.SaveUploadedFile(file, root)
			ConvertResize(root)
		}
	}
}

func ConvertResize(dir string) {
	buffer, _ := bimg.Read(dir)
	newImage, _ := bimg.NewImage(buffer).Resize(750, 1000)
	new, _ := bimg.NewImage(newImage).Convert(bimg.WEBP)
	bimg.Write(dir, new)
}

func DeleteImage(isAvatar, dir, number string, c gin.Context) string {
	id, _ := c.Cookie("id")

	if isAvatar == "1" {
		files, _ := os.ReadDir("/var/www/html/" + id + "/public")
		os.Remove("/var/www/html/" + id + "/avatar.webp")

		if len(files) > 0 {
			os.Rename("/var/www/html/"+id+"/public/"+files[0].Name(), "/var/www/html/"+id+"/avatar.webp")
			return "replace with another"
		}

		return "deleted avatar without replace"
	} else {
		os.Remove("/var/www/html/" + id + "/" + dir + "/" + number + ".webp")
		return "simple delete"
	}
}

func ChangeImageDir(isAvatar, dir, number string, c gin.Context) string {
	id, _ := c.Cookie("id")

	if isAvatar == "1" {
		private, _ := os.ReadDir("/var/www/html/" + id + "/private")
		os.Rename("/var/www/html/"+id+"/avatar.webp", "/var/www/html/"+id+"/private/"+fmt.Sprint(len(private)+1)+".webp")
		public, _ := os.ReadDir("/var/www/html/" + id + "/public")

		if len(public) > 0 {
			os.Rename("/var/www/html/"+id+"/public/"+public[0].Name(), "/var/www/html/"+id+"/avatar.webp")
			list, _ := os.ReadDir("/var/www/html/" + id + "/public")

			for _, title := range list {
				return title.Name()
			}
		}

		return "deleted avatar without replace"
	}

	return "error"
}
