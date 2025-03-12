package goflex

import "fmt"

// MediaService describes the media endpoints
type MediaService interface {
	MarkWatched(int) error
	MarkUnWatched(int) error
	MarkEpisodeWatched(ShowTitle, SeasonNumber, EpisodeNumber) error
	MarkEpisodeUnWatched(ShowTitle, SeasonNumber, EpisodeNumber) error
}

// MediaServiceOp is the operator for the MediaService
type MediaServiceOp struct {
	p *Plex
}

// MarkWatched marks a piece of media as watched
func (svc *MediaServiceOp) MarkWatched(key int) error {
	var ret struct{}
	return svc.p.sendRequestXML(mustNewRequest("GET", fmt.Sprintf("%v/:/scrobble?identifier=com.plexapp.plugins.library&key=%v", svc.p.baseURL, key)), &ret, nil)
}

// MarkUnWatched marks a piece of media as watched
func (svc *MediaServiceOp) MarkUnWatched(key int) error {
	var ret struct{}
	return svc.p.sendRequestXML(mustNewRequest("GET", fmt.Sprintf("%v/:/unscrobble?identifier=com.plexapp.plugins.library&key=%v", svc.p.baseURL, key)), &ret, nil)
}

// MarkEpisodeWatched marks an episode as watched
func (svc *MediaServiceOp) MarkEpisodeWatched(show ShowTitle, season SeasonNumber, episode EpisodeNumber) error {
	key, err := svc.p.episodeID(show, season, episode)
	if err != nil {
		return err
	}
	return svc.MarkWatched(key)
}

// MarkEpisodeUnWatched marks an episode as watched
func (svc *MediaServiceOp) MarkEpisodeUnWatched(show ShowTitle, season SeasonNumber, episode EpisodeNumber) error {
	key, err := svc.p.episodeID(show, season, episode)
	if err != nil {
		return err
	}
	return svc.MarkUnWatched(key)
}
