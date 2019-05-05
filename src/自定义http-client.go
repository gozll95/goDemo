type AuthAccountOpts struct {
	Host        string        `json:"-"`
	MgoURL      string        `json:"-"`
	Username    string        `json:"-"`
	Password    string        `json:"-"`
	ExpiresAt   time.Time     `json:"-"`
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
}

func (opts *AuthAccountOpts) Client() *http.Client {
	return &http.Client{
		Transport: opts,
	}
}




func (opts *AuthAccountOpts) RoundTrip(req *http.Request) (*http.Response, error) {
	if time.Now().After(opts.ExpiresAt) {
		resp, err := http.PostForm(opts.Host+"/oauth2/token", url.Values{
			"grant_type": {"password"},
			"username":   {opts.Username},
			"password":   {opts.Password},
		})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("acc response status %d", resp.StatusCode)
		}
		err = json.NewDecoder(resp.Body).Decode(opts)
		if err != nil {
			return nil, err
		}
		opts.ExpiresAt = time.Now().Add(time.Second * opts.ExpiresIn)
	}
	req.Header.Set("Authorization", "Bearer "+opts.AccessToken)
	return http.DefaultTransport.RoundTrip(req)
}

