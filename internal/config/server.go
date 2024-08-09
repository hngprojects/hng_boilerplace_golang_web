package config

type ServerConfiguration struct {
	Port                      string
	Secret                    string
	AccessTokenExpireDuration int
	RequestPerSecond          float64
	TrustedProxies            []string
	ExemptFromThrottle        []string
}

type App struct {
	Name                  string
	Mode                  string
	Url                   string
	MagicLinkDuration     int
	ResetPasswordDuration int
}
