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

// ActiveSessionsResponse is the active sessions
type ActiveSessionsResponse struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Text    string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
	Video   []struct {
		Text                  string `xml:",chardata"`
		AddedAt               string `xml:"addedAt,attr"`
		Art                   string `xml:"art,attr"`
		AudienceRating        string `xml:"audienceRating,attr"`
		AudienceRatingImage   string `xml:"audienceRatingImage,attr"`
		ChapterSource         string `xml:"chapterSource,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Duration              string `xml:"duration,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
		GrandparentGUID       string `xml:"grandparentGuid,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		GrandparentRatingKey  string `xml:"grandparentRatingKey,attr"`
		GrandparentSlug       string `xml:"grandparentSlug,attr"`
		GrandparentTheme      string `xml:"grandparentTheme,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		GUID                  string `xml:"guid,attr"`
		Index                 string `xml:"index,attr"`
		Key                   string `xml:"key,attr"`
		LastViewedAt          string `xml:"lastViewedAt,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		LibrarySectionKey     string `xml:"librarySectionKey,attr"`
		LibrarySectionTitle   string `xml:"librarySectionTitle,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		ParentGUID            string `xml:"parentGuid,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		ParentRatingKey       string `xml:"parentRatingKey,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		ParentTitle           string `xml:"parentTitle,attr"`
		RatingKey             string `xml:"ratingKey,attr"`
		SessionKey            string `xml:"sessionKey,attr"`
		Summary               string `xml:"summary,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Title                 string `xml:"title,attr"`
		Type                  string `xml:"type,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		ViewCount             string `xml:"viewCount,attr"`
		ViewOffset            string `xml:"viewOffset,attr"`
		Year                  string `xml:"year,attr"`
		Media                 struct {
			Text                  string `xml:",chardata"`
			AudioProfile          string `xml:"audioProfile,attr"`
			ID                    string `xml:"id,attr"`
			VideoProfile          string `xml:"videoProfile,attr"`
			AudioChannels         string `xml:"audioChannels,attr"`
			AudioCodec            string `xml:"audioCodec,attr"`
			Bitrate               string `xml:"bitrate,attr"`
			Container             string `xml:"container,attr"`
			Duration              string `xml:"duration,attr"`
			Height                string `xml:"height,attr"`
			OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
			Protocol              string `xml:"protocol,attr"`
			VideoCodec            string `xml:"videoCodec,attr"`
			VideoFrameRate        string `xml:"videoFrameRate,attr"`
			VideoResolution       string `xml:"videoResolution,attr"`
			Width                 string `xml:"width,attr"`
			Selected              string `xml:"selected,attr"`
			Part                  struct {
				Text                  string `xml:",chardata"`
				AudioProfile          string `xml:"audioProfile,attr"`
				HasThumbnail          string `xml:"hasThumbnail,attr"`
				ID                    string `xml:"id,attr"`
				VideoProfile          string `xml:"videoProfile,attr"`
				Bitrate               string `xml:"bitrate,attr"`
				Container             string `xml:"container,attr"`
				Duration              string `xml:"duration,attr"`
				Height                string `xml:"height,attr"`
				OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
				Protocol              string `xml:"protocol,attr"`
				Width                 string `xml:"width,attr"`
				Decision              string `xml:"decision,attr"`
				Selected              string `xml:"selected,attr"`
				Stream                []struct {
					Text                 string `xml:",chardata"`
					BitDepth             string `xml:"bitDepth,attr"`
					Bitrate              string `xml:"bitrate,attr"`
					ChromaLocation       string `xml:"chromaLocation,attr"`
					ChromaSubsampling    string `xml:"chromaSubsampling,attr"`
					Codec                string `xml:"codec,attr"`
					CodedHeight          string `xml:"codedHeight,attr"`
					CodedWidth           string `xml:"codedWidth,attr"`
					ColorPrimaries       string `xml:"colorPrimaries,attr"`
					ColorRange           string `xml:"colorRange,attr"`
					ColorSpace           string `xml:"colorSpace,attr"`
					ColorTrc             string `xml:"colorTrc,attr"`
					Default              string `xml:"default,attr"`
					DisplayTitle         string `xml:"displayTitle,attr"`
					ExtendedDisplayTitle string `xml:"extendedDisplayTitle,attr"`
					FrameRate            string `xml:"frameRate,attr"`
					Height               string `xml:"height,attr"`
					ID                   string `xml:"id,attr"`
					Language             string `xml:"language,attr"`
					LanguageCode         string `xml:"languageCode,attr"`
					LanguageTag          string `xml:"languageTag,attr"`
					Level                string `xml:"level,attr"`
					Profile              string `xml:"profile,attr"`
					RefFrames            string `xml:"refFrames,attr"`
					StreamType           string `xml:"streamType,attr"`
					Width                string `xml:"width,attr"`
					Decision             string `xml:"decision,attr"`
					Location             string `xml:"location,attr"`
					AudioChannelLayout   string `xml:"audioChannelLayout,attr"`
					BitrateMode          string `xml:"bitrateMode,attr"`
					Channels             string `xml:"channels,attr"`
					SamplingRate         string `xml:"samplingRate,attr"`
					Selected             string `xml:"selected,attr"`
				} `xml:"Stream"`
			} `xml:"Part"`
		} `xml:"Media"`
		UltraBlurColors struct {
			Text        string `xml:",chardata"`
			BottomLeft  string `xml:"bottomLeft,attr"`
			BottomRight string `xml:"bottomRight,attr"`
			TopLeft     string `xml:"topLeft,attr"`
			TopRight    string `xml:"topRight,attr"`
		} `xml:"UltraBlurColors"`
		Rating struct {
			Text  string `xml:",chardata"`
			Image string `xml:"image,attr"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Rating"`
		Director struct {
			Text   string `xml:",chardata"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
			TagKey string `xml:"tagKey,attr"`
			Thumb  string `xml:"thumb,attr"`
		} `xml:"Director"`
		Writer struct {
			Text   string `xml:",chardata"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
			TagKey string `xml:"tagKey,attr"`
		} `xml:"Writer"`
		Role []struct {
			Text   string `xml:",chardata"`
			Filter string `xml:"filter,attr"`
			ID     string `xml:"id,attr"`
			Tag    string `xml:"tag,attr"`
			TagKey string `xml:"tagKey,attr"`
			Thumb  string `xml:"thumb,attr"`
			Role   string `xml:"role,attr"`
		} `xml:"Role"`
		User struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Thumb string `xml:"thumb,attr"`
			Title string `xml:"title,attr"`
		} `xml:"User"`
		Player struct {
			Text                string `xml:",chardata"`
			Address             string `xml:"address,attr"`
			Device              string `xml:"device,attr"`
			MachineIdentifier   string `xml:"machineIdentifier,attr"`
			Model               string `xml:"model,attr"`
			Platform            string `xml:"platform,attr"`
			PlatformVersion     string `xml:"platformVersion,attr"`
			Product             string `xml:"product,attr"`
			Profile             string `xml:"profile,attr"`
			RemotePublicAddress string `xml:"remotePublicAddress,attr"`
			State               string `xml:"state,attr"`
			Title               string `xml:"title,attr"`
			Version             string `xml:"version,attr"`
			Local               string `xml:"local,attr"`
			Relayed             string `xml:"relayed,attr"`
			Secure              string `xml:"secure,attr"`
			UserID              string `xml:"userID,attr"`
		} `xml:"Player"`
		Session struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			Bandwidth string `xml:"bandwidth,attr"`
			Location  string `xml:"location,attr"`
		} `xml:"Session"`
		TranscodeSession struct {
			Text                    string `xml:",chardata"`
			Key                     string `xml:"key,attr"`
			Throttled               string `xml:"throttled,attr"`
			Complete                string `xml:"complete,attr"`
			Progress                string `xml:"progress,attr"`
			Size                    string `xml:"size,attr"`
			Speed                   string `xml:"speed,attr"`
			Error                   string `xml:"error,attr"`
			Duration                string `xml:"duration,attr"`
			Remaining               string `xml:"remaining,attr"`
			Context                 string `xml:"context,attr"`
			SourceVideoCodec        string `xml:"sourceVideoCodec,attr"`
			SourceAudioCodec        string `xml:"sourceAudioCodec,attr"`
			VideoDecision           string `xml:"videoDecision,attr"`
			AudioDecision           string `xml:"audioDecision,attr"`
			Protocol                string `xml:"protocol,attr"`
			Container               string `xml:"container,attr"`
			VideoCodec              string `xml:"videoCodec,attr"`
			AudioCodec              string `xml:"audioCodec,attr"`
			AudioChannels           string `xml:"audioChannels,attr"`
			Width                   string `xml:"width,attr"`
			Height                  string `xml:"height,attr"`
			TranscodeHwRequested    string `xml:"transcodeHwRequested,attr"`
			TranscodeHwFullPipeline string `xml:"transcodeHwFullPipeline,attr"`
			TimeStamp               string `xml:"timeStamp,attr"`
			MaxOffsetAvailable      string `xml:"maxOffsetAvailable,attr"`
			MinOffsetAvailable      string `xml:"minOffsetAvailable,attr"`
		} `xml:"TranscodeSession"`
	} `xml:"Video"`
}

