package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
)

func UploadImageToCloudinary(file multipart.File, filePath string, folder string) (UploadImage, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return UploadImage{}, err
	}

	uploadParams := uploader.UploadParams{
		PublicID: filePath,
		Folder:   folder,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return UploadImage{}, err
	}

	resultRes := UploadImage{
		ImageURL: result.SecureURL,
		AssetID:  result.AssetID,
	}
	return resultRes, nil
}

func UploadAudioToCloudinary(file multipart.File, filePath string, folder string) (UploadAudio, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return UploadAudio{}, err
	}

	uploadParams := uploader.UploadParams{
		PublicID: filePath,
		Folder:   folder,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return UploadAudio{}, err
	}

	resultRes := UploadAudio{
		AudioURL: result.SecureURL,
		AssetID:  result.AssetID,
	}
	return resultRes, nil
}
