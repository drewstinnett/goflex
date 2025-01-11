package plexrando

import "encoding/xml"

// SeasonsResponse is the response from a seasons request
type SeasonsResponse struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Text                string   `xml:",chardata"`
	Size                string   `xml:"size,attr"`
	AllowSync           string   `xml:"allowSync,attr"`
	Art                 string   `xml:"art,attr"`
	Identifier          string   `xml:"identifier,attr"`
	Key                 string   `xml:"key,attr"`
	LibrarySectionID    string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string   `xml:"mediaTagVersion,attr"`
	Nocache             string   `xml:"nocache,attr"`
	ParentIndex         string   `xml:"parentIndex,attr"`
	ParentTitle         string   `xml:"parentTitle,attr"`
	ParentYear          string   `xml:"parentYear,attr"`
	Summary             string   `xml:"summary,attr"`
	Theme               string   `xml:"theme,attr"`
	Thumb               string   `xml:"thumb,attr"`
	Title1              string   `xml:"title1,attr"`
	Title2              string   `xml:"title2,attr"`
	ViewGroup           string   `xml:"viewGroup,attr"`
	Directory           []struct {
		Text            string `xml:",chardata"`
		LeafCount       string `xml:"leafCount,attr"`
		Thumb           string `xml:"thumb,attr"`
		ViewedLeafCount string `xml:"viewedLeafCount,attr"`
		Key             string `xml:"key,attr"`
		Title           string `xml:"title,attr"`
		RatingKey       string `xml:"ratingKey,attr"`
		ParentRatingKey string `xml:"parentRatingKey,attr"`
		GUID            string `xml:"guid,attr"`
		ParentGUID      string `xml:"parentGuid,attr"`
		ParentSlug      string `xml:"parentSlug,attr"`
		ParentStudio    string `xml:"parentStudio,attr"`
		Type            string `xml:"type,attr"`
		ParentKey       string `xml:"parentKey,attr"`
		ParentTitle     string `xml:"parentTitle,attr"`
		Summary         string `xml:"summary,attr"`
		Index           string `xml:"index,attr"`
		ParentIndex     string `xml:"parentIndex,attr"`
		ViewCount       string `xml:"viewCount,attr"`
		LastViewedAt    string `xml:"lastViewedAt,attr"`
		ParentYear      string `xml:"parentYear,attr"`
		Art             string `xml:"art,attr"`
		ParentThumb     string `xml:"parentThumb,attr"`
		ParentTheme     string `xml:"parentTheme,attr"`
		AddedAt         string `xml:"addedAt,attr"`
		UpdatedAt       string `xml:"updatedAt,attr"`
		SkipCount       string `xml:"skipCount,attr"`
		Image           []struct {
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
	} `xml:"Directory"`
}

