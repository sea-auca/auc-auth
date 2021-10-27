package model

import (
	"crypto"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	signAuthority = "https://sea.auca.kg/auth/configuration"
	audience      = "13ce98blgm00126gqm"
)

type JwtSigner struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

//Creates the struct, which stores signing keys in memory
func CreateSigner() *JwtSigner {
	pkey_bytes, err := ioutil.ReadFile("secret/auth_key.pem")
	if err != nil {
		log.Fatalln("failed to read the private key file, err msg: " + err.Error())
	}
	pbkey_bytes, err := ioutil.ReadFile("secret/auth_key.pub")
	if err != nil {
		log.Fatalln("failed to read the public key file, err msg: " + err.Error())
	}

	pkey, err := jwt.ParseEdPrivateKeyFromPEM(pkey_bytes)
	if err != nil {
		log.Fatalln("failed to parse the EdDSA private key, err msg: " + err.Error())
	}
	pbkey, err := jwt.ParseEdPublicKeyFromPEM(pbkey_bytes)
	if err != nil {
		log.Fatalln("failed to parse the EdDSA public key, err msg: " + err.Error())
	}

	return &JwtSigner{privateKey: pkey, publicKey: pbkey}
}

//Custom payload scheme, which hanldes standard claims and application specific data
type SeaAucaClaims struct {
	Payload string `json:"payload"`
	jwt.StandardClaims
}

//Creates new token with specified ttl and payload data
func (signer *JwtSigner) NewToken(ttl time.Duration) string {

	claims := SeaAucaClaims{
		Payload: "payload",
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Id:        "1",
			Issuer:    signAuthority,
			Audience:  audience,
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(signer.privateKey)
	if err != nil {
		log.Println("Failed to sign the jwt token, error message: " + err.Error())
	}
	return token
}

type JWTError struct {
	message string
}

func newJWTError(msg string) JWTError {
	return JWTError{message: "JWT token error: " + msg}
}

func (err JWTError) Error() string {
	return err.message
}

//Handles the token validation logic and claims extraction
func (signer *JwtSigner) Validate(token string) (*SeaAucaClaims, error) {
	claims := SeaAucaClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, newJWTError("invalid signing method, method: " + t.Method.Alg())
		}

		_, ok := t.Claims.(*SeaAucaClaims)
		if !ok {
			return nil, newJWTError("claims can not be parsed")
		}

		return signer.publicKey, nil
	})
	if !t.Valid {
		return nil, err
	}

	return &claims, nil
}
