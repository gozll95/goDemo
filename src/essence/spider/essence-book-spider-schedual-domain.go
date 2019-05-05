1."正则"

var regexpForIP = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|[01]?\d?\d)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d?\d))`)

var regexpForDomains = []*regexp.Regexp{
	// *.xx or *.xxx.xx
	regexp.MustCompile(`\.(com|com\.\w{2})$`),
	regexp.MustCompile(`\.(gov|gov\.\w{2})$`),
	regexp.MustCompile(`\.(net|net\.\w{2})$`),
	regexp.MustCompile(`\.(org|org\.\w{2})$`),
	// *.xx
	regexp.MustCompile(`\.me$`),
	regexp.MustCompile(`\.biz$`),
	regexp.MustCompile(`\.info$`),
	regexp.MustCompile(`\.name$`),
	regexp.MustCompile(`\.mobi$`),
	regexp.MustCompile(`\.so$`),
	regexp.MustCompile(`\.asia$`),
	regexp.MustCompile(`\.tel$`),
	regexp.MustCompile(`\.tv$`),
	regexp.MustCompile(`\.cc$`),
	regexp.MustCompile(`\.co$`),
	regexp.MustCompile(`\.\w{2}$`),
}


	if regexpForIP.MatchString(host) {
		return host, nil
	}

	for _, re := range regexpForDomains {
		pos := re.FindStringIndex(host)
		if pos != nil {
			suffixIndex = pos[0]
			break
		}
	}