// ShowsResponse returns a response for a list of TV shows
type ShowsResponse struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Text                string   `xml:",chardata"`
	Size                string   `xml:"size,attr"`
	AllowSync           string   `xml:"allowSync,attr"`
	Art                 string   `xml:"art,attr"`
	Content             string   `xml:"content,attr"`
	Identifier          string   `xml:"identifier,attr"`
	LibrarySectionID    string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string   `xml:"mediaTagVersion,attr"`
	Nocache             string   `xml:"nocache,attr"`
	Thumb               string   `xml:"thumb,attr"`
	Title1              string   `xml:"title1,attr"`
	Title2              string   `xml:"title2,attr"`
	ViewGroup           string   `xml:"viewGroup,attr"`
	Directory           []struct {
		Text                   string `xml:",chardata"`
		RatingKey              string `xml:"ratingKey,attr"`
		Key                    string `xml:"key,attr"`
		GUID                   string `xml:"guid,attr"`
		Slug                   string `xml:"slug,attr"`
		Studio                 string `xml:"studio,attr"`
		Type                   string `xml:"type,attr"`
		Title                  string `xml:"title,attr"`
		ContentRating          string `xml:"contentRating,attr"`
		Summary                string `xml:"summary,attr"`
		Index                  string `xml:"index,attr"`
		AudienceRating         string `xml:"audienceRating,attr"`
		ViewCount              string `xml:"viewCount,attr"`
		SkipCount              string `xml:"skipCount,attr"`
		LastViewedAt           string `xml:"lastViewedAt,attr"`
		Year                   string `xml:"year,attr"`
		Tagline                string `xml:"tagline,attr"`
		Thumb                  string `xml:"thumb,attr"`
		Art                    string `xml:"art,attr"`
		Theme                  string `xml:"theme,attr"`
		Duration               string `xml:"duration,attr"`
		OriginallyAvailableAt  string `xml:"originallyAvailableAt,attr"`
		LeafCount              string `xml:"leafCount,attr"`
		ViewedLeafCount        string `xml:"viewedLeafCount,attr"`
		ChildCount             string `xml:"childCount,attr"`
		AddedAt                string `xml:"addedAt,attr"`
		UpdatedAt              string `xml:"updatedAt,attr"`
		AudienceRatingImage    string `xml:"audienceRatingImage,attr"`
		HasPremiumPrimaryExtra string `xml:"hasPremiumPrimaryExtra,attr"`
		PrimaryExtraKey        string `xml:"primaryExtraKey,attr"`
		HasPremiumExtras       string `xml:"hasPremiumExtras,attr"`
		SeasonCount            string `xml:"seasonCount,attr"`
		TitleSort              string `xml:"titleSort,attr"`
		Rating                 string `xml:"rating,attr"`
		Banner                 string `xml:"banner,attr"`
		Image                  []struct {
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
		Genre []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Genre"`
		Country []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Country"`
		Role []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Role"`
	} `xml:"Directory"`
}

// LibraryResponse is what we get back when listing the libraries
type LibraryResponse struct {
	XMLName   xml.Name `xml:"MediaContainer"`
	Text      string   `xml:",chardata"`
	Size      string   `xml:"size,attr"`
	AllowSync string   `xml:"allowSync,attr"`
	Title1    string   `xml:"title1,attr"`
	Directory []struct {
		Text             string `xml:",chardata"`
		AllowSync        string `xml:"allowSync,attr"`
		Art              string `xml:"art,attr"`
		Composite        string `xml:"composite,attr"`
		Filters          string `xml:"filters,attr"`
		Refreshing       string `xml:"refreshing,attr"`
		Thumb            string `xml:"thumb,attr"`
		Key              string `xml:"key,attr"`
		Type             string `xml:"type,attr"`
		Title            string `xml:"title,attr"`
		Agent            string `xml:"agent,attr"`
		Scanner          string `xml:"scanner,attr"`
		Language         string `xml:"language,attr"`
		UUID             string `xml:"uuid,attr"`
		UpdatedAt        string `xml:"updatedAt,attr"`
		CreatedAt        string `xml:"createdAt,attr"`
		ScannedAt        string `xml:"scannedAt,attr"`
		Content          string `xml:"content,attr"`
		Directory        string `xml:"directory,attr"`
		ContentChangedAt string `xml:"contentChangedAt,attr"`
		Hidden           string `xml:"hidden,attr"`
		Co               string `xml:"co,attr"`
		NtentChangedAt   string `xml:"ntentChangedAt,attr"`
		A                string `xml:"a,attr"`
		Gent             string `xml:"gent,attr"`
		Cont             string `xml:"cont,attr"`
		Ent              string `xml:"ent,attr"`
		Location         []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Path string `xml:"path,attr"`
		} `xml:"Location"`
	} `xml:"Directory"`
}

// PlaylistResponse is the individual playlist response
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

