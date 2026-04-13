package cloudinary

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	cld, err := cloudinary.New()
	if err != nil {
		return "", err
	}
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   "anime",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func DeleteImage(ctx context.Context, imageUrl string) error {
	cld, err := cloudinary.New()
	if err != nil {
		return err
	}
	publicID := extractPublicID(imageUrl)
	_, err = cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: []string{publicID},
	})
	return err
}

func extractPublicID(imageUrl string) string {
	// URL format: https://res.cloudinary.com/{cloud}/image/upload/v{version}/{folder}/{file}.{ext}
	parts := strings.Split(imageUrl, "/upload/")
	if len(parts) < 2 {
		return ""
	}
	// Remove version (v1234567/) and extension
	withoutVersion := strings.SplitN(parts[1], "/", 2)
	if len(withoutVersion) < 2 {
		return ""
	}
	// Remove file extension
	dotIndex := strings.LastIndex(withoutVersion[1], ".")
	if dotIndex == -1 {
		return withoutVersion[1]
	}
	return withoutVersion[1][:dotIndex]
}
