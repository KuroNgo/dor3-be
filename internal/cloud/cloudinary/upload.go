package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
)

func UploadToCloudinary(file multipart.File, filePath string, folder string) (Upload, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return Upload{}, err
	}

	uploadParams := uploader.UploadParams{
		PublicID: filePath,
		Folder:   folder,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return Upload{}, err
	}

	resultRes := Upload{
		ImageURL: result.SecureURL,
		AssetID:  result.AssetID,
	}
	return resultRes, nil
}
