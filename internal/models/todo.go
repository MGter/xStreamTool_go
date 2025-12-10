package models

import (
	"time"
)

// Todo 待办事项模型
type Todo struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	Priority    int       `json:"priority" db:"priority"`
	Category    string    `json:"category,omitempty" db:"category"`
	DueDate     time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TodoRequest 创建/更新待办事项请求
type TodoRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description" binding:"max=1000"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority" binding:"min=1,max=5"`
	Category    string    `json:"category" binding:"max=50"`
	DueDate     time.Time `json:"due_date"`
}

// TodoResponse 待办事项响应
type TodoResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority"`
	Category    string    `json:"category,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	IsOverdue   bool      `json:"is_overdue"`
}

// ToResponse 转换为响应格式
func (t *Todo) ToResponse() TodoResponse {
	now := time.Now()
	isOverdue := !t.Completed && !t.DueDate.IsZero() && t.DueDate.Before(now)

	status := "进行中"
	if t.Completed {
		status = "已完成"
	} else if isOverdue {
		status = "已过期"
	}

	return TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    t.Priority,
		Category:    t.Category,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Status:      status,
		IsOverdue:   isOverdue,
	}
}

// FromRequest 从请求创建模型
func (t *Todo) FromRequest(req *TodoRequest) {
	t.Title = req.Title
	t.Description = req.Description
	t.Completed = req.Completed
	t.Priority = req.Priority
	t.Category = req.Category
	t.DueDate = req.DueDate
	t.UpdatedAt = time.Now()
}

// User 用户模型
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
