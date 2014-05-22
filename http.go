package chef

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"
)

// ChefVersion that we pretend to emulate
const ChefVersion = "11.12.0"

// AuthConfig representing a client and a private key used for encryption
type AuthConfig struct {
	privateKey *rsa.PrivateKey
	clientName string
	cryptoHash crypto.Hash
}

// SignRequest modifies headers of an http.Request
func (ac AuthConfig) SignRequest(request *http.Request, cryptoHash crypto.Hash) error {
	// sanitize the path for the chef-server
	// chef-server doesn't support '//' in the Hash Path.
	var endpoint string
	if request.URL.Path != "" {
		endpoint = path.Clean(request.URL.Path)
		request.URL.Path = endpoint
	} else {
		endpoint = request.URL.Path
	}

	request.Header.Set("Method", request.Method)
	// Since we're not setting the encoded slice limit
	// we can safely call out [0]
	generatedHash := generateHash(endpoint, cryptoHash)
	if generatedHash == nil {
		return errors.New("Unsupported crypto hashing algorithm")
	}
	request.Header.Set("Hashed Path", base64BlockEncode(generatedHash, 0)[0])
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Chef-Version", ChefVersion)
	request.Header.Set("X-Ops-Timestamp", time.Now().UTC().Format(time.RFC3339))
	request.Header.Set("X-Ops-Userid", ac.clientName)

	var xOpsSignCrypto string
	switch cryptoHash {
	case crypto.MD5:
		xOpsSignCrypto = "algorithm=md5;"
	case crypto.SHA1:
		xOpsSignCrypto = "algorithm=sha1;"
	case crypto.SHA224:
		xOpsSignCrypto = "algorithm=sha224;"
	case crypto.SHA256:
		xOpsSignCrypto = "algorithm=sha256;"
	case crypto.SHA384:
		xOpsSignCrypto = "algorithm=sha384;"
	case crypto.SHA512:
		xOpsSignCrypto = "algorithm=sha512;"
	}

	request.Header.Set("X-Ops-Sign", xOpsSignCrypto+"version=1.0")
	request.Header.Set("X-Ops-Content-Hash", calcBodyHash(request, cryptoHash))

	// To validate the signature it seems to be very particular
	// Would like to use this loop to generate the content
	// But it causes validation to fail.. so we do it explicitly

	// var content string
	// for key, value := range request.Header {
	// 	content += fmt.Sprintf("%s:%s\n", key, request.Header.Get(key))
	// }
	var content string
	content += fmt.Sprintf("%s:%s\n", "Method", request.Header.Get("Method"))
	content += fmt.Sprintf("%s:%s\n", "Hashed Path", request.Header.Get("Hashed Path"))
	content += fmt.Sprintf("%s:%s\n", "Accept", request.Header.Get("Accept"))
	content += fmt.Sprintf("%s:%s\n", "X-Chef-Version", request.Header.Get("X-Chef-Version"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Timestamp", request.Header.Get("X-Ops-Timestamp"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Userid", request.Header.Get("X-Ops-Userid"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Sign", request.Header.Get("X-Ops-Sign"))
	content += fmt.Sprintf("%s:%s\n", "X-Ops-Content-Hash", request.Header.Get("X-Ops-Content-Hash"))

	// generate signed string of headers
	// Since we've gone through additional validation steps above,
	// we shouldn't get an error at this point
	signature, _ := generateSignature(ac.privateKey, content, cryptoHash)

	// TODO: THIS IS CHEF PROTOCOL SPECIFIC
	// Signature is made up of n 60 length chunks
	base64sig := base64BlockEncode(signature, 60)

	// roll over the auth slice and add the apropriate header
	for index, value := range base64sig {
		request.Header.Set(fmt.Sprintf("X-Ops-Authorization-%d", index+1), string(value))
	}

	return nil
}

// modified from goiardi calcBodyHash
func calcBodyHash(r *http.Request, cryptoHash crypto.Hash) string {
	var bodyStr string

	if r.Body == nil {
		bodyStr = ""
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		bodyStr = buf.String()
	}

	// Since we're not setting the encoded slice limit
	// we can safely call out [0]
	chkHash := base64BlockEncode(generateHash(bodyStr, cryptoHash), 0)
	return chkHash[0]
}