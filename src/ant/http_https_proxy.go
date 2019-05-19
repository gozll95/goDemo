/*
 * 获取http proxy
 */
func (client *Client) getHttpProxy(scheme string) (proxy *url.URL, err error) {
	if scheme == "https" {
		if client.GetHttpsProxy() != "" {
			proxy, err = url.Parse(client.httpsProxy)
		} else if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if client.GetHttpProxy() != "" {
			proxy, err = url.Parse(client.httpProxy)
		} else if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}

	return proxy, err
}


/*
 * http proxy的使用
 */
	 // Set whether to ignore certificate validation.
	 // Default InsecureSkipVerify is false.
	 if trans, ok := client.httpClient.Transport.(*http.Transport); ok && trans != nil {
		trans.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: client.getHTTPSInsecure(request),
		}
		if proxy != nil && !flag {
			trans.Proxy = http.ProxyURL(proxy)
		}
		client.httpClient.Transport = trans
	}

	var httpResponse *http.Response
	for retryTimes := 0; retryTimes <= client.config.MaxRetryTime; retryTimes++ {
		if proxy != nil && proxy.User != nil {
			if password, passwordSet := proxy.User.Password(); passwordSet {
				httpRequest.SetBasicAuth(proxy.User.Username(), password)
			}
		}