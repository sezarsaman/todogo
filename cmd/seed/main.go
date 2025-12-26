package main

import (
	"context"

	"task-manager/internal/config"
	"task-manager/internal/platform/database"
	"task-manager/internal/task/model"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := database.New(ctx, postgresDSN(cfg))
	if err != nil {
		panic(err)
	}
	defer db.Conn.Close(ctx)

	tasks := []model.Task{
		{Title: "Buy groceries", Description: "Milk, Bread, Eggs", Status: "pending", Assignee: "Jack"},
		{Title: "Finish task project", Description: "Complete Golang task", Status: "in_progress", Assignee: "Jill"},
		{Title: "Workout", Description: "30 minutes running", Status: "done", Assignee: "John"},
		{Title: "Read book", Description: "Read 'The Go Programming Language'", Status: "pending", Assignee: "Jane"},
		{Title: "Clean house", Description: "Vacuum and dust", Status: "in_progress", Assignee: "Doe"},
		{Title: "Pay bills", Description: "Electricity and Internet", Status: "done", Assignee: "Smith"},
		{Title: "Call mom", Description: "Weekly check-in", Status: "pending", Assignee: "Emily"},
		{Title: "Plan vacation", Description: "Research destinations", Status: "in_progress", Assignee: "Michael"},
		{Title: "Organize files", Description: "Sort documents on PC", Status: "done", Assignee: "Sarah"},
		{Title: "Meditate", Description: "15 minutes mindfulness", Status: "pending", Assignee: "David"},
		{Title: "Cook dinner", Description: "Try new recipe", Status: "in_progress", Assignee: "Anna"},
		{Title: "Write blog post", Description: "Topic: Go concurrency", Status: "done", Assignee: "Chris"},
		{Title: "Attend meeting", Description: "Project kickoff", Status: "pending", Assignee: "Olivia"},
		{Title: "Fix bike", Description: "Change flat tire", Status: "in_progress", Assignee: "Ethan"},
		{Title: "Study Go", Description: "Complete online course", Status: "done", Assignee: "Sophia"},
		{Title: "Gardening", Description: "Plant new flowers", Status: "pending", Assignee: "Liam"},
		{Title: "Update resume", Description: "Add recent experience", Status: "in_progress", Assignee: "Mia"},
		{Title: "Volunteer", Description: "Local community center", Status: "done", Assignee: "Noah"},
		{Title: "Learn guitar", Description: "Practice chords", Status: "pending", Assignee: "Isabella"},
		{Title: "Watch documentary", Description: "Topic: Nature", Status: "in_progress", Assignee: "James"},
		{Title: "Bake cake", Description: "Chocolate flavor", Status: "done", Assignee: "Charlotte"},
		{Title: "Go hiking", Description: "Trail in the mountains", Status: "pending", Assignee: "Alexander"},
		{Title: "Photography", Description: "Cityscape shots", Status: "in_progress", Assignee: "Amelia"},
		{Title: "Learn Spanish", Description: "Basic conversation skills", Status: "done", Assignee: "Benjamin"},
	}

	for _, task := range tasks {
		_, err := db.Conn.Exec(ctx,
			"INSERT INTO tasks (title, description, status, assignee) VALUES ($1, $2, $3, $4)",
			task.Title,
			task.Description,
			task.Status,
			task.Assignee,
		)
		if err != nil {
			panic(err)
		}
	}
}

func postgresDSN(cfg config.Config) string {
	return "postgres://" +
		cfg.DBUser + ":" +
		cfg.DBPassword + "@" +
		cfg.DBHost + ":" +
		cfg.DBPort + "/" +
		cfg.DBName
}
