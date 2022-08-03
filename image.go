package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
	"go.mongodb.org/mongo-driver/bson"
)

// Just upload new image. If avatar isn't exist, then create avatar
func UploadImage(dir string, c gin.Context) error {
	file, _ := c.FormFile("file")

	if file.Size < 20000001 {
		isAvatar := 0
		id, _ := c.Cookie("id")
		path := "/var/www/html/" + id
		os.MkdirAll(path+"/public", os.ModePerm)
		os.MkdirAll(path+"/private", os.ModePerm)
		_, avatar := os.Stat(path + "/avatar.webp")
		idInt, _ := strconv.Atoi(id)

		if avatar == nil {
			isAvatar = 1
		}

		private, _ := os.ReadDir(path + "/private")
		public, _ := os.ReadDir(path + "/public")

		all := len(public) + len(private) + isAvatar

		if all < 20 {
			if dir != "private" && avatar != nil {
				c.SaveUploadedFile(file, path+"/avatar.webp")
				ConvertResize(path + "/avatar.webp")
				UpdateDataBaseImages(idInt, path)
			} else {
				files, _ := ioutil.ReadDir(path + "/" + dir)
				root := path + "/" + dir + "/" + fmt.Sprint(len(files)+1) + ".webp"
				c.SaveUploadedFile(file, root)
				ConvertResize(root)
				UpdateDataBaseImages(idInt, path)
			}
		} else {
			return errors.New("max_image")
		}
	} else {
		return errors.New("max_size")
	}

	return nil
}

// Delete image. If target was an avatar, then make new avatar from first public image (if this image exist)
func DeleteImage(isAvatar, dir, number string, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	path := "/var/www/html/" + id

	if isAvatar == "1" {
		files, _ := os.ReadDir(path + "/public")
		os.Remove(path + "/avatar.webp")

		if len(files) > 0 {
			os.Rename(path+"/public/"+files[0].Name(), path+"/avatar.webp")
			RenameImages(path+"/", "public", idInt)
		} else {
			RenameImages(path+"/", "public", idInt)
		}
	} else {
		os.Remove(path + "/" + dir + "/" + number + ".webp")
		RenameImages(path+"/", dir, idInt)
	}
}

// Change public or private dir for image. Avatar must exist always
func ChangeImageDir(isAvatar, dir, new, number string, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	path := "/var/www/html/" + id

	if isAvatar == "1" {
		private, _ := os.ReadDir(path + "/private")
		os.Rename(path+"/avatar.webp", path+"/private/"+fmt.Sprint(len(private)+1)+".webp")
		public, _ := os.ReadDir(path + "/public")

		if len(public) > 0 {
			os.Rename(path+"/public/"+public[0].Name(), path+"/avatar.webp")
			RenameImages(path+"/", "public", idInt)
		} else {
			RenameImages(path+"/", "private", idInt)
		}
	} else {
		list, _ := os.ReadDir(path + "/" + new)
		_, err := os.Stat(path + "/avatar.webp")

		// Make avatar from private image if avatar isn't exist
		if err != nil && new == "public" {
			os.Rename(path+"/"+dir+"/"+number+".webp", path+"/avatar.webp")
			RenameImages(path+"/", dir, idInt)
		} else if err == nil {
			os.Rename(path+"/"+dir+"/"+number+".webp", path+"/"+new+"/"+fmt.Sprint(len(list)+1)+".webp")
			RenameImages(path+"/", dir, idInt)
		}
	}
}

// Make straight line of images from 1 to count of all images in dir
func RenameImages(path, dir string, id int) {
	os.MkdirAll(path+"new", os.ModePerm)
	list, _ := os.ReadDir(path + dir)

	for key, image := range list {
		os.Rename(path+dir+"/"+image.Name(), path+"new/"+fmt.Sprint(key+1)+".webp")
	}

	os.Remove(path + dir)
	os.Rename(path+"new", path+dir)
	UpdateDataBaseImages(id, path)
}

// Replaced avatar alwayc contain in public dir
func ReplaceAvatar(dir, num string, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	path := "/var/www/html/" + id
	os.Rename(path+"/"+dir+"/"+num+".webp", path+"/new.webp")

	if dir == "private" {
		files, _ := os.ReadDir(path + "public")
		os.Rename(path+"/avatar.webp", path+"/public/"+fmt.Sprint(len(files)+1)+".webp")
	} else {
		os.Rename(path+"/avatar.webp", path+"/public/"+num+".webp")
	}

	os.Rename(path+"/new.webp", path+"/avatar.webp")

	if dir == "private" {
		RenameImages(path+"/", "private", idInt)
		RenameImages(path+"/", "public", idInt)
	}
}

// Update info about images in database
func UpdateDataBaseImages(id int, path string) {
	avatar := true
	ava := 1
	_, err := os.Stat(path + "/avatar.webp")

	if err != nil {
		avatar = false
		ava = 0
	}

	public, _ := os.ReadDir(path + "/public")
	private, _ := os.ReadDir(path + "/private")
	pubLen := len(public)
	priLen := len(private)

	users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "avatar", Value: avatar}}},
		{Key: "$set", Value: bson.D{{Key: "public", Value: pubLen}}},
		{Key: "$set", Value: bson.D{{Key: "private", Value: priLen}}},
		{Key: "$set", Value: bson.D{{Key: "images", Value: pubLen + priLen + ava}}},
	})
}

// Resize each uploaded image to 1 size format and convert each image to WEBP
func ConvertResize(dir string) {
	buffer, _ := bimg.Read(dir)
	// newImage, _ := bimg.NewImage(buffer).Resize(750, 1000)
	new, _ := bimg.NewImage(buffer).Convert(bimg.WEBP)
	bimg.Write(dir, new)
}
