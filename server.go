package goflex

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

// ServerService describes the Server endpoints
type ServerService interface {
	Identity() (*IdentityResponse, error)
	MachineID() (string, error)
	Preferences() (*Preferences, error)
	Capabilities() (*Capabilities, error)
	Servers() (*Servers, error)
	Accounts() (*Accounts, error)
	Search(string) (*Search, error)
	// Notifications()
}

// ServerServiceOp is the operator for the ServerService
type ServerServiceOp struct {
	p *Plex
	// identityCache *IdentityResponse
}

// Search searches the plex libraries
func (svc *ServerServiceOp) Search(q string) (*Search, error) {
	var ret searchResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/search?query=%v", svc.p.baseURL, url.QueryEscape(q))), &ret, nil); err != nil {
		return nil, err
	}
	return &ret.Search, nil
}

// Accounts returns accounts
func (svc *ServerServiceOp) Accounts() (*Accounts, error) {
	var ret accountsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/accounts", svc.p.baseURL)), &ret, &cacheConfig{prefix: "accounts", ttl: time.Hour * 6}); err != nil {
		return nil, err
	}
	return &ret.Accounts, nil
}

// Servers returns a list of plex servers
func (svc *ServerServiceOp) Servers() (*Servers, error) {
	var ret serversResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/servers", svc.p.baseURL)), &ret, &cacheConfig{prefix: "servers", ttl: time.Hour * 6}); err != nil {
		return nil, err
	}
	return &ret.Servers, nil
}

// Capabilities returns the capabilities of a host
func (svc *ServerServiceOp) Capabilities() (*Capabilities, error) {
	var ret capabilitiesResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/", svc.p.baseURL)), &ret, &cacheConfig{prefix: "capabilities", ttl: time.Hour * 6}); err != nil {
		return nil, err
	}
	return &ret.Capabilities, nil
}

// Preferences returns server preferences
func (svc *ServerServiceOp) Preferences() (*Preferences, error) {
	var ret prefsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/:/prefs", svc.p.baseURL)), &ret, &cacheConfig{prefix: "prefs", ttl: time.Hour * 1}); err != nil {
		return nil, err
	}
	return &ret.Preferences, nil
}

// MachineID returns the ServerID
func (svc *ServerServiceOp) MachineID() (string, error) {
	got, err := svc.Identity()
	if err != nil {
		return "", err
	}
	return got.MachineIdentifier, nil
}

func (svc *ServerServiceOp) Identity() (*IdentityResponse, error) {
	req := mustNewRequest("GET", fmt.Sprintf("%v/identity", svc.p.baseURL))
	var ret IdentityResponse
	if err := svc.p.sendRequestXML(req, &ret, nil); err != nil {
		return nil, err
	}
	return &ret, nil
}

// IdentityResponse is the response back from the identity endpoint
type IdentityResponse struct {
	XMLName           xml.Name `xml:"MediaContainer"`
	Text              string   `xml:",chardata"`
	Size              string   `xml:"size,attr"`
	APIVersion        string   `xml:"apiVersion,attr"`
	Claimed           string   `xml:"claimed,attr"`
	MachineIdentifier string   `xml:"machineIdentifier,attr"`
	Version           string   `xml:"version,attr"`
}

type capabilitiesResponse struct {
	Capabilities Capabilities `json:"MediaContainer"`
}

