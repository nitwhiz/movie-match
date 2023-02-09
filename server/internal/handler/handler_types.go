package handler

type MatchResult struct {
	OtherUserID string `gorm:"other_user_id" json:"otherUserId"`
}

type UserMatchParams struct {
	MediaType string `form:"type"`
	UserId    string `uri:"userId"`
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
	UserID  string `uri:"userId"`
	MediaID string `uri:"mediaId"`
}

type MediaParams struct {
	MediaID string `uri:"mediaId"`
}

type MediaRecommendationParams struct {
	UserId string `uri:"userId"`
	Page   int    `form:"page"`
}
