package src

type Resps struct {
	url        string
	body       string
	header     map[string][]string
	server     string
	statuscode int
	length     int
	title      string
	mmh3icon   string
	md5icon    string
	jsfiles    []string
}

type Fingerprint struct {
	Fp      string   `json:"fp"`
	Headers []string `json:"headers"`
	Body    []string `json:"body"`
	Icon    []string `json:"icon"`
	JS      []string `json:"js"`
	Title   []string `json:"title"`
	Regexp  string   `json:"regexp"`
}

type Data struct {
	Fingerprints []Fingerprint `json:"fingerprints"`
}
