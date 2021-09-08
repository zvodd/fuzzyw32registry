package fuzzyw32registry

import (
	"os"
	"strconv"
	"testing"
)

func pad(x int) string {
	r := ""
	if x < 1 {
		return r
	}
	for ; x > 0; x-- {
		r += " "
	}
	return r
}

func TestFuzzyRegKey(t *testing.T) {
	tests := []string{
		`Computer\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HKEYLOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HKEY_LM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HKEYLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HK_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HKLOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		`HK_LM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`,

		// `Computer\HK_CU\Printers\Settings\Wizard`,
		// `Computer\HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\NlaSvc\Parameters\Internet`
	}

	for i, test := range tests {
		_, res, err := FuzzyRegKey(test)
		t.Logf("Test %d -> \"%s\"\n %s %s", i, test, pad(len(strconv.FormatInt(int64(i), 10))), res)
		if err != nil {
			t.Logf("Failed with err:\n%s\n-----------\n", err)
		}
	}
}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
