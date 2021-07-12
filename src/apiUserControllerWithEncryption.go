package main

import (
	"crypto/aes"
	_ "crypto/aes"
	"crypto/cipher"
	_ "crypto/cipher"
	"crypto/rand"
	_ "crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"io/ioutil"
	"os"
	_ "strings"
)

type apiUserControllerWithEncryption struct {
	PAGINATION_LIMIT uint16
}

func (self*apiUserControllerWithEncryption) servers(ctx iris.Context) {
	type Server struct {
		Id 			int64				`json:"id"`
		Ip 			string				`json:"ip"`
		Location 	string				`json:"location"`
		Id_Server 	string				`json:"id_server"`
		Port		string				`json:"port"`
	}
	var serversArray []Server

	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serversArray)

	//msg(ctx,0,serversArray ,"servers")

	b, err := json.Marshal(serversArray)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	fmt.Println("Encryption Program v0.01")

	text := []byte(string(b))
	key := []byte("marsphrasewhichneedstobe32bytes!")

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	fmt.Println(gcm.Seal(nonce, nonce, text, nil))
	msg(ctx,0,gcm.Seal(nonce, nonce, text, nil) ,"servers")

	// the WriteFile method returns an error if unsuccessful
	err = ioutil.WriteFile("d://vpnfiles/myfile.data", gcm.Seal(nonce, nonce, text, nil), 0777)
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}
}

func (self*apiUserControllerWithEncryption) status(ctx iris.Context) {

}
func (self*apiUserControllerWithEncryption) readFromFile(ctx iris.Context) {
	fmt.Println("Decryption Program v0.01")

	key := []byte("marsphrasewhichneedstobe32bytes!")
	ciphertext, err := ioutil.ReadFile("d://vpnfiles/myfile.data")
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		fmt.Println(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plaintext))
}