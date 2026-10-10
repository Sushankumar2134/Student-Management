package main

type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	City   string `json:"city"`
}

type User struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	PasswordHash string  `json:"-"`
	Role         string  `json:"role"`
	StudentID    *int    `json:"student_id,omitempty"`
	Email        string  `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message   string `json:"message"`
	Token     string `json:"token"`
	Role      string `json:"role"`
	StudentID *int   `json:"student_id"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	City     string `json:"city"`
}

type RegisterResponse struct {
	Message   string `json:"message"`
	StudentID int    `json:"student_id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ProfileResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	City   string `json:"city"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
	City  string `json:"city"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ForgotPasswordRequest struct {
	Username        string `json:"username"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type StatisticsResponse struct {
	TotalStudents  int `json:"total_students"`
	MaleStudents   int `json:"male_students"`
	FemaleStudents int `json:"female_students"`
}

type Claims struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
	StudentID *int   `json:"student_id"`
	Exp      int64   `json:"exp"`
}