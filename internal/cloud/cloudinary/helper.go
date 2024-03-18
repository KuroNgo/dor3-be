package cloudinary

import (
	"clean-architecture/bootstrap"
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"time"
)

var (
	Database *bootstrap.Database
)

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cld, err := cloudinary.NewFromParams(Database.CloudinaryCloudName, Database.CloudinaryAPIKey, Database.CloudinaryAPISecret)
	if err != nil {
		return "", err
	}

	//upload file
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: Database.CloudinaryUploadFolder})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
