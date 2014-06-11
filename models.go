package main

import (
	"time"
)

type User struct {
	Id       string `db:"id" json:"id"`             // *PK VARCHAR(255)
	Email    string `db:"email" json:"email"`       // VARCHAR(255)
	Name     string `db:"name" json:"name"`         // VARCHAR(255)
	Picture  string `db:"picture" json:"picture"`   // VARCHAR(255)
	Gender   string `db:"gender" json:"gender"`     // ENUM('male', 'female')
	Link     string `db:"link" json:"link"`         // VARCHAR(255)
	Locale   string `db:"locale" json:"locale"`     // VARCHAR(5)
	Timezone int    `db:"timezone" json:"timezone"` // TINYINT
	Verified bool   `db:"verified" json:"verified"` // BOOL

	AccessToken          string    `db:"access_token" json:"access_token"`                       // VARCHAR(255)
	LongLivedAccessToken string    `db:"long_lived_access_token" json:"long_lived_access_token"` // VARCHAR(255)
	TokenExpiresAt       time.Time `db:"token_expires_at" json:"token_expires_at"`               // DATETIME

	EmailPermission       bool `db:"email_permission" json:"email_permission"`               // BOOL
	ReadStreamPermission  bool `db:"read_stream_permission" json:"read_stream_permission"`   // BOOL
	UserFriendsPermission bool `db:"user_friends_permission" json:"user_friends_permission"` // BOOL

	LastRead    time.Time `db:"last_read" json:"-"`               // DATETIME
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
}

type Link struct {
	Link        string    `db:"link" json:"link"`                 // *PK VARCHAR(255)
	Id          string    `db:"id" json:"id"`                     // VARCHAR(255) Facebook Id
	Likes       int       `db:"likes" json:"likes"`               // INT Likes or Shares (Likes + Shares + Comments)
	Name        string    `db:"name" json:"name"`                 // VARCHAR(255)
	Description string    `db:"description" json:"description"`   // VARCHAR(255)
	Caption     string    `db:"caption" json:"caption"`           // VARCHAR(255)
	Picture     string    `db:"picture" json:"picture"`           // VARCHAR(255)
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
}

type Shared struct {
	Id          string    `db:"id" json:"id"`                     // *PK VARCHAR(255)
	Link        string    `db:"link" json:"link"`                 // *FK VARCHAR(255)
	From        string    `db:"from" json:"from"`                 // *FK VARCHAR(255)
	Name        string    `db:"name" json:"name"`                 // VARCHAR(255)
	Description string    `db:"description" json:"description"`   // VARCHAR(255)
	Message     string    `db:"message" json:"message"`           // VARCHAR(255)
	Caption     string    `db:"caption" json:"caption"`           // VARCHAR(255)
	Picture     string    `db:"picture" json:"picture"`           // VARCHAR(255)
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
}

type Photo struct {
	Id          string    `db:"id" json:"id"`                     // *PK VARCHAR(255)
	Likes       int       `db:"likes" json:"likes"`               // INT Likes
	Name        string    `db:"name" json:"name"`                 // VARCHAR(255)
	Height      int       `db:"height" json:"height"`             // INT
	Width       int       `db:"width" json:"width"`               // INT
	Link        string    `db:"link" json:"link"`                 // VARCHAR(255)
	Source      string    `db:"source" json:"source"`             // VARCHAR(255)
	Picture     string    `db:"picture" json:"picture"`           // VARCHAR(255)
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
	FromId      string    `db:"from_id" json:"from_id"`           // VARCHAR(255)
	FromName    string    `db:"from_name" json:"from_name"`       // VARCHAR(255)
	FromPicture string    `db:"from_picture" json:"from_picture"` // VARCHAR(255)
}

type Video struct {
	Id          string    `db:"id" json:"id"`                     // *PK VARCHAR(255)
	Likes       int       `db:"likes" json:"likes"`               // INT Likes
	Name        string    `db:"name" json:"name"`                 // VARCHAR(255)
	Description string    `db:"description" json:"description"`   // VARCHAR(255)
	Source      string    `db:"source" json:"source"`             // VARCHAR(255)
	Picture     string    `db:"picture" json:"picture"`           // VARCHAR(255)
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
	FromId      string    `db:"from_id" json:"from_id"`           // VARCHAR(255)
	FromName    string    `db:"from_name" json:"from_name"`       // VARCHAR(255)
	FromPicture string    `db:"from_picture" json:"from_picture"` // VARCHAR(255)
}

type YouTubeVideo struct {
	Id             string    `db:"id" json:"id"`                           // *PK VARCHAR(255)
	Link           string    `db:"link" json:"link"`                       // VARCHAR(255)
	Source         string    `db:"source" json:"source"`                   // VARCHAR(255)
	Likes          int       `db:"likes" json:"likes"`                     // INT Likes
	YouTubeLikes   int       `db:"youtubelikes" json:"youtubelikes"`       // INT Youtube Likes
	Title          string    `db:"title" json:"title"`                     // VARCHAR(255)
	Description    string    `db:"description" json:"description"`         // VARCHAR(255)
	PictureDefault string    `db:"picture_default" json:"picture_default"` // VARCHAR(255)
	PictureMedium  string    `db:"picture_medium" json:"picture_medium"`   // VARCHAR(255)
	PictureHigh    string    `db:"picture_high" json:"picture_high"`       // VARCHAR(255)
	CreatedTime    time.Time `db:"created_time" json:"created_time"`       // DATETIME
}

// type Likes struct {
// 	Id          string    `db:"id" json:"id"`                     // *PK VARCHAR(255)
// 	Link        string    `db:"link" json:"link"`                 // *FK VARCHAR(255)
// 	From        string    `db:"from" json:"from"`                 // *FK VARCHAR(255)
// 	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
// }
