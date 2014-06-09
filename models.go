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
	Shares      int       `db:"shares" json:"shares"`             // INT (Likes + Shares + Comments)
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

type Likes struct {
	Id          string    `db:"id" json:"id"`                     // *PK VARCHAR(255)
	Link        string    `db:"link" json:"link"`                 // *FK VARCHAR(255)
	From        string    `db:"from" json:"from"`                 // *FK VARCHAR(255)
	CreatedTime time.Time `db:"created_time" json:"created_time"` // DATETIME
}