// Capabilities describes the capabilities of the connected server
type Capabilities struct {
	Size                          int    `json:"size"`
	Allowcameraupload             bool   `json:"allowCameraUpload"`
	Allowchannelaccess            bool   `json:"allowChannelAccess"`
	Allowmediadeletion            bool   `json:"allowMediaDeletion"`
	Allowsharing                  bool   `json:"allowSharing"`
	Allowsync                     bool   `json:"allowSync"`
	Allowtuners                   bool   `json:"allowTuners"`
	Apiversion                    string `json:"apiVersion"`
	Backgroundprocessing          bool   `json:"backgroundProcessing"`
	Certificate                   bool   `json:"certificate"`
	Companionproxy                bool   `json:"companionProxy"`
	Countrycode                   string `json:"countryCode"`
	Diagnostics                   string `json:"diagnostics"`
	Eventstream                   bool   `json:"eventStream"`
	Friendlyname                  string `json:"friendlyName"`
	Hubsearch                     bool   `json:"hubSearch"`
	Itemclusters                  bool   `json:"itemClusters"`
	Livetv                        int    `json:"livetv"`
	Machineidentifier             string `json:"machineIdentifier"`
	Mediaproviders                bool   `json:"mediaProviders"`
	Multiuser                     bool   `json:"multiuser"`
	Musicanalysis                 int    `json:"musicAnalysis"`
	Myplex                        bool   `json:"myPlex"`
	Myplexmappingstate            string `json:"myPlexMappingState"`
	Myplexsigninstate             string `json:"myPlexSigninState"`
	Myplexsubscription            bool   `json:"myPlexSubscription"`
	Myplexusername                string `json:"myPlexUsername"`
	Offlinetranscode              int    `json:"offlineTranscode"`
	Ownerfeatures                 string `json:"ownerFeatures"`
	Platform                      string `json:"platform"`
	Platformversion               string `json:"platformVersion"`
	Pluginhost                    bool   `json:"pluginHost"`
	Pushnotifications             bool   `json:"pushNotifications"`
	Readonlylibraries             bool   `json:"readOnlyLibraries"`
	Streamingbrainabrversion      int    `json:"streamingBrainABRVersion"`
	Streamingbrainversion         int    `json:"streamingBrainVersion"`
	Sync                          bool   `json:"sync"`
	Transcoderactivevideosessions int    `json:"transcoderActiveVideoSessions"`
	Transcoderaudio               bool   `json:"transcoderAudio"`
	Transcoderlyrics              bool   `json:"transcoderLyrics"`
	Transcoderphoto               bool   `json:"transcoderPhoto"`
	Transcodersubtitles           bool   `json:"transcoderSubtitles"`
	Transcodervideo               bool   `json:"transcoderVideo"`
	Transcodervideobitrates       string `json:"transcoderVideoBitrates"`
	Transcodervideoqualities      string `json:"transcoderVideoQualities"`
	Transcodervideoresolutions    string `json:"transcoderVideoResolutions"`
	Updatedat                     int    `json:"updatedAt"`
	Updater                       bool   `json:"updater"`
	Version                       string `json:"version"`
	Voicesearch                   bool   `json:"voiceSearch"`
	Directory                     []struct {
		Count int    `json:"count"`
		Key   string `json:"key"`
		Title string `json:"title"`
	} `json:"Directory"`
}

type prefsResponse struct {
	Preferences Preferences `json:"MediaContainer"`
}

// Preferences are the preferences set on the server
type Preferences struct {
	Size    int `json:"size"`
	Setting []struct {
		ID         string `json:"id"`
		Label      string `json:"label"`
		Summary    string `json:"summary"`
		Type       string `json:"type"`
		Default    any    `json:"default"`
		Value      any    `json:"value"`
		Hidden     bool   `json:"hidden"`
		Advanced   bool   `json:"advanced"`
		Group      string `json:"group"`
		EnumValues string `json:"enumValues,omitempty"`
	} `json:"Setting"`
}

type serversResponse struct {
	Servers Servers `json:"MediaContainer"`
}

// Servers is multiple plex servers
type Servers struct {
	Size   int      `json:"size"`
	Server []Server `json:"Server"`
}

// Server represents a plex server
type Server struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	Address           string `json:"address"`
	Port              int    `json:"port"`
	MachineIdentifier string `json:"machineIdentifier"`
	Version           string `json:"version"`
}

type accountsResponse struct {
	Accounts Accounts `json:"MediaContainer"`
}

// Accounts are accounts tied to a given server
type Accounts struct {
	Size       int    `json:"size"`
	Identifier string `json:"identifier"`
	Account    []struct {
		ID                      int    `json:"id"`
		Key                     string `json:"key"`
		Name                    string `json:"name"`
		DefaultAudioLanguage    string `json:"defaultAudioLanguage"`
		AutoSelectAudio         bool   `json:"autoSelectAudio"`
		DefaultSubtitleLanguage string `json:"defaultSubtitleLanguage"`
		SubtitleMode            int    `json:"subtitleMode"`
		Thumb                   string `json:"thumb"`
	} `json:"Account"`
}

type searchResponse struct {
	Search Search `json:"MediaContainer"`
}

