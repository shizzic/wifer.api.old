package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nickalie/go-webpbin"
	"go.mongodb.org/mongo-driver/bson"
)

type File struct {
	ID        string
	EntryPath string
	FullPath  string
	Dir       string `form:"dir"`
	NewDir    string `form:"newDir"`
	Name      string `form:"name"`
	Extension string `form:"extension"`
	IsAvatar  bool   `form:"isAvatar"`
}

// Just upload new image. If avatar isn't exist, then create avatar
func UploadImage(data File, c *gin.Context) error {
	os.MkdirAll(data.EntryPath+"/public", os.ModePerm)
	os.MkdirAll(data.EntryPath+"/private", os.ModePerm)
	file, _ := c.FormFile("file")

	if file.Size < 20000001 {
		public_images, _ := os.ReadDir(data.EntryPath + "/public")
		private_images, _ := os.ReadDir(data.EntryPath + "/private")
		public_images_count := len(public_images)
		private_images_count := len(private_images)
		images_count := public_images_count + private_images_count
		f, err := os.Open(data.EntryPath + "/avatar.webp")
		f.Close()
		if err == nil {
			images_count++
		}

		if images_count < 20 {
			if data.Dir == "public" && err != nil {
				data.FullPath = data.EntryPath + "/avatar.webp"
				c.SaveUploadedFile(file, data.FullPath)
				ConvertToWebP(data.FullPath)
				RecountImages(data)
			} else {
				var number int
				if data.Dir == "public" {
					number = public_images_count + 1
				} else {
					number = private_images_count + 1
				}
				data.Name = strconv.Itoa(number) + ".webp"
				data.FullPath = data.EntryPath + "/" + data.Dir + "/" + data.Name

				c.SaveUploadedFile(file, data.FullPath)
				ConvertToWebP(data.FullPath)
				RecountImages(data)
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
func DeleteImage(data File) error {
	if data.IsAvatar {
		public_images, _ := os.ReadDir(data.EntryPath + "/public")
		os.Remove(data.EntryPath + "/avatar.webp")

		if len(public_images) > 0 {
			os.Rename(data.EntryPath+"/public/"+public_images[0].Name(), data.EntryPath+"/avatar.webp")
			RenameImages(data, "public")
		} else {
			RecountImages(data)
		}
	} else {
		os.Remove(data.EntryPath + "/" + data.Dir + "/" + data.Name + ".webp")
		RenameImages(data, data.Dir)
	}

	return nil
}

// Change public or private dir for image. Avatar must exist always
func ChangeImageDir(data File) {
	if data.IsAvatar {
		private, _ := os.ReadDir(data.EntryPath + "/private")
		os.Rename(data.EntryPath+"/avatar.webp", data.EntryPath+"/private/"+strconv.Itoa(len(private)+1)+".webp")
		public, _ := os.ReadDir(data.EntryPath + "/public")

		if len(public) > 0 {
			os.Rename(data.EntryPath+"/public/"+public[0].Name(), data.EntryPath+"/avatar.webp")
			RenameImages(data, "public")
		} else {
			RenameImages(data, "private")
		}
	} else {
		f, err := os.Open(data.EntryPath + "/avatar.webp")
		f.Close()

		// Make avatar from private image if avatar isn't exist
		if err != nil && data.NewDir == "public" {
			os.Rename(data.EntryPath+"/"+data.Dir+"/"+data.Name+".webp", data.EntryPath+"/avatar.webp")
			RenameImages(data, data.Dir)
		} else {
			files, _ := os.ReadDir(data.EntryPath + "/" + data.NewDir)
			os.Rename(data.EntryPath+"/"+data.Dir+"/"+data.Name+".webp", data.EntryPath+"/"+data.NewDir+"/"+strconv.Itoa(len(files)+1)+".webp")
			RenameImages(data, data.Dir)
		}
	}
}

// // Replaced avatar alwayc contain in public dir
func ReplaceAvatar(data File) {
	os.Rename(data.EntryPath+"/"+data.Dir+"/"+data.Name+".webp", data.EntryPath+"/new_avatar.webp")

	if data.Dir == "private" {
		files, _ := os.ReadDir(data.EntryPath + "/public")
		os.Rename(data.EntryPath+"/avatar.webp", data.EntryPath+"/public/"+strconv.Itoa(len(files)+1)+".webp")
	} else {
		os.Rename(data.EntryPath+"/avatar.webp", data.EntryPath+"/public/"+data.Name+".webp")
	}

	os.Rename(data.EntryPath+"/new_avatar.webp", data.EntryPath+"/avatar.webp")

	if data.Dir == "private" {
		RenameImages(data, "private")
		RenameImages(data, "public")
	}
}

// Make straight line of images from 1 to count of all images in dir
func RenameImages(data File, target_dir string) {
	os.MkdirAll(data.EntryPath+"/new", os.ModePerm)
	files, _ := os.ReadDir(data.EntryPath + "/" + target_dir)
	for index, image := range files {
		os.Rename(data.EntryPath+"/"+target_dir+"/"+image.Name(), data.EntryPath+"/new/"+strconv.Itoa(index+1)+".webp")
	}
	os.Remove(data.EntryPath + "/" + target_dir)
	os.Rename(data.EntryPath+"/new", data.EntryPath+"/"+target_dir)
	RecountImages(data)
}

// Update info about images in database
func RecountImages(data File) {
	public_images, _ := os.ReadDir(data.EntryPath + "/public")
	private_images, _ := os.ReadDir(data.EntryPath + "/private")
	public_images_count := len(public_images)
	private_images_count := len(private_images)
	images_count := public_images_count + private_images_count
	has_avatar := false
	f, err := os.Open(data.EntryPath + "/avatar.webp")
	f.Close()
	if err == nil {
		has_avatar = true
		images_count++
	}

	user_id, _ := strconv.Atoi(data.ID)
	DB["users"].UpdateOne(ctx, bson.M{"_id": user_id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "avatar", Value: has_avatar}}},
		{Key: "$set", Value: bson.D{{Key: "public", Value: public_images_count}}},
		{Key: "$set", Value: bson.D{{Key: "private", Value: private_images_count}}},
		{Key: "$set", Value: bson.D{{Key: "images", Value: images_count}}},
	})
}

// Convert image to WebP format
func ConvertToWebP(full_path_to_file string) {
	webpbin.NewCWebP().
		Quality(80).
		InputFile(full_path_to_file).
		OutputFile(full_path_to_file).
		Run()
}
