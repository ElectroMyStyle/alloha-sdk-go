package alloha

// FindOneResponse represents the structure of the API response to searching for data by ID
type FindOneResponse struct {
	// Request status ("success" or "error")
	Status string `json:"status"`
	// Error information (only for the "error" status)
	ErrorInfo string `json:"error_info,omitempty"`
	// The movie or TV series data field
	Data *MovieData `json:"data,omitempty"`
}

// FindManyResponse represents a structure for processing the API response to data search
type FindManyResponse struct {
	// Request status ("success" or "error")
	Status string `json:"status"`
	// Error information (only for the "error" status)
	ErrorInfo string `json:"error_info,omitempty"`
	// The movie or TV series data field
	Data []*MovieData `json:"data,omitempty"`
	// Next page (optional)
	NextPage *int `json:"next_page,omitempty"`
	// Previous page (optional)
	PrevPage *int `json:"prev_page,omitempty"`
}

// MovieData represents the structure of information about a movie or TV series
type MovieData struct {
	Name                  string                       `json:"name"`
	OriginalName          string                       `json:"original_name"`
	AlternativeName       string                       `json:"alternative_name"`
	Year                  int                          `json:"year"`
	Category              int                          `json:"category"`
	IDKp                  int                          `json:"id_kp"`
	AlternativeIDKp       int                          `json:"alternative_id_kp"`
	IDImdb                string                       `json:"id_imdb"`
	IDTmdb                int                          `json:"id_tmdb"`
	IDWorldArt            int                          `json:"id_world_art"`
	TokenMovie            string                       `json:"token_movie"`
	Country               string                       `json:"country"`
	Genre                 string                       `json:"genre"`
	Actors                string                       `json:"actors"`
	Directors             string                       `json:"directors"`
	Producers             string                       `json:"producers"`
	PremiereRu            string                       `json:"premiere_ru"`
	Premiere              string                       `json:"premiere"`
	AgeRestrictions       int                          `json:"age_restrictions"`
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

// TranslationIframe represents the structure of the translation iframe
type TranslationIframe struct {
	Name    string `json:"name"`
	Iframe  string `json:"iframe"`
	Quality string `json:"quality"`
	Adv     bool   `json:"adv"`
	Date    string `json:"date"`
}