// Metadata is metadata for a search item
type Metadata struct {
	AllowSync             bool    `json:"allowSync"`
	LibrarySectionID      int     `json:"librarySectionID"`
	LibrarySectionTitle   string  `json:"librarySectionTitle"`
	LibrarySectionUUID    string  `json:"librarySectionUUID"`
	Personal              bool    `json:"personal"`
	SourceTitle           string  `json:"sourceTitle"`
	RatingKey             string  `json:"ratingKey"`
	Key                   string  `json:"key"`
	GUID                  string  `json:"guid"`
	Studio                string  `json:"studio,omitempty"`
	Type                  string  `json:"type"`
	Title                 string  `json:"title"`
	Summary               string  `json:"summary"`
	Rating                float64 `json:"rating,omitempty"`
	Year                  int     `json:"year,omitempty"`
	Tagline               string  `json:"tagline,omitempty"`
	Thumb                 string  `json:"thumb"`
	Art                   string  `json:"art"`
	Duration              int     `json:"duration"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	AddedAt               int     `json:"addedAt"`
	UpdatedAt             int     `json:"updatedAt,omitempty"`
	Media                 []struct {
		ID              int     `json:"id"`
		Duration        int     `json:"duration"`
		Bitrate         int     `json:"bitrate"`
		Width           int     `json:"width"`
		Height          int     `json:"height"`
		AspectRatio     float64 `json:"aspectRatio"`
		AudioChannels   int     `json:"audioChannels"`
		AudioCodec      string  `json:"audioCodec"`
		VideoCodec      string  `json:"videoCodec"`
		VideoResolution string  `json:"videoResolution"`
		Container       string  `json:"container"`
		VideoFrameRate  string  `json:"videoFrameRate"`
		VideoProfile    string  `json:"videoProfile"`
		Part            []struct {
			ID           int    `json:"id"`
			Key          string `json:"key"`
			Duration     int    `json:"duration"`
			File         string `json:"file"`
			Size         int    `json:"size"`
			Container    string `json:"container"`
			VideoProfile string `json:"videoProfile"`
		} `json:"Part"`
	} `json:"Media,omitempty"`
	Image []struct {
		Alt  string `json:"alt"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"Image"`
	UltraBlurColors struct {
		TopLeft     string `json:"topLeft"`
		TopRight    string `json:"topRight"`
		BottomRight string `json:"bottomRight"`
		BottomLeft  string `json:"bottomLeft"`
	} `json:"UltraBlurColors"`
	Genre []struct {
		Tag string `json:"tag"`
	} `json:"Genre,omitempty"`
	Country []struct {
		Tag string `json:"tag"`
	} `json:"Country,omitempty"`
	Director []struct {
		Tag string `json:"tag"`
	} `json:"Director,omitempty"`
	Writer []struct {
		Tag string `json:"tag"`
	} `json:"Writer,omitempty"`
	Role []struct {
		Tag string `json:"tag"`
	} `json:"Role,omitempty"`
	Slug                   string  `json:"slug,omitempty"`
	ContentRating          string  `json:"contentRating,omitempty"`
	Index                  int     `json:"index,omitempty"`
	AudienceRating         float64 `json:"audienceRating,omitempty"`
	ViewCount              int     `json:"viewCount,omitempty"`
	SkipCount              int     `json:"skipCount,omitempty"`
	LastViewedAt           int     `json:"lastViewedAt,omitempty"`
	Theme                  string  `json:"theme,omitempty"`
	LeafCount              int     `json:"leafCount,omitempty"`
	ViewedLeafCount        int     `json:"viewedLeafCount,omitempty"`
	ChildCount             int     `json:"childCount,omitempty"`
	AudienceRatingImage    string  `json:"audienceRatingImage,omitempty"`
	HasPremiumExtras       string  `json:"hasPremiumExtras,omitempty"`
	HasPremiumPrimaryExtra string  `json:"hasPremiumPrimaryExtra,omitempty"`
	PrimaryExtraKey        string  `json:"primaryExtraKey,omitempty"`
	ParentRatingKey        string  `json:"parentRatingKey,omitempty"`
	GrandparentRatingKey   string  `json:"grandparentRatingKey,omitempty"`
	ParentGUID             string  `json:"parentGuid,omitempty"`
	GrandparentGUID        string  `json:"grandparentGuid,omitempty"`
	GrandparentSlug        string  `json:"grandparentSlug,omitempty"`
	TitleSort              string  `json:"titleSort,omitempty"`
	GrandparentKey         string  `json:"grandparentKey,omitempty"`
	ParentKey              string  `json:"parentKey,omitempty"`
	GrandparentTitle       string  `json:"grandparentTitle,omitempty"`
	ParentTitle            string  `json:"parentTitle,omitempty"`
	ParentIndex            int     `json:"parentIndex,omitempty"`
	ParentThumb            string  `json:"parentThumb,omitempty"`
	GrandparentThumb       string  `json:"grandparentThumb,omitempty"`
	GrandparentArt         string  `json:"grandparentArt,omitempty"`
	GrandparentTheme       string  `json:"grandparentTheme,omitempty"`
	ChapterSource          string  `json:"chapterSource,omitempty"`
	ViewOffset             int     `json:"viewOffset,omitempty"`
}

// Search represents search results
type Search struct {
	Size            int        `json:"size"`
	Identifier      string     `json:"identifier"`
	MediaTagPrefix  string     `json:"mediaTagPrefix"`
	MediaTagVersion int        `json:"mediaTagVersion"`
	Metadata        []Metadata `json:"Metadata"`
	Provider        []struct {
		Key    string `json:"key"`
		MTitle string `json:"m_title"`
		MType  string `json:"m_type"`
	} `json:"Provider"`
}

// Episodes returns episodes from a search result
func (s Search) Episodes() (EpisodeList, error) {
	ret := EpisodeList{}
	for _, item := range s.Metadata {
		if item.Type == MediaTypeEpisode {
			e, err := episodeWith(item)
			if err != nil {
				return nil, err
			}
			ret = append(ret, *e)
		}
	}
	return ret, nil
}

// Shows returns a list of shows from the search
func (s Search) Shows() (ShowList, error) {
	ret := ShowList{}
	for _, item := range s.Metadata {
		if item.Type == MediaTypeShow {
			e, err := showWith(item)
			if err != nil {
				return nil, err
			}
			ret = append(ret, e)
		}
	}
	return ret, nil
}
