package configuration

type Config struct {
	Profiles map[string]Profile `yaml:"profiles,omitempty"`
}

type Profile struct {
	AuthServer        string `yaml:"auth_server,omitempty"`
	RefreshToken      string `yaml:"refresh_token,omitempty"`
	CurrenAccessToken string `yaml:"current_access_token,omitempty"`
}
