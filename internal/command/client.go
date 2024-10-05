package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/morphy76/cgnapi/internal/configuration"
)

func ValidateToken(name string) error {
	profile, err := configuration.GetProfile(name)
	if err != nil {
		return err
	}

	if profile.RefreshToken == "" {
		return fmt.Errorf("refresh token not initialized for profile %s", name)
	}

	if profile.CurrenAccessToken == "" {
		return fmt.Errorf("access token not initialized for profile %s", name)
	}

	// introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", profile.AuthServer, profile.Realm)
	// data := fmt.Sprintf("token=%s", profile.CurrenAccessToken)
	// req, err := http.NewRequest("POST", introspectionURL, bytes.NewBufferString(data))
	// if err != nil {
	// 	return fmt.Errorf("failed to create introspection request: %v", err)
	// }
	// req.SetBasicAuth(profile.ClientID, profile.ClientSecret)
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("failed to perform introspection request: %v", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("introspection request failed with status: %s", resp.Status)
	// }

	// var introspectionResponse struct {
	// 	Active bool `json:"active"`
	// }

	// if err := json.NewDecoder(resp.Body).Decode(&introspectionResponse); err != nil {
	// 	return fmt.Errorf("failed to decode introspection response: %v", err)
	// }

	// if !introspectionResponse.Active {
	// 	return fmt.Errorf("access token is not active")
	// }

	return nil
}

func RenewToken(name string) error {
	profile, err := configuration.GetProfile(name)
	if err != nil {
		return err
	}

	if profile.RefreshToken == "" {
		return fmt.Errorf("refresh token not initialized for profile %s", name)
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", profile.AuthServer, profile.Realm)
	data := fmt.Sprintf("client_id=%s&grant_type=refresh_token&refresh_token=%s", profile.ClientID, profile.RefreshToken)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to create token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status: %s", resp.Status)
	}

	var tokenRenewResponse struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		ErrorURI         string `json:"error_uri"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenRenewResponse); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	if tokenRenewResponse.Error != "" {
		return fmt.Errorf("failed to renew token: %s", tokenRenewResponse.ErrorDescription)
	}

	profile.CurrenAccessToken = tokenRenewResponse.AccessToken
	profile.RefreshToken = tokenRenewResponse.RefreshToken

	if err := configuration.UpdateProfile(name, profile); err != nil {
		return err
	}

	return nil
}
