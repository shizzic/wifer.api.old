package main

import (
	"context"
	"os"

	"github.com/kothar/go-backblaze"
	"github.com/mholt/archiver/v4"
)

// Make comprasion of database to gzip file
func zip_images() {
	files, _ := archiver.FilesFromDisk(nil, map[string]string{
		config.PATH + "/images": "",
	})
	to, _ := os.Create(config.PATH + "/cron_dump/images.tar.gz")
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err := format.Archive(context.Background(), to, files)
	to.Close()
	if err == nil {
		upload_images_zip_to_backblaze()
	}
}

/*
1. Open connection to backblaze
2. Getting fielId that was uploaded before to delete him
3. Upload the new file
4. After upload i save the fileId of uploaded file in txt document to delete it later
*/
func upload_images_zip_to_backblaze() {
	b2, connection_error := backblaze.NewB2(backblaze.Credentials{
		AccountID:      config.BACKBLAZE_ID,
		ApplicationKey: config.BACKBLAZE_KEY,
	})

	if connection_error == nil {
		b2.CreateBucket(config.PRODUCT_NAME, backblaze.AllPrivate)
		bucket, bucket_error := b2.Bucket(config.PRODUCT_NAME)

		if bucket_error == nil {
			old_file_id, err := os.Open(config.PATH + "/cron_dump/images_dump_id.txt")
			if err == nil {
				content, _ := os.ReadFile(config.PATH + "/cron_dump/images_dump_id.txt")
				file_id := string(content)
				old_file_id.Close()
				bucket.DeleteFileVersion("images.tar.gz", file_id)
			}

			reader, _ := os.Open(config.PATH + "/cron_dump/images.tar.gz")
			metadata := make(map[string]string)
			res, _ := bucket.UploadFile("images.tar.gz", metadata, reader)
			reader.Close()

			os.Remove(config.PATH + "/cron_dump/images_dump_id.txt")
			new_file_id, _ := os.Create(config.PATH + "/cron_dump/images_dump_id.txt")
			new_file_id.WriteString(res.ID)
			new_file_id.Close()
			os.RemoveAll(config.PATH + "/cron_dump/images.tar.gz")
		}
	}
}
