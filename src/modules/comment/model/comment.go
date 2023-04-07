package model

import "time"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Photo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserID   string `json:"user_id"`
}

type FetchedComment struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PhotoID   string     `json:"photo_id"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	User      *User      `json:"user"`
	Photo     *Photo     `json:"photo"`
}

type ResponseDataFetchedComment struct {
	Status string           `json:"status" example:"success"`
	Data   []FetchedComment `json:"data"`
}

type AddComment struct {
	Message string `json:"message" example:"A comment"`
	PhotoID string `json:"photo_id" example:"photo-123"`
}

type AddedComment struct {
	ID        string     `json:"id" example:"here is the generated comment id"`
	UserID    string     `json:"user_id" example:"here is the generated user id"`
	PhotoID   string     `json:"photo_id" example:"here is the generated photo id"`
	Message   string     `json:"message" example:"A comment"`
	CreatedAt *time.Time `json:"created_at" example:"the created at generated here"`
}

type ResponseDataAddedComment struct {
	Status string       `json:"status" example:"success"`
	Data   AddedComment `json:"data"`
}

type UpdateComment struct {
	Message string `json:"message" example:"A new comment"`
}

type UpdatedComment struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	UserID    string     `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ResponseDataUpdatedComment struct {
	Status string         `json:"status" example:"success"`
	Data   UpdatedComment `json:"data"`
}

type ResponseMessageDeletedComment struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your comment has been successfully deleted"`
}

type ResponseMessage struct {
	Status string `json:"status" example:"fail"`
	Data   string `json:"data" example:"the error explained here"`
}
