package alloha

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllohaClient_buildApiURL(t *testing.T) {
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

func TestAllohaClient_FindByKPId_InvalidKPIdParameter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAllohaClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByKPId(t.Context(), -1)

	// Проверяем результат
	assert.Nil(t, movie)

	assert.ErrorIs(t, errMovie, InvalidKPIdParameterError)
}

func TestAllohaClient_FindByKPId_StatusResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "5119525", r.URL.Query().Get("kp"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status:    "error",
			ErrorInfo: "not valid token",
		}

		w.WriteHeader(http.StatusOK)
		errEncode := json.NewEncoder(w).Encode(movie)
		if errEncode != nil {
			t.Errorf("failed to encode movie response: %v", errEncode)
		}
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAllohaClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByKPId(t.Context(), 5119525)

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "error", movie.Status)
	assert.Equal(t, "not valid token", movie.ErrorInfo)
	assert.Nil(t, movie.Data)
}

func TestAllohaClient_FindByKPId_StatusResponseSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "5119525", r.URL.Query().Get("kp"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status: "success",
			Data: &MovieData{
				Name: "Преступники",
				IDKp: 5119525,
			},
		}

		w.WriteHeader(http.StatusOK)
		errEncode := json.NewEncoder(w).Encode(movie)
		if errEncode != nil {
			t.Errorf("failed to encode movie response: %v", errEncode)
		}
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAllohaClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByKPId(t.Context(), 5119525)

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "success", movie.Status)
	assert.Empty(t, movie.ErrorInfo)

	assert.Equal(t, "Преступники", movie.Data.Name)
	assert.Equal(t, 5119525, movie.Data.IDKp)
}
