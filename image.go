package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadImage(dir string, c gin.Context) {
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
			UpdateDataBaseImages(id, dir, true, true, 0)
		} else {
			files, _ := ioutil.ReadDir(path + "/" + dir)
			root := path + "/" + dir + "/" + fmt.Sprint(len(files)+1) + ".webp"
			c.SaveUploadedFile(file, root)
			ConvertResize(root)
			UpdateDataBaseImages(id, dir, false, true, len(files)+1)
		}
	}
}

func DeleteImage(isAvatar, dir, number string, c gin.Context) {
	id, _ := c.Cookie("id")

	if isAvatar == "1" {
		files, _ := os.ReadDir("/var/www/html/" + id + "/public")
		os.Remove("/var/www/html/" + id + "/avatar.webp")

		if len(files) > 0 {
			os.Rename("/var/www/html/"+id+"/public/"+files[0].Name(), "/var/www/html/"+id+"/avatar.webp")
			RenameImages("/var/www/html/"+id+"/", "public", id, false, true)
		} else {
			RenameImages("/var/www/html/"+id+"/", "public", id, true, false)
		}
	} else {
		os.Remove("/var/www/html/" + id + "/" + dir + "/" + number + ".webp")
		RenameImages("/var/www/html/"+id+"/", dir, id, false, true)
	}
}

func ChangeImageDir(isAvatar, dir, new, number string, c gin.Context) {
	id, _ := c.Cookie("id")

	if isAvatar == "1" {
		private, _ := os.ReadDir("/var/www/html/" + id + "/private")
		os.Rename("/var/www/html/"+id+"/avatar.webp", "/var/www/html/"+id+"/private/"+fmt.Sprint(len(private)+1)+".webp")
		public, _ := os.ReadDir("/var/www/html/" + id + "/public")

		if len(public) > 0 {
			os.Rename("/var/www/html/"+id+"/public/"+public[0].Name(), "/var/www/html/"+id+"/avatar.webp")
			RenameImages("/var/www/html/"+id+"/", "public", id, true, true)
		} else {
			RenameImages("/var/www/html/"+id+"/", "private", id, true, false)
		}
	} else {
		// сделать аватаркой приват фотку, если была попытка сделать ее публичной без аватарки
		list, _ := os.ReadDir("/var/www/html/" + id + "/" + new)
		// _, err := os.Stat("/var/www/html/" + id + "/avatar.webp")

		os.Rename("/var/www/html/"+id+"/"+dir+"/"+number+".webp", "/var/www/html/"+id+"/"+new+"/"+fmt.Sprint(len(list)+1)+".webp")
		RenameImages("/var/www/html/"+id+"/", dir, id, false, true)

		UpdateDataBaseImages(id, new, false, true, len(list)+1)
	}
}

func RenameImages(path, dir, id string, isAvatar, avatar bool) {
	os.MkdirAll(path+"new", os.ModePerm)
	list, _ := os.ReadDir(path + dir)

	for key, image := range list {
		os.Rename(path+dir+"/"+image.Name(), path+"new/"+fmt.Sprint(key+1)+".webp")
	}

	os.Remove(path + dir)
	os.Rename(path+"new", path+dir)
	UpdateDataBaseImages(id, dir, isAvatar, avatar, len(list))
}

func UpdateDataBaseImages(id, dir string, isAvatar, avatar bool, quantity int) {
	if isAvatar {
		users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "avatar", Value: avatar}}}})
	} else {
		users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: bson.D{{Key: dir, Value: quantity}}}})
	}
}

// Resize each uploaded image to 1 size format and convert each image to WEBP
func ConvertResize(dir string) {
	buffer, _ := bimg.Read(dir)
	newImage, _ := bimg.NewImage(buffer).Resize(750, 1000)
	new, _ := bimg.NewImage(newImage).Convert(bimg.WEBP)
	bimg.Write(dir, new)
}
