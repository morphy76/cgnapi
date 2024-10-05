package command

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/morphy76/cgnapi/internal/configuration"
)

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

func GetToken(name string, decoded bool) error {
	profile, err := configuration.GetProfile(name)
	if err != nil {
		return err
	}

	if profile.CurrenAccessToken == "" {
		return fmt.Errorf("access token not initialized for profile %s", name)
	}

	if decoded {
		claims, err := parseToken(profile.CurrenAccessToken)
		if err != nil {
			return err
		}

		claimsJSON, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal claims: %v", err)
		}

		fmt.Println(string(claimsJSON))
	} else {
		fmt.Println(profile.CurrenAccessToken)
	}

	return nil
}

func GetTokenExp(name string) error {
	profile, err := configuration.GetProfile(name)
	if err != nil {
		return err
	}

	if profile.CurrenAccessToken == "" {
		return fmt.Errorf("access token not initialized for profile %s", name)
	}

	claims, err := parseToken(profile.CurrenAccessToken)
	if err != nil {
		return err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("failed to parse exp claim")
	}

	fmt.Println(time.Unix(int64(exp), 0).Format(time.RFC3339))

	return nil
}

func parseToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %v", err)
	}

	return claims, nil
}
