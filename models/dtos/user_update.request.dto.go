package dtos

type UserUpdateRequestDto struct {
	Username *string `json:"username"`
	IsPublic *bool   `json:"isPublic"`
}
