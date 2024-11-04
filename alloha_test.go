package alloha

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildApiURL(t *testing.T) {
	tests := []struct {
		name        string
		apiToken    string
		baseApiUrl  string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "empty API token",
			apiToken:    "",
			baseApiUrl:  "https://example.com",
			expectedErr: ApiTokenEmptyError,
		},
		{
			name:        "empty base API URL",
			apiToken:    "testToken",
			baseApiUrl:  "",
			expectedErr: BaseApiUrlEmptyError,
		},
		{
			name:        "invalid base API URL scheme",
			apiToken:    "testToken",
			baseApiUrl:  "://example.com",
			expectedErr: errors.New("missing protocol scheme"),
		},
		{
			name:        "invalid base API URL host",
			apiToken:    "testToken",
			baseApiUrl:  "http://",
			expectedErr: BaseApiUrlInvalidHostError,
		},
		{
			name:        "valid input",
			apiToken:    "testToken",
			baseApiUrl:  "https://example.com/api",
			expectedURL: "https://example.com/?token=testToken",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotErr := buildApiURL(tt.apiToken, tt.baseApiUrl)

			if tt.expectedErr != nil {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.expectedURL, gotURL)
			}
		})
	}
}
