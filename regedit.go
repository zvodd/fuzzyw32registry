package fuzzyw32registry

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var ERR_BASE_END = errors.New("EOF")

func FuzzyRegKey(inkey string) (*registry.Key, string, error) {
	inkey = strings.Trim(inkey, "\n\r \t")
	inkey = strings.TrimPrefix(inkey, "Computer")
	inkey = strings.TrimPrefix(inkey, `\`)
	prekey, key, err := AttemptPrefix(inkey)
	if err != nil {
		return nil, key, err
	}
	key, err = LocateFirstValid(prekey, key)
	if err != nil {

		return nil, key, err
	}
	return prekey, key, nil
}

func AttemptPrefix(str string) (*registry.Key, string, error) {
	//TODO check length
	parts := strings.Split(str, `\`)
	keystart, therestparts := parts[0], parts[1:]
	therest := strings.Join(therestparts, `\`)
	keystart = strings.ToUpper(keystart)
	log.Printf(`
	parts: %v
	keystart: %s
	therest: %s
	`, parts, keystart, therest)

	ignore_start := 0
	for _, pre := range []string{"HKEY", "HK"} {
		if len(keystart) <= len(pre) {
			break
		}
		if strings.HasPrefix(keystart, pre) {
			ignore_start = len(pre)
			if string(keystart[len(pre)]) == "_" {
				if len(keystart) > len(pre)+1 {
					ignore_start += 1
				}
			}
			break
		}
	}
	reduced_keystart := keystart[ignore_start:]

	log.Printf(`
	keystart: %v [%v]
	reduced_keystart: %v
	`, keystart, ignore_start, reduced_keystart)

	kmap := map[string]registry.Key{
		"CLASSES_ROOT":     registry.CLASSES_ROOT,
		"CURRENT_USER":     registry.CURRENT_USER,
		"LOCAL_MACHINE":    registry.LOCAL_MACHINE,
		"USERS":            registry.USERS,
		"CURRENT_CONFIG":   registry.CURRENT_CONFIG,
		"PERFORMANCE_DATA": registry.PERFORMANCE_DATA,
		"CR":               registry.CLASSES_ROOT,
		"CU":               registry.CURRENT_USER,
		"LM":               registry.LOCAL_MACHINE,
		"U":                registry.USERS,
		"CC":               registry.CURRENT_CONFIG,
	}
	for k, v := range kmap {
		if strings.HasPrefix(reduced_keystart, k) {
			return &v, therest, nil
		}
	}
	return nil, "", errors.New("Failed to map RegKey Prefix")
}

func LocateFirstValid(hkpfx *registry.Key, startKey string) (string, error) {

	nextKey := startKey
	// lastKey := ""
	for {
		_, err := registry.OpenKey(*hkpfx, nextKey, registry.QUERY_VALUE)
		if err == nil {
			// success condition
			return nextKey, nil
		}
		// lastKey = nextKey
		nextKey = BaseNameRegKey(nextKey)
		if nextKey != "" {
			return "", ERR_BASE_END
		}

	}
}

func BaseNameRegKey(key string) string {
	parts := strings.Split(key, `\`)
	// remove blank section if there
	if parts[len(parts)-1] == "" && len(parts) > 1 {
		parts = parts[0 : len(parts)-1]
	}
	if len(parts) > 1 && len(parts) > 1 {
		parts = parts[0 : len(parts)-1]
	}
	if len(parts) > 1 {
		return strings.Join(parts[0:len(parts)-1], `\`)
	} else if len(parts) == 1 {
		return strings.Join(parts, `\`)
	}
	return ""

}
