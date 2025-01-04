package plexrando

import "encoding/xml"

// PlaylistResponse is what we get back for a playlist from the API
type PlaylistResponse struct {
	XMLName      xml.Name `xml:"MediaContainer"`
	Text         string   `xml:",chardata"`
	Size         string   `xml:"size,attr"`
	Composite    string   `xml:"composite,attr"`
	Duration     string   `xml:"duration,attr"`
	LeafCount    string   `xml:"leafCount,attr"`
	PlaylistType string   `xml:"playlistType,attr"`
	RatingKey    string   `xml:"ratingKey,attr"`
	Smart        string   `xml:"smart,attr"`
	Title        string   `xml:"title,attr"`
	Video        []struct {
		Text                  string `xml:",chardata"`
		RatingKey             string `xml:"ratingKey,attr"`
		Key                   string `xml:"key,attr"`
		ParentRatingKey       string `xml:"parentRatingKey,attr"`
		GrandparentRatingKey  string `xml:"grandparentRatingKey,attr"`
		GUID                  string `xml:"guid,attr"`
		ParentGUID            string `xml:"parentGuid,attr"`
		GrandparentGUID       string `xml:"grandparentGuid,attr"`
		GrandparentSlug       string `xml:"grandparentSlug,attr"`
		Type                  string `xml:"type,attr"`
		Title                 string `xml:"title,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		LibrarySectionTitle   string `xml:"librarySectionTitle,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		LibrarySectionKey     string `xml:"librarySectionKey,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		ParentTitle           string `xml:"parentTitle,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Summary               string `xml:"summary,attr"`
		Index                 string `xml:"index,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		ViewOffset            string `xml:"viewOffset,attr"`
		LastViewedAt          string `xml:"lastViewedAt,attr"`
		Year                  string `xml:"year,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Art                   string `xml:"art,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
		GrandparentTheme      string `xml:"grandparentTheme,attr"`
		PlaylistItemID        string `xml:"playlistItemID,attr"`
		Duration              string `xml:"duration,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		AddedAt               string `xml:"addedAt,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		TitleSort             string `xml:"titleSort,attr"`
		AudienceRating        string `xml:"audienceRating,attr"`
		AudienceRatingImage   string `xml:"audienceRatingImage,attr"`
		ViewCount             string `xml:"viewCount,attr"`
		SkipCount             string `xml:"skipCount,attr"`
		Media                 []struct {
			Text            string `xml:",chardata"`
			ID              string `xml:"id,attr"`
			Duration        string `xml:"duration,attr"`
			Bitrate         string `xml:"bitrate,attr"`
			Width           string `xml:"width,attr"`
			Height          string `xml:"height,attr"`
			AspectRatio     string `xml:"aspectRatio,attr"`
			AudioChannels   string `xml:"audioChannels,attr"`
			AudioCodec      string `xml:"audioCodec,attr"`
			VideoCodec      string `xml:"videoCodec,attr"`
			VideoResolution string `xml:"videoResolution,attr"`
			Container       string `xml:"container,attr"`
			VideoFrameRate  string `xml:"videoFrameRate,attr"`
			VideoProfile    string `xml:"videoProfile,attr"`
			Part            struct {
				Text         string `xml:",chardata"`
				ID           string `xml:"id,attr"`
				Key          string `xml:"key,attr"`
				Duration     string `xml:"duration,attr"`
				File         string `xml:"file,attr"`
				Size         string `xml:"size,attr"`
				Container    string `xml:"container,attr"`
				VideoProfile string `xml:"videoProfile,attr"`
			} `xml:"Part"`
		} `xml:"Media"`
		Image []struct {
			Text string `xml:",chardata"`
			Alt  string `xml:"alt,attr"`
			Type string `xml:"type,attr"`
			URL  string `xml:"url,attr"`
		} `xml:"Image"`
		UltraBlurColors struct {
			Text        string `xml:",chardata"`
			TopLeft     string `xml:"topLeft,attr"`
			TopRight    string `xml:"topRight,attr"`
			BottomRight string `xml:"bottomRight,attr"`
			BottomLeft  string `xml:"bottomLeft,attr"`
		} `xml:"UltraBlurColors"`
		Role []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Role"`
		Director []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Director"`
		Writer struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Writer"`
	} `xml:"Video"`
}
