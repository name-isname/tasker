package core

import (
	"time"
)

// Task represents a task in the system
type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddTask adds a new task
func AddTask(title string) (*Task, error) {
	task := &Task{
		Title:     title,
		Completed: false,
	}
	if err := DB.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// ListTasks returns all tasks
func ListTasks() ([]Task, error) {
	var tasks []Task
	if err := DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// CompleteTask marks a task as completed
func CompleteTask(id uint) error {
	return DB.Model(&Task{}).Where("id = ?", id).Update("completed", true).Error
}

// DeleteTask removes a task
func DeleteTask(id uint) error {
	return DB.Delete(&Task{}, id).Error
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return DB.AutoMigrate(&Task{})
}
