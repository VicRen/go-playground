package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gofrs/uuid"
)

const (
	defUUIDFileName = "uuid.txt"
)

func main() {
	uuidPath := "/usr/local/var/galaxy"
	fmt.Println("reading server UUID from /usr/local/var/galaxy/uuid.txt")
	fullPath := path.Join(uuidPath, defUUIDFileName)
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("failed to read server UUID: %v\n", err)
	}
	var uuidBuf [16]byte
	copy(uuidBuf[:], data)
	uuidValue, err := UnmarshalUUID(data)
	if err != nil {
		fmt.Printf("failed to read server UUID: %v\n", err)
	}
	fmt.Printf("server UUID: %s(%v)\n", uuidValue, uuidBuf)

	uuidPath = "/usr/local/var/galaxy/client"
	fmt.Println("reading client UUID from /usr/local/var/galaxy/client/uuid.txt")
	fullPath = path.Join(uuidPath, defUUIDFileName)
	data, err = ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("failed to read client UUID: %v\n", err)
	}
	copy(uuidBuf[:], data)
	uuidValue, err = UnmarshalUUID(data)
	if err != nil {
		fmt.Printf("failed to read client UUID: %v\n", err)
	}
	fmt.Printf("client UUID: %s(%v)\n", uuidValue, uuidBuf)
}

// UnmarshalUUID format data to string format id.
func UnmarshalUUID(id []byte) (string, error) {
	uid, err := uuid.FromBytes(id)
	if err != nil {
		return "", err
	}
	return uid.String(), err
}
