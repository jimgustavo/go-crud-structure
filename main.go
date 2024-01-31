package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define a struct for your data model.
type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Create a slice to store your tasks (simulating a database):
var tasks []Task

// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	_ = json.NewDecoder(r.Body).Decode(&newTask)
	newTask.ID = fmt.Sprintf("%d", len(tasks)+1)
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// Get all tasks
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get a specific task by ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, task := range tasks {
		if task.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found"))
}

// Update a task by ID
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, task := range tasks {
		if task.ID == params["id"] {
			var updatedTask Task
			_ = json.NewDecoder(r.Body).Decode(&updatedTask)
			updatedTask.ID = task.ID
			tasks[index] = updatedTask
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found"))
}

// Delete a task by ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, task := range tasks {
		if task.ID == params["id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found"))
}

func main() {
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks", getAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*

curl -X POST -H "Content-Type: application/json" -d '{"title":"Task 1","content":"Description 1"}' http://localhost:8080/tasks

curl http://localhost:8080/tasks

# Replace {task_id} with the actual task ID returned from the previous command
curl http://localhost:8080/tasks/{task_id}

# Replace {task_id} with the actual task ID returned from the previous command
curl -X PUT -H "Content-Type: application/json" -d '{"title":"Updated Task","content":"Updated Description"}' http://localhost:8080/tasks/{task_id}

# Replace {task_id} with the actual task ID returned from the previous command
curl -X DELETE http://localhost:8080/tasks/{task_id}


*/
