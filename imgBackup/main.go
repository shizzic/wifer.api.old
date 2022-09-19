package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/kothar/go-backblaze"
	"github.com/mholt/archiver/v4"
)

func main() {
	ZipImages()
}

// Make comprasion of database to gzip file
func ZipImages() {
	files, _ := archiver.FilesFromDisk(nil, map[string]string{
		"/var/www/images": "Images",
	})
	out, _ := os.Create("/var/www/default/site/images.tar.gz")
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err := format.Archive(context.Background(), out, files)
	if err == nil {
		ToBackblaze()
	}
}

/*
	1. Open connection to backblaze
	2. Getting fielId that was uploaded before to delete him
	3. Upload the new file
	4. After upload i save the fileId of uploaded file in txt document to delete him later
*/
func ToBackblaze() {
	b2, _ := backblaze.NewB2(backblaze.Credentials{
		AccountID:      "69119b753b60",
		ApplicationKey: "00496e3b8b3f04f576df1b96fe9c5c9136ac28e711",
	})
	bucket, _ := b2.Bucket("my-wifer")

	_, err := os.Stat("/var/www/default/site/ImagesID.txt")
	if err == nil {
		content, _ := ioutil.ReadFile("/var/www/default/site/ImagesID.txt")
		fileId := string(content)
		bucket.DeleteFileVersion("images.tar.gz", fileId)
	}

	reader, _ := os.Open("/var/www/default/site/images.tar.gz")
	metadata := make(map[string]string)
	res, _ := bucket.UploadFile("images.tar.gz", metadata, reader)

	os.Remove("/var/www/default/site/ImagesID.txt")
	file, _ := os.Create("/var/www/default/site/ImagesID.txt")
	file.WriteString(res.ID)

	os.RemoveAll("/var/www/default/site/images.tar.gz")
}
