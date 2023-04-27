package almanack

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/carlmjohnson/crockford"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/resperr"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spotlightpa/almanack/internal/aws"
)

func GetSignedImageUpload(ctx context.Context, is aws.BlobStore, ct string) (signedURL, filename string, err error) {
	filename = makeImageName(ct)
	h := make(http.Header, 1)
	h.Set("Content-Type", ct)
	signedURL, err = is.SignPutURL(ctx, filename, h)
	return
}

func makeImageName(ct string) string {
	now := time.Now().UTC()
	ext := "bin"
	if _, tempext, ok := strings.Cut(ct, "/"); ok && len(tempext) >= 3 {
		ext = tempext
	}
	var sb strings.Builder
	sb.Grow(len("2006/01/1234-5678-9abc-defg.") + len(ext))
	sb.WriteString(now.Format("2006/01/"))
	buf := make([]byte, 0, len("1234-5678-9abc-defg"))
	buf = crockford.AppendTime(crockford.Lower, now, buf)
	buf = crockford.AppendRandom(crockford.Lower, buf)
	buf = crockford.AppendPartition(buf[:0], buf, 4)
	sb.Write(buf)
	sb.WriteByte('.')
	sb.WriteString(ext)
	return sb.String()
}

func hashURLpath(srcPath, ext string) string {
	return fmt.Sprintf("external/%s.%s",
		crockford.MD5(crockford.Lower, []byte(srcPath)),
		ext,
	)
}

func FetchImageURL(ctx context.Context, c *http.Client, srcurl string) (body []byte, ctype string, err error) {
	var buf bytes.Buffer
	if err = requests.
		URL(srcurl).
		Client(c).
		CheckStatus(http.StatusOK).
		CheckPeek(512, func(peek []byte) error {
			ct := mimetype.Detect(peek)
			if ct.Is("image/jpeg") ||
				ct.Is("image/png") ||
				ct.Is("image/tiff") ||
				ct.Is("image/webp") ||
				ct.Is("image/avif") ||
				ct.Is("image/heic") {
				ctype = ct.String()
				return nil
			}
			return resperr.WithUserMessage(
				fmt.Errorf("%q did not have proper MIME type: %s",
					srcurl, ct.String()),
				"URL must be an image",
			)
		}).
		ToBytesBuffer(&buf).
		Fetch(ctx); err != nil {
		if resperr.StatusCode(err) != http.StatusInternalServerError {
			return nil, "", err
		}
		return nil, "", resperr.New(http.StatusBadGateway, "problem fetching image: %w", err)
	}

	return buf.Bytes(), ctype, nil
}
