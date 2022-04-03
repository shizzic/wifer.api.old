package main

import (
	"context"
	"os"
	"os/exec"

	"github.com/kothar/go-backblaze"
	"github.com/mholt/archiver/v4"
)

func main() {
	err := exec.Command("mongodump", "mongodb://shizzic:WebDev77@wifer-test.ru:27017/wifer?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false", "-d", "wifer", "-o", "/var/www/default/site").Run()

	if err == nil {
		ZipBackUp()
	}
}

// Make comprasion of database to gzip file
func ZipBackUp() {
	files, _ := archiver.FilesFromDisk(nil, map[string]string{
		"/var/www/default/site/wifer": "Archive",
	})
	out, _ := os.Create("/var/www/default/site/archive.tar.gz")
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err := format.Archive(context.Background(), out, files)
	if err == nil {
		os.RemoveAll("/var/www/default/site/wifer")
		ToBackblaze()
	}
}

func ToBackblaze() {
	b2, _ := backblaze.NewB2(backblaze.Credentials{
		AccountID:      "69119b753b60",
		ApplicationKey: "00496e3b8b3f04f576df1b96fe9c5c9136ac28e711",
	})

	bucket, _ := b2.Bucket("my-wifer")
	reader, _ := os.Open("/var/www/default/site/archive.tar.gz")
	metadata := make(map[string]string)
	bucket.UploadFile("db.tar.gz", metadata, reader)
}
