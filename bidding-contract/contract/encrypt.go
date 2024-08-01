package contract

import (
	"encoding/pem"
	"fmt"
	"log"
	"os"

	ecies "github.com/ecies/go/v2"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func ConvertSecp256k1ToEcies(pubKey *secp256k1.PublicKey) (*ecies.PublicKey, error) {
	// Extract the X and Y coordinates by calling the functions
	x := pubKey.X()
	y := pubKey.Y()

	// Create an ECIES public key from the X and Y coordinates
	eciesPubKey := &ecies.PublicKey{
		X:     x,
		Y:     y,
		Curve: secp256k1.S256(),
	}

	return eciesPubKey, nil
}
func Ecies_encryption(pubkey_path string, data []byte) (ciphertext []byte) {
	read_pubKey, err := os.ReadFile(pubkey_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("publickey which is read from given pubkey.pem file is ", read_pubKey)

	pemdecoded_pubkey, rest := pem.Decode(read_pubKey)
	fmt.Println("pemdecodedpublic key is  ", pemdecoded_pubkey)
	fmt.Println("rest part is ", rest)
	pubkeyback, _ := secp256k1.ParsePubKey(pemdecoded_pubkey.Bytes)
	eciesPubKey, err := ConvertSecp256k1ToEcies(pubkeyback)
	if err != nil {
		fmt.Println("Error converting public key:", err)
		return
	}

	ciphertext, err = ecies.Encrypt(eciesPubKey, data)
	if err != nil {
		panic(err)
	}
	//fmt.Println("ciphertext is  ", ciphertext)
	return ciphertext
}
