package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromJSON(t *testing.T) {
	cases := []struct {
		name                 string
		specifiedJSON        string
		expectedErrorMessage string
		expected             []cryptoBalanceCheckerConfig
	}{
		{"case #1", `[{"symbol": "BTC", "addresses": ["a"]},{"symbol": "DASH","addresses": ["b","c"],"api_key": "apikey1"},{"symbol": "ETH","addresses": ["d"],"api_key": "apikey2"}]`,
			"",
			[]cryptoBalanceCheckerConfig{
				{btc, []string{"a"}, ""},
				{dash, []string{"b", "c"}, "apikey1"},
				{eth, []string{"d"}, "apikey2"},
			},
		},
		{"case #2", `[{"symbol": "UNO", "addresses": ["asdkfhjkadfghds"]}]`,
			"",
			[]cryptoBalanceCheckerConfig{
				{uno, []string{"asdkfhjkadfghds"}, ""},
			},
		},
		{"case #3 (invalid JSON)", `[{"symbol": "UNO", "addresses": ["asdkfhjkadfghds",]}]`,
			"invalid character ']' looking for beginning of value",
			[]cryptoBalanceCheckerConfig{},
		},
	}

	for _, testCase := range cases {
		testCaseName := fmt.Sprintf("Test case %s", testCase.name)
		config, err := loadConfigFromJSON([]byte(testCase.specifiedJSON))

		if testCase.expectedErrorMessage != "" {
			require.Error(t, err, testCaseName)
			require.Equal(t, testCase.expectedErrorMessage, err.Error(), testCaseName)
		} else {
			require.NotNil(t, config, testCaseName)
			require.Len(t, config, len(testCase.expected), "%s: expected %d crypto-currencies", testCaseName, len(testCase.expected))
		}

		for idx, expected := range testCase.expected {
			require.Equal(t, expected.Symbol, config[idx].Symbol, testCaseName)
			require.Equal(t, expected.APIKey, config[idx].APIKey, testCaseName)
			require.EqualValues(t, expected.Addresses, config[idx].Addresses, testCaseName)
		}
	}
}
