package models

type User struct{
	ID int64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email string `json:"email,omitempty"`
	Name string `json:"name,omitempty"`
	Photo string `json:"photo,omitempty"`
	Bio string `json:"bio,omitempty"`
}

type FollowStats struct{
	Followers int32 `json:"followers"`
	Following int32 `json:"following"`
}