package google

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/spotlightpa/almanack/internal/stringutils"
)

func TestListDriveFiles(t *testing.T) {
	svc := Service{}
	svc.l = log.Default()
	svc.driveID = stringutils.First(os.Getenv("ALMANACK_GOOGLE_TEST_DRIVE"), "1")
	ctx := context.Background()
	cl := *http.DefaultClient
	cl.Transport = requests.Replay("testdata")
	if os.Getenv("ALMANACK_GOOGLE_TEST_RECORD_REQUEST") != "" {
		gcl, err := svc.DriveClient(ctx)
		if err != nil {
			t.Fatal(err)
		}
		cl.Transport = requests.Record(gcl.Transport, "testdata")
	}
	files, err := svc.Files(ctx, &cl)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 1 {
		t.Fatal("got no files")
	}
}
