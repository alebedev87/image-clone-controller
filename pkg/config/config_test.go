package config

import (
	"os"
	"reflect"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name          string
		input         *Config
		userVar       bool
		pwdVar        bool
		expectedError bool
	}{
		{
			name:          "No empty values",
			input:         newTestConfig("1", "1", "1", "1"),
			expectedError: false,
		},
		{
			name:          "No registry",
			input:         newTestConfig("", "1", "1", "1"),
			expectedError: true,
		},
		{
			name:          "No organization",
			input:         newTestConfig("1", " ", "1", "1"),
			expectedError: true,
		},
		{
			name:          "No user",
			input:         newTestConfig("1", "1", "", "1"),
			expectedError: true,
		},
		{
			name:          "No password",
			input:         newTestConfig("1", "1", "1", "  "),
			expectedError: true,
		},
		{
			name:          "User var",
			input:         newTestConfig("1", "1", "", "1"),
			userVar:       true,
			expectedError: false,
		},
		{
			name:          "Password var",
			input:         newTestConfig("1", "1", "1", "  "),
			pwdVar:        true,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.userVar {
				os.Setenv(usernameVar, "1")
			} else {
				os.Unsetenv(usernameVar)
			}
			if tc.pwdVar {
				os.Setenv(passwordVar, "1")
			} else {
				os.Unsetenv(passwordVar)
			}
			err := tc.input.Validate()
			if err != nil {
				if !tc.expectedError {
					t.Error("Got not expected error")
				}
			} else {
				if tc.expectedError {
					t.Error("Got no error while one is expected")
				}
			}
		})
	}
}

func TestNamespaceBlacklist(t *testing.T) {
	testCases := []struct {
		name     string
		config   *Config
		expected map[string]bool
	}{
		{
			name: "Nominal",
			config: &Config{
				MandatoryNamespaceBlacklist: []string{"must"},
			},
			expected: map[string]bool{
				"must": true,
			},
		},
		{
			name: "Nominal with additional",
			config: &Config{
				MandatoryNamespaceBlacklist:  []string{"must"},
				AdditionalNamespaceBlacklist: []string{"should"},
			},
			expected: map[string]bool{
				"must":   true,
				"should": true,
			},
		},
		{
			name: "Duplication",
			config: &Config{
				MandatoryNamespaceBlacklist:  []string{"must"},
				AdditionalNamespaceBlacklist: []string{"must", "should", "shouldtoo", "should"},
			},
			expected: map[string]bool{
				"must":      true,
				"should":    true,
				"shouldtoo": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.config.NamespaceBlacklist()
			if !reflect.DeepEqual(output, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, output)
			}
		})
	}
}

func newTestConfig(reg, org, usr, pwd string) *Config {
	return &Config{
		Registry:                     reg,
		Organization:                 org,
		Username:                     usr,
		Password:                     pwd,
		ImageCopyTimeoutSeconds:      0,
		MandatoryNamespaceBlacklist:  []string{},
		AdditionalNamespaceBlacklist: []string{},
	}
}
