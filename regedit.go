package fuzzyw32registry

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var ERR_BASE_END = errors.New("EOF")

var PrefixConstMap = map[string]registry.Key{
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

var ConstPrefixMap = map[registry.Key]string{
	registry.CLASSES_ROOT:     "CLASSES_ROOT",
	registry.CURRENT_USER:     "CURRENT_USER",
	registry.LOCAL_MACHINE:    "LOCAL_MACHINE",
	registry.USERS:            "USERS",
	registry.CURRENT_CONFIG:   "CURRENT_CONFIG",
	registry.PERFORMANCE_DATA: "PERFORMANCE_DATA",
}

func AddPrefix(pf *registry.Key, k string) (string, error) {
	if pstr, ok := ConstPrefixMap[*pf]; ok {
		return fmt.Sprintf(`HKEY_%s\%s`, pstr, k), nil
	}
	return "", errors.New("Prefix did not match any.")
}

func FuzzyRegKey(inkey string) (*registry.Key, string, error) {
	inkey = strings.Trim(inkey, "\n\r \t")
	inkey = strings.TrimPrefix(inkey, "Computer")
	prekey, key, err := SplitKeyPrefix(inkey)
	if err != nil {
		return nil, key, err
	}
	key, err = TransverseUntilRealKey(prekey, key)
	if err != nil {

		return nil, key, err
	}
	return prekey, key, nil
}

func SplitKeyPrefix(str string) (*registry.Key, string, error) {
	str = strings.TrimPrefix(str, `\`)
	parts := strings.Split(str, `\`)
	keystart, therestparts := parts[0], parts[1:]
	therest := strings.Join(therestparts, `\`)
	keystart = strings.ToUpper(keystart)

	ignore_start := 0
	for _, prefx := range []string{"HKEY", "HK"} {
		if len(keystart) <= len(prefx) {
			break
		}
		if strings.HasPrefix(keystart, prefx) {
			ignore_start = len(prefx)
			if string(keystart[len(prefx)]) == "_" {
				if len(keystart) > len(prefx)+1 {
					ignore_start += 1
				}
			}
			break
		}
	}
	reduced_keystart := keystart[ignore_start:]

	for k, v := range PrefixConstMap {
		if strings.HasPrefix(reduced_keystart, k) {
			return &v, therest, nil
		}
	}
	return nil, "", errors.New("Failed to map Registry Prefix \"" + keystart + `"`)
}

func TransverseUntilRealKey(hkpfx *registry.Key, startKey string) (string, error) {
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
