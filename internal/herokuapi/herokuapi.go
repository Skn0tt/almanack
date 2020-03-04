package herokuapi

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func FlagVar(fl *flag.FlagSet) func() (connectionURL string, err error) {
	apiKey := fl.String("heroku-api-key", "", "`API key` for retreiving config from Heroku")
	addOn := fl.String("heroku-add-on-id", "", "`ID` for Heroku Add-On to get config from")
	return func() (connectionURL string, err error) {
		if *apiKey == "" || *addOn == "" {
			return "", nil
		}

		return Request(*apiKey, *addOn)
	}
}

const timeout = 5 * time.Second

func Request(apiKey, addOn string) (connectionURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.heroku.com/addons/"+
		addOn+"/config", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var values []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err = json.Unmarshal(b, &values); err != nil {
		return
	}

	for _, val := range values {
		if val.Name == "URL" || val.Name == "url" {
			connectionURL = val.Value
			err = nil
			return
		}
	}
	return "", fmt.Errorf("could not find connection URL for %s", addOn)
}
