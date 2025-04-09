package alloha

import (
	"database/sql"
	"encoding/json"
)

// NullInt32 represents a nullable int32
type NullInt32 sql.NullInt32

// UnmarshalJSON implements the json.Unmarshaler interface
func (n *NullInt32) UnmarshalJSON(b []byte) error {
	var v *int32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if v != nil {
		n.Int32 = *v
		n.Valid = true
	} else {
		n.Valid = false
	}

	return nil
}

// FindOneResponse represents the structure of the API response to searching for data by ID
type FindOneResponse struct {
	// Request status ("success" or "error")
	Status string `json:"status"`
	// Error information (only for the "error" status)
	ErrorInfo string `json:"error_info"`
	// The movie or TV series data field
	Data *MovieData `json:"data"`
}

// FindListResponse represents a structure for processing the API response to data search
type FindListResponse struct {
	// Request status ("success" or "error")
	Status string `json:"status"`
	// Error information (only for the "error" status)
	ErrorInfo string `json:"error_info"`
	// The movie or TV series data field
	Data []*MovieSearchData `json:"data"`
	// Next page (optional)
	NextPage NullInt32 `json:"next_page"`
	// Previous page (optional)
	PrevPage NullInt32 `json:"prev_page"`
}

// ListOfLatestSeriesResponse represents a structure for handling an API response to a list of the latest episodes of a TV series
type ListOfLatestSeriesResponse struct {
	// Request status ("success" or "error")
	Status string `json:"status"`
	// Error information (only for the "error" status)
	ErrorInfo string `json:"error_info"`
	// TV series latest episode data field
	Data []*SeriesData `json:"data"`
	// Next page (optional)
	NextPage NullInt32 `json:"next_page"`
	// Previous page (optional)
	PrevPage NullInt32 `json:"prev_page"`
}

// MovieData represents the structure of information about a movie or TV series
type MovieData struct {
	Name                  string                       `json:"name"`
	OriginalName          string                       `json:"original_name"`
	AlternativeName       string                       `json:"alternative_name"`
	Year                  int                          `json:"year"`
	Category              int                          `json:"category"`
	IDKp                  int                          `json:"id_kp"`
	AlternativeIDKp       NullInt32                    `json:"alternative_id_kp"`
	IDImdb                string                       `json:"id_imdb"`
	IDTmdb                NullInt32                    `json:"id_tmdb"`
	IDWorldArt            NullInt32                    `json:"id_world_art"`
	TokenMovie            string                       `json:"token_movie"`
	Country               string                       `json:"country"`
	Genre                 string                       `json:"genre"`
	Actors                string                       `json:"actors"`
	Directors             string                       `json:"directors"`
	Producers             string                       `json:"producers"`
	PremiereRu            string                       `json:"premiere_ru"`
	Premiere              string                       `json:"premiere"`
	AgeRestrictions       NullInt32                    `json:"age_restrictions"`
	RatingMpaa            string                       `json:"rating_mpaa"`
	RatingKp              float64                      `json:"rating_kp"`
	RatingImdb            float64                      `json:"rating_imdb"`
	Time                  string                       `json:"time"`
	Tagline               string                       `json:"tagline"`
	Poster                string                       `json:"poster"`
	Description           string                       `json:"description"`
	SeasonsCount          int                          `json:"seasons_count"`
	Seasons               map[string]SeasonIframe      `json:"seasons"`
	Quality               string                       `json:"quality"`
	Translation           string                       `json:"translation"`
	TranslationIframe     map[string]TranslationIframe `json:"translation_iframe"`
	Iframe                string                       `json:"iframe"`
	IframeTrailer         string                       `json:"iframe_trailer"`
	Lgbt                  bool                         `json:"lgbt"`
	Uhd                   bool                         `json:"uhd"`
	AvailableDirectorsCut bool                         `json:"available_directors_cut"`
}

