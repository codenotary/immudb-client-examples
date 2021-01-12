package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"github.com/codenotary/immudb/pkg/api/schema"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc/metadata"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"log"
)

// encryption AES example: https://github.com/moisoto/crypt/blob/master/crypt.go
func createPBKDF(key string, salt []byte) []byte {
	return pbkdf2.Key([]byte(key), salt, 4096, 32, sha1.New)
}

// randomSalt for use on calls to Encrypt and Decrypt functions
func randomSalt(size int) (salt []byte) {
	salt = make([]byte, size)
	return salt
}

// encrypts Values of a key-value map
func encryptMapValues(data map[string][]byte, password string, salt[]byte) map[string][]byte{
	block, _ := aes.NewCipher([]byte(createPBKDF(password,salt)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	var m map[string][]byte
	m = make(map[string][]byte)

	for key, value := range data {
		encryptedValue := gcm.Seal(nonce, nonce, value, nil)
		m[key] = encryptedValue
	}

	return m
}

// decrypts values
func decryptValue(data []byte, password string, salt[]byte) []byte {
	key := []byte(createPBKDF(password, salt))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, encryptedValue := data[:nonceSize], data[nonceSize:]
	decryptedValue, err := gcm.Open(nil, nonce, encryptedValue, nil)
	if err != nil {
		log.Fatal(err)
	}
	return decryptedValue
}

// basic print of unencrypted map
func printMap(data map[string][]byte){
	for key, value := range data {
		fmt.Println("Key:", key, "=>", "value:", string(value))
	}
}

// prints encrypted Map by decrypting values
func printEncryptedMap(data map[string][]byte,password string,salt[]byte){
	for key, value := range data {
		fmt.Println("Key:", key,string(decryptValue(value,password,salt)))
	}
}



func main() {

	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// login with default username and password
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}

	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// creating new database
		err = client.CreateDatabase(ctx, &schema.Database{
			Databasename: "yellowcard",
		})
		if err != nil {
			log.Fatal(err)
		}

	// switch to database
	resp, err := client.UseDatabase(ctx, &schema.Database{
		Databasename: "yellowcard",
	})

	md = metadata.Pairs("authorization", resp.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)


	// key-value map for personal data
	person := map[string][]byte{			"345623:name": []byte("Smarty McGopher"),
		"345623:date_of_birth":  []byte("20/10/1980")}

	var password = "SmartyMcGopherisSmart-and-encrypts!#vaccinationdata!"
	var encryptedPersonMap map[string][]byte=encryptMapValues(person,password,randomSalt(10))
	printEncryptedMap(encryptedPersonMap,password,randomSalt(10))

	// sets key-values of map in immudb
		for key, value := range encryptedPersonMap {
			if _, err := client.Set(ctx, []byte(key), value); err != nil {
				log.Fatal(err)
			}
		}

	for key := range encryptedPersonMap {
		if item, err := client.Get(ctx, []byte(key)); err != nil {
			log.Fatal(err)
		} else {
			// immudb sdk provides structured data. https://github.com/codenotary/immudb#structured-value
			fmt.Printf("%s\n", decryptValue(item.Value, password, randomSalt(10)))
		}
	}


	// key-value map for personal data
	vaccinations := map[string][]byte{	"345623:IPV:manufacturer": []byte("Vaccine int. corp"),
		"345623:HEPB:manufacturer" : []byte("Medco"),
		"345623:COV19:manufacturer": []byte("Gopher Immunization inc."),
		"345623:IPV:date_if_vaccination": []byte("15/09/2012"),
		"345623:HEPB:date_if_vaccination": []byte("03/06/2015"),
		"345623:COV19:date_if_vaccination": []byte("15/03/2021"),
		"345623:IPV:doctor_id": []byte("100234"),
		"345623:HEPB:doctor_id": []byte("100234"),
		"345623:COV19:doctor_id": []byte("100956"),
		"345623:IPV:product_id": []byte("1230"),
		"345623:HEPB:product_id": []byte("3309"),
		"345623:COV19:product_id":  []byte("1097")}
	var encryptedVaccinationMap map[string][]byte=encryptMapValues(vaccinations,password,randomSalt(10))
	printEncryptedMap(encryptedVaccinationMap,password,randomSalt(10))



	// sets key-values of map in immudb
	for key, value := range encryptedVaccinationMap {
		verifiedIndex, err := client.VerifiedSet(ctx, []byte(key), value)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("SafeSet - add and verify entry: \n", key, value, verifiedIndex.Id)

	}

	for key := range encryptedVaccinationMap {
		if item, err := client.Get(ctx, []byte(key)); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%s\n", decryptValue(item.Value, password, randomSalt(10)))
		}
	}
}
