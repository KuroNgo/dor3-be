package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func DeleteToCloudinary(filename string) (string, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return "", err
	}

	uploadParams := uploader.DestroyParams{
		PublicID: filename,
	}

	result, err := cld.Upload.Destroy(ctx, uploadParams)
	if err != nil {
		return "", err
	}

	return result.Result, nil
}
