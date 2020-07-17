package almanack

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/carlmjohnson/crockford"
	"github.com/carlmjohnson/resperr"
	"github.com/spotlightpa/almanack/pkg/common"
	"golang.org/x/net/context/ctxhttp"
)

func GetSignedUpload(is common.ImageStore, ext string) (signedURL, filename string, err error) {
	filename = makeFilename(ext)
	signedURL, err = is.GetSignedURL(filename)
	return
}

func GetSignedHashedUrl(is common.ImageStore, srcurl, ext string) (signedURL, filename string, err error) {
	filename = hashURLpath(srcurl, ext)
	signedURL, err = is.GetSignedURL(filename)
	return
}

func makeFilename(ext string) string {
	var sb strings.Builder
	sb.Grow(len("2006/01/123456789abcdefg.") + len(ext))
	sb.WriteString(time.Now().Format("2006/01/"))
	sb.Write(crockford.Time(crockford.Lower, time.Now()))
	sb.Write(crockford.AppendRandom(crockford.Lower, nil))
	sb.WriteString(".")
	sb.WriteString(ext)
	return sb.String()
}

func hashURLpath(srcPath, ext string) string {
	return fmt.Sprintf("external/%s.%s",
		crockford.AppendMD5(crockford.Lower, nil, []byte(srcPath)),
		ext,
	)
}

func UploadFromURL(ctx context.Context, c *http.Client, is common.ImageStore, srcurl string) (filename, ext string, err error) {
	res, err := ctxhttp.Get(ctx, c, srcurl)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	const (
		megabyte = 1 << 20
		maxSize  = 25 * megabyte
	)
	buf := bufio.NewReader(http.MaxBytesReader(nil, res.Body, maxSize))
	// http.DetectContentType only uses first 512 bytes
	peek, err := buf.Peek(512)
	if err != nil && err != io.EOF {
		return "", "", err
	}

	if ct := http.DetectContentType(peek); strings.HasPrefix(ct, "image/jpeg") {
		ext = "jpeg"
	} else if strings.HasPrefix(ct, "image/png") {
		ext = "png"
	} else {
		return "", "", resperr.WithCodeAndMessage(
			fmt.Errorf("%q did not have proper MIME type", srcurl),
			http.StatusBadRequest,
			"URL must be an image",
		)
	}

	slurp, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", "", err
	}

	var signedURL string
	signedURL, filename, err = GetSignedHashedUrl(is, srcurl, ext)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest(http.MethodPut, signedURL, bytes.NewReader(slurp))
	if err != nil {
		return "", "", err
	}
	res, err = ctxhttp.Do(ctx, c, req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", resperr.WithCodeAndMessage(
			fmt.Errorf("unexpected S3 status: %d", res.StatusCode),
			http.StatusBadGateway,
			"Could not upload image from URL",
		)
	}
	return filename, ext, nil
}
