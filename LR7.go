package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type User struct {
	Username string `json:"username"`
}

var users = make(map[string]User)

func main() {
	// Обработчик для статической HTML-страницы
	http.HandleFunc("/", serveStaticHTML)

	// Обработчик для управления пользователями
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUser)

	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Функция для обработки запросов к статическому HTML
func serveStaticHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>User Management</title>
</head>
<body>
    <h1>User Management</h1>

    <!-- Просмотр пользователей -->
    <button onclick="viewUsers()">View Users</button>
    <div id="userList"></div>

    <!-- Добавление пользователя -->
    <h2>Add User</h2>
    <input type="text" id="username" placeholder="Enter username">
    <button onclick="addUser()">Add User</button>

    <!-- Удаление пользователя -->
    <h2>Delete User</h2>
    <input type="text" id="deleteUsername" placeholder="Enter username to delete">
    <button onclick="deleteUser()">Delete User</button>

    <script>
        function viewUsers() {
            fetch('/users')
                .then(response => response.json())
                .then(data => {
                    let userList = '<h3>Users:</h3><ul>';
                    data.forEach(user => {
                        userList += '<li>' + user.username + '</li>';
                    });
                    userList += '</ul>';
                    document.getElementById('userList').innerHTML = userList;
                })
                .catch(error => console.error('Error:', error));
        }

        function addUser() {
            const username = document.getElementById('username').value;
            fetch('/users', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: username })
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                viewUsers();
            })
            .catch(error => console.error('Error:', error));
        }

        function deleteUser() {
            const username = document.getElementById('deleteUsername').value;
            fetch('/users/' + username, {
                method: 'DELETE'
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                viewUsers();
            })
            .catch(error => console.error('Error:', error));
        }
    </script>
</body>
</html>
`
	w.Write([]byte(html))
}

// Обработчик для работы с коллекцией пользователей (GET, POST)
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w)
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Обработчик для работы с конкретным пользователем (DELETE)
func handleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deleteUser(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Получение списка пользователей
func getUsers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	// Сортировка пользователей по имени
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].Username < userList[j].Username
	})

	json.NewEncoder(w).Encode(userList)
}

// Добавление пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if newUser.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	users[newUser.Username] = newUser
	fmt.Fprintf(w, "User %s added", newUser.Username)
}

// Удаление пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/users/")
	if _, exists := users[username]; exists {
		delete(users, username)
		fmt.Fprintf(w, "User %s deleted", username)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
