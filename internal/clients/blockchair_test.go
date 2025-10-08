package clients

import "testing"

func TestIsValidAddress(t *testing.T) {
	client := NewBlockchairClient()

	testCases := []struct {
		address string
		valid   bool
	}{
		{"bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5", true},  // Bech32
		{"3E8ociqZa9mZUSwGdSmAEMAoAxBK3FNDcd", true},           // P2SH
		{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", true},           // P2PKH
		{"invalid", false},                                      // Too short
		{"", false},                                             // Empty
		{"2N1234567890abcdef", false},                           // Wrong prefix
	}

	for _, tc := range testCases {
		result := client.IsValidAddress(tc.address)
		if result != tc.valid {
			t.Errorf("IsValidAddress(%s) = %v; want %v", tc.address, result, tc.valid)
		}
	}
}
