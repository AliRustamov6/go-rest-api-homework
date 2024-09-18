package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// getTask возвращает все задачи
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Ошибка ", http.StatusInternalServerError)
		return
	}

}

// postTask создает новую запись
func postTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	// прверка, существует ли задача с таким ID
	if _, ok := tasks[newTask.ID]; ok {
		http.Error(w, "Задача с таким ID уже существует", http.StatusConflict)
		return
	}

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

// getTaskID добавляет задачу по ID
func getTaskID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	// Если задача не найдена, возвращаем 400 Bad Request
	if !ok {
		http.Error(w, "Задача с таким ID не найдена", http.StatusBadRequest)
		return
	}

	// если задача ненайдена, отправляем её вответ
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при оброботке запроса", http.StatusInternalServerError)
		return
	}

}

// deleteTaskID удаляет задачу по ID
func deleteTaskID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", getTask)              // получение всех задач
	r.Post("/tasks", postTask)            // Создание новых задач
	r.Get("/tasks/{id}", getTaskID)       // Получение задач по ID
	r.Delete("/tasks/{id}", deleteTaskID) // Удаление задачи по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
