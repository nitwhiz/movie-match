package handler

type MatchResult struct {
	MediaID     string `gorm:"media_id" json:"mediaId"`
	OtherUserID string `gorm:"other_user_id" json:"otherUserId"`
}

type UserMatchParams struct {
	UserId string `uri:"userId"`
}

type UserVoteParams struct {
	UserId   string `uri:"userId"`
	MediaId  string `uri:"mediaId"`
	VoteType string `json:"voteType"`
}

type MediaPosterParams struct {
	MediaID string `uri:"mediaId"`
}

type UserSeenParams struct {
	UserId  string `uri:"userId"`
	MediaId string `uri:"mediaId"`
}

type MediaParams struct {
	MediaId string `uri:"mediaId"`
}

type MediaRecommendationParams struct {
	UserId string `uri:"userId"`
	Page   int    `form:"page"`
}