// EpisodesResponse is is the response for episode listings
type EpisodesResponse struct {
	XMLName                  xml.Name `xml:"MediaContainer"`
	Text                     string   `xml:",chardata"`
	Size                     string   `xml:"size,attr"`
	AllowSync                string   `xml:"allowSync,attr"`
	Art                      string   `xml:"art,attr"`
	GrandparentContentRating string   `xml:"grandparentContentRating,attr"`
	GrandparentRatingKey     string   `xml:"grandparentRatingKey,attr"`
	GrandparentStudio        string   `xml:"grandparentStudio,attr"`
	GrandparentTheme         string   `xml:"grandparentTheme,attr"`
	GrandparentThumb         string   `xml:"grandparentThumb,attr"`
	GrandparentTitle         string   `xml:"grandparentTitle,attr"`
	Identifier               string   `xml:"identifier,attr"`
	Key                      string   `xml:"key,attr"`
	LibrarySectionID         string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle      string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID       string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix           string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion          string   `xml:"mediaTagVersion,attr"`
	Nocache                  string   `xml:"nocache,attr"`
	ParentIndex              string   `xml:"parentIndex,attr"`
	ParentTitle              string   `xml:"parentTitle,attr"`
	Theme                    string   `xml:"theme,attr"`
	Thumb                    string   `xml:"thumb,attr"`
	Title1                   string   `xml:"title1,attr"`
	Title2                   string   `xml:"title2,attr"`
	ViewGroup                string   `xml:"viewGroup,attr"`
	Video                    []struct {
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
		TitleSort             string `xml:"titleSort,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		ParentTitle           string `xml:"parentTitle,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Summary               string `xml:"summary,attr"`
		Index                 string `xml:"index,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		ViewCount             string `xml:"viewCount,attr"`
		LastViewedAt          string `xml:"lastViewedAt,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Art                   string `xml:"art,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
		GrandparentTheme      string `xml:"grandparentTheme,attr"`
		Duration              string `xml:"duration,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		AddedAt               string `xml:"addedAt,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		Media                 struct {
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
		Writer struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Writer"`
		Director struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Director"`
		Role []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Role"`
	} `xml:"Video"`
}

// PlaylistsResponse lists all playlists
type PlaylistsResponse struct {
	XMLName  xml.Name `xml:"MediaContainer"`
	Text     string   `xml:",chardata"`
	Size     string   `xml:"size,attr"`
	Playlist []struct {
		Text         string `xml:",chardata"`
		RatingKey    string `xml:"ratingKey,attr"`
		Key          string `xml:"key,attr"`
		GUID         string `xml:"guid,attr"`
		Type         string `xml:"type,attr"`
		Title        string `xml:"title,attr"`
		Summary      string `xml:"summary,attr"`
		Smart        string `xml:"smart,attr"`
		PlaylistType string `xml:"playlistType,attr"`
		Composite    string `xml:"composite,attr"`
		ViewCount    string `xml:"viewCount,attr"`
		LastViewedAt string `xml:"lastViewedAt,attr"`
		Duration     string `xml:"duration,attr"`
		LeafCount    string `xml:"leafCount,attr"`
		AddedAt      string `xml:"addedAt,attr"`
		UpdatedAt    string `xml:"updatedAt,attr"`
		TitleSort    string `xml:"titleSort,attr"`
	} `xml:"Playlist"`
}

// CreatePlaylistResponse is what we get back from a playlist create
type CreatePlaylistResponse struct {
	XMLName  xml.Name `xml:"MediaContainer"`
	Text     string   `xml:",chardata"`
	Size     string   `xml:"size,attr"`
	Playlist struct {
		Text         string `xml:",chardata"`
		RatingKey    string `xml:"ratingKey,attr"`
		Key          string `xml:"key,attr"`
		GUID         string `xml:"guid,attr"`
		Type         string `xml:"type,attr"`
		Title        string `xml:"title,attr"`
		Summary      string `xml:"summary,attr"`
		Smart        string `xml:"smart,attr"`
		PlaylistType string `xml:"playlistType,attr"`
		LeafCount    string `xml:"leafCount,attr"`
		AddedAt      string `xml:"addedAt,attr"`
		UpdatedAt    string `xml:"updatedAt,attr"`
	} `xml:"Playlist"`
}
