package models

import (
	"time"
)

type Tweet struct {
	ID int64 `json:"id"`
	Message string `json:"message"`
	Likes int32 `json:"likes"`
	Retweet int32 `json:"retweets"`
	Comments int32 `json:"comments"`
	Type string `json:"type"`
	CreatedBy int64 `json:"createdby"`
	CreatedOn *time.Time `json:"createdon"`
}

type Tweets struct{
	Original Tweet `json:"original,omitempty"`
	Parent Tweet `json:"parent,omitempty"`
}