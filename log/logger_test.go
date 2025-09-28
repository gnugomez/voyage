package log

import "testing"

func TestParseLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected LogLevel
	}{
		{
			name:     "Debug",
			input:    "debug",
			expected: DebugLevel,
		},
		{
			name:     "Info",
			input:    "info",
			expected: InfoLevel,
		},
		{
			name:     "Error",
			input:    "error",
			expected: ErrorLevel,
		},
		{
			name:     "Fatal",
			input:    "fatal",
			expected: FatalLevel,
		},
		{
			name:     "Default to Info",
			input:    "some-unknown-level",
			expected: InfoLevel,
		},
		{
			name:     "Empty string defaults to Info",
			input:    "",
			expected: InfoLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := ParseLogLevel(tc.input); result != tc.expected {
				t.Errorf("For input '%s', expected log level %s, but got %s", tc.input, tc.expected, result)
			}
		})
	}
}
