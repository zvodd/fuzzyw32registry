package main

import (
	"fmt"
	"log"

	reg32 "fuzzyw32registry"

	"golang.design/x/clipboard"
	"golang.org/x/sys/windows/registry"
)

const REGKEY_LASTOPEN = `Software\Microsoft\Windows\CurrentVersion\Applets\Regedit\`

func main() {
	fmt.Printf("Accessing registry key:\n\t"+`"%s"`+"\n", REGKEY_LASTOPEN)
	k, err := registry.OpenKey(registry.CURRENT_USER, REGKEY_LASTOPEN, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	location, _, err := k.GetStringValue("LastKey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Regedit current location:\n\t"+`"%s"`+"\n", location)
	clipkey := string(clipboard.Read(clipboard.FmtText))

	// prekey, key, err := reg32.AttemptPrefix(clipkey)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// key, err = reg32.LocateFirstValid(prekey, key)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	prefix, key, err := reg32.FuzzyRegKey(clipkey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found reg Key:\n\t"+`%+v "%s"`+"\n", *prefix, key)
}