// MovieSearchData represents the structure of information about a movie or TV series
type MovieSearchData struct {
	LastSeason            NullInt32                    `json:"last_season"`
	LastEpisode           NullInt32                    `json:"last_episode"`
	Name                  string                       `json:"name"`
	OriginalName          string                       `json:"original_name"`
	AlternativeName       string                       `json:"alternative_name"`
	Year                  int                          `json:"year"`
	CategoryId            int                          `json:"category_id"`
	IDKp                  int                          `json:"id_kp"`
	AlternativeIDKp       NullInt32                    `json:"alternative_id_kp"`
	IDImdb                string                       `json:"id_imdb"`
	IDTmdb                NullInt32                    `json:"id_tmdb"`
	IDWorldArt            NullInt32                    `json:"id_world_art"`
	TokenMovie            string                       `json:"token_movie"`
	Country               string                       `json:"country"`
	Genre                 string                       `json:"genre"`
	Actors                string                       `json:"actors"`
	Directors             string                       `json:"directors"`
	Producers             string                       `json:"producers"`
	PremiereRu            string                       `json:"premiere_ru"`
	Premiere              string                       `json:"premiere"`
	AgeRestrictions       NullInt32                    `json:"age_restrictions"`
	RatingMpaa            string                       `json:"rating_mpaa"`
	RatingKp              float64                      `json:"rating_kp"`
	RatingImdb            float64                      `json:"rating_imdb"`
	Time                  string                       `json:"time"`
	Tagline               string                       `json:"tagline"`
	Poster                string                       `json:"poster"`
	Description           string                       `json:"description"`
	SeasonsCount          int                          `json:"seasons_count"`
	Seasons               map[string]SeasonIframe      `json:"seasons"`
	Quality               string                       `json:"quality"`
	Translation           string                       `json:"translation"`
	TranslationIframe     map[string]TranslationIframe `json:"translation_iframe"`
	Iframe                string                       `json:"iframe"`
	IframeTrailer         string                       `json:"iframe_trailer"`
	Lgbt                  bool                         `json:"lgbt"`
	Uhd                   bool                         `json:"uhd"`
	AvailableDirectorsCut bool                         `json:"available_directors_cut"`
}

// EpisodeIframe represents the structure of the episode iframe
type EpisodeIframe struct {
	Iframe      string                       `json:"iframe"`
	Episode     int                          `json:"episode"`
	Translation map[string]TranslationIframe `json:"translation"`
}

// SeasonIframe represents the structure of the season iframe
type SeasonIframe struct {
	Iframe   string                   `json:"iframe"`
	Season   int                      `json:"season"`
	Episodes map[string]EpisodeIframe `json:"episodes"`
}

// SeriesData represents the structure of information about TV series
type SeriesData struct {
	Season          int       `json:"season"`
	Episode         int       `json:"episode"`
	Translation     int       `json:"translation"`
	Quality         string    `json:"quality"`
	AdvPresence     int       `json:"adv_presence"`
	Name            string    `json:"name"`
	IDItem          int       `json:"id_item"`
	OriginalName    string    `json:"original_name"`
	Category        string    `json:"category"`
	AlternativeName string    `json:"alternative_name"`
	Year            int       `json:"year"`
	IDKp            int       `json:"id_kp"`
	AlternativeIDKp NullInt32 `json:"alternative_id_kp"`
	IDImdb          string    `json:"id_imdb"`
	IDTmdb          NullInt32 `json:"id_tmdb"`
	IDWorldArt      NullInt32 `json:"id_world_art"`
	TokenMovie      string    `json:"token_movie"`
	Date            string    `json:"date"`
	Iframe          string    `json:"iframe"`
	Adv             bool      `json:"adv"`
	CategoryId      int       `json:"category_id"`
	IframeLast      string    `json:"iframe_last"`
	IframeTrailer   string    `json:"iframe_trailer"`
	Lgbt            bool      `json:"lgbt"`
	Uhd             bool      `json:"uhd"`
}

// TranslationIframe represents the structure of the translation iframe
type TranslationIframe struct {
	Name    string `json:"name"`
	Iframe  string `json:"iframe"`
	Quality string `json:"quality"`
	Adv     bool   `json:"adv"`
	Date    string `json:"date"`
	Lgbt    bool   `json:"lgbt"`
	Uhd     bool   `json:"uhd"`
}
