package almanack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/errutil"
	"github.com/carlmjohnson/resperr"
	"github.com/spotlightpa/almanack/internal/db"
)

func (svc Service) ReplaceImageURL(ctx context.Context, srcURL, description, credit string) (s string, err error) {
	defer errutil.Trace(&err)

	if srcURL == "" {
		return "", fmt.Errorf("no image provided")
	}
	image, err := svc.Queries.GetImageBySourceURL(ctx, srcURL)
	if err != nil && !db.IsNotFound(err) {
		return "", err
	}
	if !db.IsNotFound(err) && image.IsUploaded {
		return image.Path, nil
	}
	var path, ext string
	if path, ext, err = UploadFromURL(ctx, svc.Client, svc.ImageStore, srcURL); err != nil {
		return "", resperr.New(
			http.StatusBadGateway,
			"could not upload image %s: %w", srcURL, err,
		)
	}
	_, err = svc.Queries.CreateImage(ctx, db.CreateImageParams{
		Path:        path,
		Type:        ext,
		Description: description,
		Credit:      credit,
		SourceURL:   srcURL,
		IsUploaded:  true,
	})
	return path, err
}

func (svc Service) UpdateMostPopular(ctx context.Context) (err error) {
	defer errutil.Trace(&err)

	svc.Logger.Printf("updating most popular")
	cl, err := svc.gsvc.GAClient(ctx)
	if err != nil {
		return err
	}
	pages, err := svc.gsvc.MostPopularNews(ctx, cl)
	if err != nil {
		return err
	}
	data := struct {
		Pages []string `json:"pages"`
	}{pages}
	return UploadJSON(
		ctx,
		svc.FileStore,
		"feeds/most-popular.json",
		"public, max-age=300",
		&data,
	)
}