package saml

type SamlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (x *SamlConfig) IsValid() bool {
	return x.Username != "" && x.Password != ""
}
