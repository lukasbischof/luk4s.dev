package hcaptcha

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const URL = "https://hcaptcha.com/siteverify"

func Verify(response string) (bool, error) {
	secret := os.Getenv("HCAPTCHA_SECRET_KEY")

	resp, err := http.Post(
		URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("secret=%s&response=%s", secret, response)),
	)

	if err != nil {
		return false, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result["success"] == true, nil
}
