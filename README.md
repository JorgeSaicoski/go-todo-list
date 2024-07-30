
# Go Todo List
This project is a simple RestAPI with a basic frontend, focusing on the backend using Golang. The main objective is to teach and learn Golang as a backend language.

## Features
- **RestAPI**: Provides endpoints to manage a todo list.
- **Frontend**: A simple interface to interact with the API, built using Golang, HTML, and CSS.
- **Database**: Uses SQLite for client-side storage. Users can download their database to use on another computer.
- **Deployment**: Open for deployment on platforms like Azure or Google Cloud.

## Getting Started

### Prerequisites
- Go 1.8 or higher
- SQLite

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/JorgeSaicoski/go-todo-list.git
   cd go-todo-list
   ```
2. Run the application:
   ```bash
   go run main.go
   ```

## API Endpoints
- **Get all todos:**
  ```bash
  curl -H "Content-Type: application/json" -X GET http://localhost:8000/todos
  ```

- **Get a single todo:**
  ```bash
  curl -H "Content-Type: application/json" -X GET http://localhost:8000/todos/{id}
  ```

- **Add a new todo:**
  ```bash
  curl -H "Content-Type: application/json" -X POST -d '{"content":"New Todo"}' http://localhost:8000/todos
  ```

- **Update a todo:**
  ```bash
  curl -H "Content-Type: application/json" -X PUT -d '{"content":"Updated Todo"}' http://localhost:8000/todos/{id}
  ```

- **Delete a todo:**
  ```bash
  curl -H "Content-Type: application/json" -X DELETE http://localhost:8000/todos/{id}
  ```

## Todo List
- [ ] Set up the project structure
- [x] Implement the main.go file
- [ ] Create the route to get all todos
- [ ] Create the route to get a single todo
- [ ] Create the route to add a new todo
- [ ] Create the route to update a todo
- [ ] Create the route to delete a todo
- [ ] Build the frontend interface using Golang, HTML, and CSS
- [ ] Test the application
- [ ] Prepare for deployment on Azure or Google Cloud

## Deployment
This project can be deployed on platforms like Azure or Google Cloud. Ensure that the necessary configurations are made for the respective platforms.

## Contributing
Feel free to fork this repository and contribute by submitting a pull request. For major changes, please open an issue first to discuss what you would like to change. It is recommended to create a version with full server-side implementation.

## License
This project is licensed under the MIT License.
```

You can copy and paste this Markdown into your README file. Let me know if you need any more changes or additional sections!


