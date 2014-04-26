func SignV2(creds Creds, service Service, req map[string]string,
	date time.Time, expires bool) map[string]string {

	params := map[string]string{
		"AWSAccessKeyId":   creds.Access,
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
		"Version":          service.Version,
	}
	//... minor code elided, copy req into params
	// handle SecurityToken, Expires vs Timestamp
	toSign := strings.Join([]string{
		"GET", // Method
		strings.ToLower(service.Endpoint),
		"/", // Path
		QueryString(params)}, "\n")
	signature := hmac.New(sha256.New, []byte(creds.Secret))
	signature.Write([]byte(toSign))
	params["Signature"] = base64.StdEncoding.EncodeToString(signature.Sum(nil))
	return params
}
