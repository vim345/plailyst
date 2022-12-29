package crawler

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"google.golang.org/api/youtube/v3"
)

const maxResults = 50

// Video represents a simple video object.
type Video struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ChannelTitle string `json:"channel_title"`
}

// Channel represents youtube channel schema.
type Channel struct {
	ID string `yaml:"id"`
}

// Configs is the config file for the crawler.
type Configs struct {
	Channels []Channel `yaml:"channels"`
	Teams    []string  `yaml:"teams"`
	Terms    []string  `yaml:"terms"`
	Playlist string    `yaml:"playlist"`
}

// Crawler crawls youtube channels and filters uploads that match the requirements.
type Crawler struct {
	service *youtube.Service
	configs *Configs
}

// NewCrawler Creates a new crawler object.
func NewCrawler(configs *Configs, service *youtube.Service) *Crawler {
	return &Crawler{service, configs}
}

func matcher(title string, items []string) bool {
	found := false
	for _, team := range items {
		if strings.Contains(strings.ToLower(title), strings.ToLower(team)) {
			found = true
			continue
		}
	}
	return found
}

func matches(title string, teams []string, keywords []string) bool {
	return matcher(title, teams) && matcher(title, keywords)
}

func (c *Crawler) exists(video *Video, videos []*youtube.PlaylistItem) bool {
	for _, item := range videos {
		if video.ID == item.Snippet.ResourceId.VideoId {
			return true
		}
	}
	return false
}

func (c *Crawler) getExistingVideos(playlistID string) ([]*youtube.PlaylistItem, error) {
	part := []string{"snippet"}
	call := c.service.PlaylistItems.List(part)
	call = call.PlaylistId(playlistID)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

// Run runs the crawler.
func (c *Crawler) Run() {
	existingVideos, err := c.getExistingVideos(c.configs.Playlist)
	if err != nil {
		log.Errorf("Cannot find existing videos %+v", err)
	}
	var wg sync.WaitGroup

	for _, ch := range c.configs.Channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			playlistID, err := c.GetPlaylistID(ch.ID)
			if err != nil {
				log.Errorf("Cannot find uploads for channel %+v, %+v\n", ch.ID, err)
				return
			}
			uploads, err := c.GetUploads(playlistID)
			if err != nil {
				log.Errorf("Cannot find uploads for channel %+v, %+v\n", ch.ID, err)
			}
			for _, upload := range uploads {
				if matches(upload.Title, c.configs.Teams, c.configs.Terms) == true {
					if !c.exists(upload, existingVideos) {
						log.Infof("Adding video to playlist = %+v\n", upload.Title)
						c.addToPlayList(upload)
					}
				}
			}
		}()

		wg.Wait()
	}

}

func (c *Crawler) addToPlayList(video *Video) {
	part := []string{"snippet"}
	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: c.configs.Playlist,
			Position:   0,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: video.ID,
			},
		},
	}
	call := c.service.PlaylistItems.Insert(part, playlistItem)
	_, err := call.Do()
	if err != nil {
		log.Warnf("Couldn't add %+v to playlist. Got %+v", video, err)
	} else {
		log.Infof("Added %+v to playlist", video)
	}
}

// GetPlaylistID gets uploads ID.
func (c *Crawler) GetPlaylistID(channelID string) (string, error) {
	channelParts := []string{"contentDetails"}
	call := c.service.Channels.List(channelParts)
	call = call.Id(channelID)
	resp, err := call.Do()
	if err != nil {
		return "", err
	}
	return resp.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

// GetUploads gets the uploaded videos for the given playlist.
func (c *Crawler) GetUploads(playlistID string) ([]*Video, error) {
	parts := []string{"snippet"}
	call := c.service.PlaylistItems.List(parts)
	call = call.PlaylistId(playlistID).MaxResults(maxResults)
	resp, err := call.Do()
	var videos []*Video
	if err != nil {
		return videos, err
	}
	for _, item := range resp.Items {
		video := &Video{
			item.Snippet.ResourceId.VideoId,
			item.Snippet.Title,
			item.Snippet.VideoOwnerChannelTitle,
		}
		videos = append(videos, video)
	}
	return videos, nil
}
