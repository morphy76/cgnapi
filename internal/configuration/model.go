package configuration

type Config struct {
	Profiles map[string]Profile `yaml:"profiles,omitempty"`
}

type Profile struct {
	AuthServer        string `yaml:"auth_server,omitempty"`
	ClientID          string `yaml:"client_id,omitempty"`
	CurrenAccessToken string `yaml:"current_access_token,omitempty"`
	Realm             string `yaml:"realm,omitempty"`
	RefreshToken      string `yaml:"refresh_token,omitempty"`
}
