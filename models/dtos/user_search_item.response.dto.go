package dtos

type UserSearchItemResponseDto struct {
	Username  *string `json:"username"`
	UserPubId *string `json:"userPubId" gorm:"user_pub_id"`
}
