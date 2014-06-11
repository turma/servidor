package main

import (
	"github.com/huandu/facebook"
)

const FBTimeLayout = "2006-01-02T15:04:05+0000"

var FBUserParams = facebook.Params{
	"fields": "id,email,name,gender,link,locale,timezone,verified,picture,permissions",
}

var FBFeedParams = facebook.Params{
	"limit": 100,
	// Removed LIKES from the fields: likes.fields(id,name,picture)
	// For likes or comments count: likes.limit(1).summary(true),comments.limit(1).summary(true)
	"fields": "id,name,description,message,link,source,picture,caption,created_time,privacy,type,object_id",
}

var FBVideoParams = facebook.Params{
	// Field format removed, we will use just standard format
	// Not using embed_html, couse we will open a new iFram with source
	"fields": "id,from.fields(id,name,picture),name,description,source,picture,created_time",
}

var FBPhotoParams = facebook.Params{
	// Another image sizes not in use, removed from fields: images
	"fields": "id,from.fields(id,name,picture),name,height,width,link,source,picture,created_time",
}

var FBLinkParams = facebook.Params{
	"fields": "id,name,description,link,likes,picture",
}

var FBSummaryParams = facebook.Params{
	"limit":   0, // I don't need the likes information
	"summary": true,
}

type FBUser struct {
	Id          string
	Email       string
	Name        string
	Gender      string
	Link        string
	Locale      string
	Timezone    int
	Verified    bool
	Picture     FBPicture
	Permissions FBPermissions
}

// This struct gets just the picture of the object
type FBObjectPicture struct {
	Picture FBPicture
}

type FBPicture struct {
	Data FBPictureData
}

type FBPictureData struct {
	Url string
}

type FBPermissions struct {
	Data []FBPermissionsData
}

type FBPermissionsData struct {
	Permission string
	Status     string
}

type FBFeed struct {
	Data []FBFeedData
}

type FBFeedData struct {
	Id          string // user id _ post id
	Name        string
	Description string
	Message     string
	Link        string
	Source      string // For video player
	Picture     string
	Caption     string // www.youtube.com for videos
	CreatedTime string
	Privacy     FBPrivacy
	Type        string // photo, video, swf
	ObjectId    string // If exist's, just share if it's accessible

}

type FBPrivacy struct {
	Value string
}

type FBLinkInfo struct {
	Id     string
	Link   string
	Shares int
}

type FBLink struct {
	Id          string
	Name        string
	Description string
	Link        string
	Likes       int
	Picture     FBPicture
}

type FBImage struct {
	Id          string // user id _ post id
	Name        string
	Height      int
	Width       int
	Link        string
	Source      string // Photo URL
	Picture     string
	CreatedTime string
	From        FBFrom

	// We will use just the standard image size
	// Images []FBImages
}

// Not in use, we will use de default 720x720 max image
//
// type FBImages struct {
// 	Source string
// 	Height int
// 	Width  int
// }

type FBVideo struct {
	Id          string // user id _ post id
	Name        string
	Description string
	EmbedHtml   string
	Source      string // For video player
	Picture     string
	CreatedTime string
	From        FBFrom

	// We will use just the standard movie format
	//Format []FBFormat
}

// type FBFormat struct {
// 	EmbedHtml string
// 	Width     int
// 	Height    int
// 	Filter    string
// 	// 130x130 (real:130x98), 480x480 (real:480x380), native (real:640x480)
// }

type FBFrom struct {
	Name    string
	Id      string
	Picture FBPicture
}

// type FBLikes struct {
// 	Summary FBSummary
// }

type FBLikesSummary struct {
	Summary FBSummary
}

type FBSummary struct {
	TotalCount int
}

type FBShares struct {
	Shares int
}
