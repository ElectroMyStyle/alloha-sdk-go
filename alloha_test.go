package alloha

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_buildApiURL(t *testing.T) {
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

func TestAPIClient_FindByIMDbId_EmptyIMDbIdParameter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByIMDbId(t.Context(), "")

	// Проверяем результат
	assert.Nil(t, movie)

	assert.ErrorIs(t, errMovie, EmptyIMDbIdParameterError)
}

func TestAPIClient_FindByIMDbId_StatusResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "tt-id", r.URL.Query().Get("imdb"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status:    "error",
			ErrorInfo: "not movie",
		}

		w.WriteHeader(http.StatusOK)
		errEncode := json.NewEncoder(w).Encode(movie)
		if errEncode != nil {
			t.Errorf("failed to encode movie response: %v", errEncode)
		}
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByIMDbId(t.Context(), "tt-id")

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "error", movie.Status)
	assert.Equal(t, "not movie", movie.ErrorInfo)
	assert.Nil(t, movie.Data)
}

func TestAPIClient_FindByIMDbId_StatusResponseSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "tt0110912", r.URL.Query().Get("imdb"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status: "success",
			Data: &MovieData{
				Name:   "Криминальное чтиво",
				IDKp:   342,
				IDImdb: "tt0110912",
				IDTmdb: 680,
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
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByIMDbId(t.Context(), "tt0110912")

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "success", movie.Status)
	assert.Empty(t, movie.ErrorInfo)

	assert.Equal(t, "Криминальное чтиво", movie.Data.Name)
	assert.Equal(t, 342, movie.Data.IDKp)
	assert.Equal(t, "tt0110912", movie.Data.IDImdb)
	assert.Equal(t, 680, movie.Data.IDTmdb)
}

func TestAPIClient_FindByKPId_InvalidKPIdParameter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByKPId(t.Context(), -1)

	// Проверяем результат
	assert.Nil(t, movie)

	assert.ErrorIs(t, errMovie, InvalidKPIdParameterError)
}

func TestAPIClient_FindByKPId_StatusResponseError(t *testing.T) {
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
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
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

func TestAPIClient_FindByKPId_StatusResponseSuccess(t *testing.T) {
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
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
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

func TestAPIClient_FindByTMDbId_InvalidTMDbIdParameter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByTMDbId(t.Context(), -1)

	// Проверяем результат
	assert.Nil(t, movie)

	assert.ErrorIs(t, errMovie, InvalidTMDbIdParameterError)
}

func TestAPIClient_FindByTMDbId_StatusResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "21029674", r.URL.Query().Get("tmdb"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status:    "error",
			ErrorInfo: "not movie",
		}

		w.WriteHeader(http.StatusOK)
		errEncode := json.NewEncoder(w).Encode(movie)
		if errEncode != nil {
			t.Errorf("failed to encode movie response: %v", errEncode)
		}
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByTMDbId(t.Context(), 21029674)

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "error", movie.Status)
	assert.Equal(t, "not movie", movie.ErrorInfo)
	assert.Nil(t, movie.Data)
}

func TestAPIClient_FindByTMDbId_StatusResponseSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "57532", r.URL.Query().Get("tmdb"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status: "success",
			Data: &MovieData{
				Name:   "Щенячий патруль",
				IDKp:   790343,
				IDTmdb: 57532,
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
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.FindByTMDbId(t.Context(), 57532)

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "success", movie.Status)
	assert.Empty(t, movie.ErrorInfo)

	assert.Equal(t, "Щенячий патруль", movie.Data.Name)
	assert.Equal(t, 790343, movie.Data.IDKp)
	assert.Equal(t, 57532, movie.Data.IDTmdb)
}

func TestAPIClient_SearchForOneByName_EmptyMovieNameParameter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.SearchForOneByName(t.Context(), "")

	// Проверяем результат
	assert.Nil(t, movie)

	assert.ErrorIs(t, errMovie, EmptyMovieNameParameterError)
}

func TestAPIClient_SearchForOneByName_StatusResponseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "Преступники", r.URL.Query().Get("name"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status:    "error",
			ErrorInfo: "not movie",
		}

		w.WriteHeader(http.StatusOK)
		errEncode := json.NewEncoder(w).Encode(movie)
		if errEncode != nil {
			t.Errorf("failed to encode movie response: %v", errEncode)
		}
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.SearchForOneByName(t.Context(), "Преступники")

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "error", movie.Status)
	assert.Equal(t, "not movie", movie.ErrorInfo)
	assert.Nil(t, movie.Data)
}

func TestAPIClient_SearchForOneByName_StatusResponseSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие конкретных параметров в URL
		assert.Equal(t, "Преступники", r.URL.Query().Get("name"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("token"))

		// Возвращаем тестовые данные
		movie := &FindOneResponse{
			Status: "success",
			Data: &MovieData{
				Name:   "Преступники",
				IDKp:   4859936,
				IDImdb: "tt14531774",
				IDTmdb: 201076,
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
	client, err := NewAPIClient(ts.Client(), "test-api-key", ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	movie, errMovie := client.SearchForOneByName(t.Context(), "Преступники")

	// Проверяем результат
	assert.Nil(t, errMovie)

	assert.Equal(t, "success", movie.Status)
	assert.Empty(t, movie.ErrorInfo)

	assert.Equal(t, "Преступники", movie.Data.Name)
	assert.Equal(t, 4859936, movie.Data.IDKp)
	assert.Equal(t, "tt14531774", movie.Data.IDImdb)
	assert.Equal(t, 201076, movie.Data.IDTmdb)
}
