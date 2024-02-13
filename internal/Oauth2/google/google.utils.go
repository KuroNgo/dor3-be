package google

import (
	"bytes"
	"clean-architecture/bootstrap"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GetGoogleOauthToken Retrieve the OAuth2 Access Token
func GetGoogleOauthToken(code string) (*OauthToken, error) {
	const rootURL = "https://oauth2.googleapis.com/token"

	app := bootstrap.App()
	env := app.Env
	values := url.Values{}
	// grant_type is the type of grant being trequested, which is typically authorization_code
	values.Add("grant_type", "authorization_code")

	// the authorization code obtained from the authorization endpoint
	values.Add("code", code)

	// the secret associated with the client ID
	values.Add("client_id", env.GoogleClientID)
	values.Add("client_secret", env.GoogleClientSecret)

	// the authorized callback URL registered with the client
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

	//resBody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, err
	//}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &OauthToken{
		AccessToken: GoogleOauthTokenRes["access_token"].(string),
		IDToken:     GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

// GetGoogleUser Get the Google User's Account Information
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

	//resBody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, err
	//}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
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
