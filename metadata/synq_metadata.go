package metadata

type MetaData struct {
  Version         string              `json:"metadata_version"`
  Title           LanguageList        `json:"title"`
  Description     LanguageList        `json:"description"`
  Year            int                 `json:"production_year"`
  Type            string              `json:"type"`
  Series          Series              `json:"series,omitempty"`
  Genres          []string            `json:"genres"`
  Credits         []Credit            `json:"credits"`
  Regional        bool                `json:"regional_content"`
  Rating          string              `json:"parental_rating"`
  Ratio           string              `json:"aspect_ratio,omitempty"`
  Duration        string              `json:"expected_duration,omitempty"`
  Countries       []string            `json:"country_of_origin"`
}

type Series struct {
  Episode         int                 `json:"episode_number,omitempty"`
  Season          int                 `json:"season,omitempty"`
  ExternalId      string              `json:"external_id,omitempty"`
  InternalId      string              `json:"internal_id,omitempty"`
  EpisodeCount    int                 `json:"episodes_in_season,omitempty"`
}

type Credit struct {
  Name            string              `json:"name"`
  Function        string              `json:"role"`
}

type Language     map[string]string
type LanguageList map[string]Language