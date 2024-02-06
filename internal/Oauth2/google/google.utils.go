package google

import (
	"bytes"
	"clean-architecture/bootstrap"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GetGoogleOauthToken(code string) (*OauthToken, error) {
	const rootURL = "https://oauth2.googleapis.com/token"

	app := bootstrap.App()
	env := app.Env
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", env.GoogleClientID)
	values.Add("client_secret", env.GoogleClientSecret)
	values.Add("redirect_uri", env.GoogleOAuthRedirectUrl)

	query := values.Encode()

	req, err := http.NewRequest("POST", rootURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody, &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &OauthToken{
		AccessToken: GoogleOauthTokenRes["access_token"].(string),
		IDToken:     GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetGoogleUser(accessToken string, idToken string) (*UserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", accessToken)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody, &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &UserResult{
		Id:              GoogleUserRes["id"].(string),
		Email:           GoogleUserRes["email"].(string),
		IsVerifiedEmail: GoogleUserRes["verified_email"].(bool),
		Name:            GoogleUserRes["name"].(string),
		GivenName:       GoogleUserRes["given_name"].(string),
		Picture:         GoogleUserRes["picture"].(string),
		Locale:          GoogleUserRes["locale"].(string),
	}

	return userBody, nil
}
