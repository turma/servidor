package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const YouTubeTimeLayout = "2006-01-02T15:04:05.000Z"

//
// This function inserts and/or updates videos
//
func InsertYouTubeVideo(fb FB, user *User, feed FBFeedData) {

	chunks := strings.Split(feed.Source, "?")
	chunks = strings.Split(chunks[0], "/")

	youtubeId := chunks[len(chunks)-1]

	//log.Printf("id: %s type: %s", youtubeId, feed.Type)

	// Check if this youtube video is already saved
	videoobj, err := db.Get(YouTubeVideo{}, youtubeId)
	if err != nil {
		log.Printf("Error getting the youtube video '%s'. %s", youtubeId, err)
		return
	}

	if videoobj == nil {

		// Try to get likes from this Youtube video url
		res, err := fb.Get("https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3D"+youtubeId, nil)
		if err != nil { // Facebook object_id not exists
			log.Printf("Error getting the youtube video '%s' likes! %s.", youtubeId, err)
			return
		}

		var fbShares FBShares
		err = res.Decode(&fbShares)
		if err != nil {
			log.Printf("Error decoding the youtube video '%s' likes! %s.", youtubeId, err)
			return
		}

		response, err := http.Get("https://www.googleapis.com/youtube/v3/videos?id=" + youtubeId + "&key=" + EnvDev.Youtube + "&part=snippet,statistics&fields=items(id,snippet(title,description,publishedAt,thumbnails(default(url),medium(url),high(url))),statistics(likeCount,dislikeCount,favoriteCount))")
		if err != nil {
			log.Printf("Error getting info in youtube %s", err)
			return
		}
		defer response.Body.Close()

		var item YouTube
		json.NewDecoder(response.Body).Decode(&item)

		//log.Printf("item: %#v", item)

		videodata := item.Items[0]

		// Extract a time format from the youtube json time format
		createdTime, err := time.Parse(YouTubeTimeLayout, videodata.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error extracting CreatedTime from a youtube video. %s", err)
			return
		}

		youtubeLikes, err := GetYouTubeLikes(videodata)
		if err != nil {
			log.Printf("Error a youtube video likes. %s", err)
			return
		}

		video := &YouTubeVideo{
			Id:             videodata.Id,
			Link:           "https://www.youtube.com/watch?v=" + videodata.Id,
			Source:         "https://www.youtube.com/v/" + videodata.Id + "?version=3&autohide=1&autoplay=1",
			Likes:          fbShares.Shares,
			YouTubeLikes:   youtubeLikes,
			Title:          videodata.Snippet.Title,
			Description:    videodata.Snippet.Description,
			PictureDefault: videodata.Snippet.Thumbnails["default"].Url,
			PictureMedium:  videodata.Snippet.Thumbnails["medium"].Url,
			PictureHigh:    videodata.Snippet.Thumbnails["high"].Url,
			CreatedTime:    createdTime,
		}

		//log.Printf("video: %#v", video)

		err = db.Insert(video)
		if err != nil {
			log.Printf("Error insterting a new video '%s'! %s.", video.Id, err)
			return
		}

		log.Printf("Video inserted '%s'!.", video.Id)

	} else {
		log.Printf("Youtube video '%s' already exists.", youtubeId)
		return
	}

	//
	// I'm not saving that this user is sharing this video
	//
}

//
// It's my formula to get an YouTube video likes
// Just sum Likes - Dislikes + Favorite
//
func GetYouTubeLikes(videodata YouTubeItem) (int, error) {

	// Convert YouTube stupid strings to integer
	like, err := strconv.Atoi(videodata.Statistics.LikeCount)
	if err != nil {
		return 0, err
	}

	dislike, err := strconv.Atoi(videodata.Statistics.DislikeCount)
	if err != nil {
		return 0, err
	}

	favorite, err := strconv.Atoi(videodata.Statistics.FavoriteCount)
	if err != nil {
		return 0, err
	}

	return like - dislike + favorite, nil
}

type YouTube struct {
	Items []YouTubeItem `json:"items"`
}

type YouTubeItem struct {
	Id         string            `json:"id"`
	Snippet    YouTubeSnippet    `json:"snippet"`
	Statistics YouTubeStatistics `json:"statistics"`
}

type YouTubeSnippet struct {
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	PublishedAt string                      `json:"publishedAt"`
	Thumbnails  map[string]YouTubeThumbnail `json:"thumbnails"`
}

// Not using full Thumbnail
type YouTubeThumbnail struct {
	Url string `json:"url"`
}

// Not using full Statistics statistics(viewCount,likeCount,dislikeCount,commentCount)
type YouTubeStatistics struct {
	LikeCount     string `json:"likeCount"`
	DislikeCount  string `json:"dislikeCount"`
	FavoriteCount string `json:"favoriteCount"`
}