// HistorySessionResponse is the response on history stuff
type HistorySessionResponse struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Text    string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
	Video   []struct {
		Text                  string `xml:",chardata"`
		HistoryKey            string `xml:"historyKey,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		Title                 string `xml:"title,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		Type                  string `xml:"type,attr"`
		Index                 string `xml:"index,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		ViewedAt              string `xml:"viewedAt,attr"`
		AccountID             string `xml:"accountID,attr"`
		DeviceID              string `xml:"deviceID,attr"`
		Key                   string `xml:"key,attr"`
		RatingKey             string `xml:"ratingKey,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		Thumb                 string `xml:"thumb,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
	} `xml:"Video"`
	Track []struct {
		Text             string `xml:",chardata"`
		HistoryKey       string `xml:"historyKey,attr"`
		Key              string `xml:"key,attr"`
		RatingKey        string `xml:"ratingKey,attr"`
		LibrarySectionID string `xml:"librarySectionID,attr"`
		ParentKey        string `xml:"parentKey,attr"`
		GrandparentKey   string `xml:"grandparentKey,attr"`
		Title            string `xml:"title,attr"`
		ParentTitle      string `xml:"parentTitle,attr"`
		GrandparentTitle string `xml:"grandparentTitle,attr"`
		Type             string `xml:"type,attr"`
		Thumb            string `xml:"thumb,attr"`
		ParentThumb      string `xml:"parentThumb,attr"`
		GrandparentThumb string `xml:"grandparentThumb,attr"`
		GrandparentArt   string `xml:"grandparentArt,attr"`
		Index            string `xml:"index,attr"`
		ParentIndex      string `xml:"parentIndex,attr"`
		ViewedAt         string `xml:"viewedAt,attr"`
		AccountID        string `xml:"accountID,attr"`
		DeviceID         string `xml:"deviceID,attr"`
	} `xml:"Track"`
}
