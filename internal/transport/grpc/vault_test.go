package grpc

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestAbs(t *testing.T) {

	test := "-----BEGIN PUBLIC KEY-----\n" +
		"MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuII+UtaMVIjFnUfyubqwPQPtF2WuPNPpdQyqJUoVHOuZ8b3Q1rQrAXVnt5jIbolfUIXpsurkgawEoeDNVxmoIaxr3sA44d17UqzcXxXTW3xGt4XXkIyZrDhZOi+wO6gebd1/LPe2MYbodfhAa1koBZzQLqXRvKk2M4LyHO8ndwN/3j2fnjHMrgQXv7nKkeVgb/b9K7ciVA4ftmJe+5ZYsH3JNVuQTi8b257xpgCYjQHtP130hyX0ysg8weCUbzUXN3a6R4zbk4iTby2xFKWjE3AxVe9Vqlh+XtH7rREOvthIej0znDb0VuBkssIkf/nx3Izqmm9NaitTb8AZTHlTUT8ztcuhsOFVLeIpVakrnvJ1dS/CxltewlsfYchGti0Y8CKohMEG7Rp+DUMyKRkeOz3bSAzyeAdQVBsb2JJij/drNuAMPEIE8mRY/70sBtYG2BGpqlvVxY7jf2IJKrSIiqfIhZhMbhjArCWjq1aS+4+ZcY6m4AqsOwOQu/VuupZ8Xb7+SQxdB5WAtsJRQCHnFyfkz77BR0GLEjqsqRPwOoCbNcI+oM9RiJgYVkeanLJQ5lFhmwbJmbphHM6wRJ8bNVKqkG564NLXspsJakEmTQHFt1inBbnB4jJ6W+XnYCFlisuD/hpYhSozZ5mUqksiT0pQMLZ2RqFG6WQMNLS1O2cCAwEAAQ==" +
		"\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode([]byte(test))

	_, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Errorf("%s", err)
	}
}
