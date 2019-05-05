package credential

// A Value is the AWS credential value for individual credential fields.
type Value struct {
	AccessKeyId     string // aws access key id
	AccessKeySecret string // aws access secret key

	SessionToken string // aws session token
}
