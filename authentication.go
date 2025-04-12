package goflex

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"

	"github.com/drewstinnett/inspectareq"
	"github.com/google/uuid"
)

// AuthenticationService is the description of the Authentication endpoints.
type AuthenticationService interface {
	Token(string, string) (string, error)
}

// AuthenticationServiceOp is the operator for the AuthenticationService.
type AuthenticationServiceOp struct {
	p *Flex
}

// Token returns a new token using username and password authentication.
func (svc *AuthenticationServiceOp) Token(username, password string) (string, error) {
	body := url.Values{}
	body.Add("login", username)
	body.Add("password", password)
	body.Add("noGuest", "true")
	body.Add("skipAuthentication", "true")

	req, err := http.NewRequest(
		http.MethodPost,
		"https://plex.tv/api/v2/users/signin",
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return "", err
	}
	cid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Plex-Client-Identifier", cid.String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Accept", "application/xml")

	if err := inspectareq.Print(req); err != nil {
		svc.p.logger.Warn("error printing request", "error", err)
	}

	got, err := svc.p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer dclose(got.Body)

	text, err := io.ReadAll(got.Body)
	if err != nil {
		return "", err
	}
	var ret TokenResponse
	if err := xml.Unmarshal(text, &ret); err != nil {
		return "", err
	}
	return ret.AuthToken, nil
}

// TokenResponse is what we get back from the new token request.
type TokenResponse struct {
	XMLName                 xml.Name `xml:"user"`
	Text                    string   `xml:",chardata"`
	ID                      string   `xml:"id,attr"`
	UUID                    string   `xml:"uuid,attr"`
	Username                string   `xml:"username,attr"`
	Title                   string   `xml:"title,attr"`
	Email                   string   `xml:"email,attr"`
	FriendlyName            string   `xml:"friendlyName,attr"`
	Locale                  string   `xml:"locale,attr"`
	Confirmed               string   `xml:"confirmed,attr"`
	JoinedAt                string   `xml:"joinedAt,attr"`
	EmailOnlyAuth           string   `xml:"emailOnlyAuth,attr"`
	HasPassword             string   `xml:"hasPassword,attr"`
	Protected               string   `xml:"protected,attr"`
	Thumb                   string   `xml:"thumb,attr"`
	AuthToken               string   `xml:"authToken,attr"`
	MailingListStatus       string   `xml:"mailingListStatus,attr"`
	MailingListActive       string   `xml:"mailingListActive,attr"`
	ScrobbleTypes           string   `xml:"scrobbleTypes,attr"`
	Country                 string   `xml:"country,attr"`
	SubscriptionDescription string   `xml:"subscriptionDescription,attr"`
	Restricted              string   `xml:"restricted,attr"`
	Anonymous               string   `xml:"anonymous,attr"`
	Home                    string   `xml:"home,attr"`
	Guest                   string   `xml:"guest,attr"`
	HomeSize                string   `xml:"homeSize,attr"`
	HomeAdmin               string   `xml:"homeAdmin,attr"`
	MaxHomeSize             string   `xml:"maxHomeSize,attr"`
	RememberExpiresAt       string   `xml:"rememberExpiresAt,attr"`
	AdsConsent              string   `xml:"adsConsent,attr"`
	AdsConsentSetAt         string   `xml:"adsConsentSetAt,attr"`
	AdsConsentReminderAt    string   `xml:"adsConsentReminderAt,attr"`
	ExperimentalFeatures    string   `xml:"experimentalFeatures,attr"`
	TwoFactorEnabled        string   `xml:"twoFactorEnabled,attr"`
	BackupCodesCreated      string   `xml:"backupCodesCreated,attr"`
	AttributionPartner      string   `xml:"attributionPartner,attr"`
	Subscription            struct {
		Text           string `xml:",chardata"`
		Active         string `xml:"active,attr"`
		SubscribedAt   string `xml:"subscribedAt,attr"`
		Status         string `xml:"status,attr"`
		PaymentService string `xml:"paymentService,attr"`
		Plan           string `xml:"plan,attr"`
		Features       struct {
			Text    string `xml:",chardata"`
			Feature []struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"feature"`
		} `xml:"features"`
	} `xml:"subscription"`
	Profile struct {
		Text                         string `xml:",chardata"`
		AutoSelectAudio              string `xml:"autoSelectAudio,attr"`
		DefaultAudioAccessibility    string `xml:"defaultAudioAccessibility,attr"`
		DefaultAudioLanguage         string `xml:"defaultAudioLanguage,attr"`
		DefaultAudioLanguages        string `xml:"defaultAudioLanguages,attr"`
		DefaultSubtitleLanguage      string `xml:"defaultSubtitleLanguage,attr"`
		DefaultSubtitleLanguages     string `xml:"defaultSubtitleLanguages,attr"`
		AutoSelectSubtitle           string `xml:"autoSelectSubtitle,attr"`
		DefaultSubtitleAccessibility string `xml:"defaultSubtitleAccessibility,attr"`
		DefaultSubtitleForced        string `xml:"defaultSubtitleForced,attr"`
		WatchedIndicator             string `xml:"watchedIndicator,attr"`
		MediaReviewsVisibility       string `xml:"mediaReviewsVisibility,attr"`
		MediaReviewsLanguages        string `xml:"mediaReviewsLanguages,attr"`
	} `xml:"profile"`
	Entitlements      string `xml:"entitlements"`
	Subscriptions     string `xml:"subscriptions"`
	PastSubscriptions string `xml:"pastSubscriptions"`
	Trials            string `xml:"trials"`
	Services          struct {
		Text    string `xml:",chardata"`
		Service []struct {
			Text       string `xml:",chardata"`
			Identifier string `xml:"identifier,attr"`
			Endpoint   string `xml:"endpoint,attr"`
			Token      string `xml:"token,attr"`
			Secret     string `xml:"secret,attr"`
			Status     string `xml:"status,attr"`
		} `xml:"service"`
	} `xml:"services"`
}
