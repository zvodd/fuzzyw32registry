package main

import (
	"log"

	reg32 "fuzzyw32registry"

	"golang.design/x/clipboard"
	"golang.org/x/sys/windows/registry"
)

const REGKEY_LASTOPEN = `SOFTWARE\Microsoft\Windows\CurrentVersion\Applets\Regedit\`

func main() {
	log.SetFlags(0)
	log.Printf("Accessing registry key:\n\t"+`"%s"`+"\n", REGKEY_LASTOPEN)
	klastOpenR, err := registry.OpenKey(registry.CURRENT_USER, REGKEY_LASTOPEN, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}

	location, _, err := klastOpenR.GetStringValue("LastKey")
	if err != nil {
		log.Fatal(err)
	}

	klastOpenR.Close()

	log.Printf("Regedit current location:\n\t"+`"%s"`+"\n", location)
	clipkey := string(clipboard.Read(clipboard.FmtText))

	prefix, key, err := reg32.FuzzyRegKey(clipkey)
	if err != nil {
		log.Fatal(err)
	}
	fullkey, err := reg32.AddPrefix(prefix, key)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found reg Key:\n\t"+`"%s"`+"\n", fullkey)

	fullkey = `Computer\` + fullkey
	klastOpenW, err := registry.OpenKey(registry.CURRENT_USER, REGKEY_LASTOPEN, registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	klastOpenW.SetStringValue("LastKey", fullkey)
}
