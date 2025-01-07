package echogy

type Config struct {
	HttpAddr   string `json:"httpAddr"`
	SSHAddr    string `json:"SSHAddr"`
	Domain     string `json:"domain"`
	PrivateKey string `json:"privateKey"`
}
