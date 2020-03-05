package registry

import (
	"testing"
)

func TestBelongs(t *testing.T) {
	testCases := []struct {
		name     string
		cli      *Client
		input    string
		expected bool
	}{
		{
			name:     "Nominal docker",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "docker.io/alebedev87/coredns:1.3.1",
			expected: true,
		},
		{
			name:     "Nominal quay",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "quay.io/alebedev87/coredns:1.3.1",
			expected: true,
		},
		{
			name:     "Spaces removed",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "  quay.io/alebedev87/coredns:1.3.1  ",
			expected: true,
		},
		{
			name:     "Different registry",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "quay.io/kubermatic/openvpn:v0.5",
			expected: false,
		},
		{
			name:     "Different registry",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "docker.io/coredns/coredns:1.3.1",
			expected: false,
		},
		{
			name:     "Different organization",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "docker.io/coredns/coredns:1.3.1",
			expected: false,
		},
		{
			name:     "Rubbish",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "docker.io",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.cli.Belongs(tc.input)
			if output != tc.expected {
				t.Errorf("Test case %q: output didn't match", tc.name)
			}
		})
	}
}

func TestNewFullName(t *testing.T) {
	testCases := []struct {
		name     string
		cli      *Client
		input    string
		expected string
	}{
		{
			name:     "Nominal all different",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "quay.io/kubermatic/openvpn:v0.5",
			expected: "docker.io/alebedev87/quay.io-kubermatic-openvpn:v0.5",
		},
		{
			name:     "Nominal organization different",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "docker.io/coredns/coredns:1.3.1",
			expected: "docker.io/alebedev87/docker.io-coredns-coredns:1.3.1",
		},
		{
			name:     "Nominal all different 2",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "docker.io/coredns/coredns:1.3.1",
			expected: "quay.io/alebedev87/docker.io-coredns-coredns:1.3.1",
		},
		{
			name:     "Nominal organization different 2",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "quay.io/coredns/coredns:1.3.1",
			expected: "quay.io/alebedev87/quay.io-coredns-coredns:1.3.1",
		},
		{
			name:     "Nominal no compacting",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "openvpn:v0.5",
			expected: "docker.io/alebedev87/openvpn:v0.5",
		},
		{
			name:     "Nominal no compacting no tag",
			cli:      NewClient("registry-1.docker.io", "alebedev87", "", "", 0),
			input:    "openvpn",
			expected: "docker.io/alebedev87/openvpn",
		},
		{
			name:     "Spaces removed",
			cli:      NewClient("quay.io", "alebedev87", "", "", 0),
			input:    "  quay.io/coredns:1.3.1  ",
			expected: "quay.io/alebedev87/quay.io-coredns:1.3.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.cli.newFullName(tc.input)

			if output != tc.expected {
				t.Errorf("Output didn't match. Expected: %q, got: %q", tc.expected, output)
			}
		})
	}
}

func TestSkopeoCopyCmd(t *testing.T) {
	testCases := []struct {
		name     string
		cli      *Client
		src      string
		dst      string
		expected string
	}{
		{
			name:     "Nominal",
			cli:      NewClient("", "", "here", "there", 0),
			src:      "quay.io/coredns:1.3.1",
			dst:      "docker.io/alebedev87/coredns:1.3.1",
			expected: "skopeo copy --dest-creds here:there docker://quay.io/coredns:1.3.1 docker://docker.io/alebedev87/coredns:1.3.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.cli.skopeoCopyCmd(tc.src, tc.dst)
			if output != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, output)
			}
		})
	}
}
