package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MGter/xStreamTool_go/internal/models"
	"github.com/MGter/xStreamTool_go/internal/store"
	"github.com/gorilla/mux"
)

// Handler HTTP 处理器
type Handler struct {
	store store.TodoStore
}

// NewHandler 创建新的处理器
func NewHandler(store store.TodoStore) *Handler {
	return &Handler{store: store}
}

// SetupRoutes 设置路由
func SetupRoutes(h *Handler) *mux.Router {
	router := mux.NewRouter()

	// 全局中间件
	router.Use(loggingMiddleware)

	// Web 页面路由
	router.HandleFunc("/", h.HomePage).Methods("GET")
	router.HandleFunc("/todos", h.TodosPage).Methods("GET")
	router.HandleFunc("/api/docs", h.APIDocsPage).Methods("GET")

	// API 路由
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/todos", h.GetTodos).Methods("GET")
	api.HandleFunc("/todos", h.CreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id}", h.GetTodo).Methods("GET")
	api.HandleFunc("/todos/{id}", h.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id}", h.DeleteTodo).Methods("DELETE")
	api.HandleFunc("/todos/{id}/complete", h.CompleteTodo).Methods("PATCH")
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")

	return router
}

// HomePage 首页
func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>xStreamTool Go</title>
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
			h1 { color: #333; }
			.card { background: #f9f9f9; padding: 20px; margin: 20px 0; border-radius: 8px; }
			.btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
		</style>
	</head>
	<body>
		<h1>🚀 xStreamTool Go HTTP 服务器</h1>
		<div class="card">
			<h2>欢迎使用</h2>
			<p>这是一个简单的 Go HTTP 服务器示例</p>
			<a href="/todos" class="btn">查看待办事项</a>
			<a href="/api/docs" class="btn">API 文档</a>
		</div>
		<div class="card">
			<h3>📋 API 端点</h3>
			<ul>
				<li><code>GET /api/todos</code> - 获取所有待办事项</li>
				<li><code>GET /api/todos/{id}</code> - 获取单个待办事项</li>
				<li><code>POST /api/todos</code> - 创建新待办事项</li>
				<li><code>PUT /api/todos/{id}</code> - 更新待办事项</li>
				<li><code>DELETE /api/todos/{id}</code> - 删除待办事项</li>
			</ul>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// TodosPage 待办事项页面
func (h *Handler) TodosPage(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetAllTodos()
	if err != nil {
		sendError(w, "获取待办事项失败", http.StatusInternalServerError)
		return
	}

	tmplStr := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>待办事项</title>
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
			.todo-item { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
			.completed { background: #e8f5e8; }
			.btn { padding: 5px 10px; margin-right: 5px; border: none; border-radius: 3px; cursor: pointer; }
			.btn-primary { background: #007bff; color: white; }
			.btn-success { background: #28a745; color: white; }
			.btn-danger { background: #dc3545; color: white; }
		</style>
	</head>
	<body>
		<h1>📋 待办事项列表</h1>
		<div id="todoList">
			{{range .}}
			<div class="todo-item {{if .Completed}}completed{{end}}">
				<h3>{{.Title}} {{if .Completed}}✅{{end}}</h3>
				<p>ID: {{.ID}} | 创建时间: {{.CreatedAt.Format "2006-01-02 15:04"}}</p>
				<p>优先级: {{.Priority}} | 分类: {{.Category}}</p>
				<button class="btn btn-success" onclick="completeTodo({{.ID}})">标记完成</button>
				<button class="btn btn-danger" onclick="deleteTodo({{.ID}})">删除</button>
			</div>
			{{else}}
			<p>暂无待办事项</p>
			{{end}}
		</div>
		
		<div style="margin-top: 30px; background: #f8f9fa; padding: 20px; border-radius: 8px;">
			<h3>添加新待办事项</h3>
			<input type="text" id="title" placeholder="标题" style="width: 100%; padding: 10px; margin: 10px 0;">
			<textarea id="description" placeholder="描述" style="width: 100%; padding: 10px; margin: 10px 0;" rows="3"></textarea>
			<button class="btn btn-primary" onclick="createTodo()">添加</button>
		</div>

		<script>
			async function createTodo() {
				const title = document.getElementById('title').value;
				if (!title) {
					alert('请输入标题');
					return;
				}
				
				const response = await fetch('/api/todos', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ title: title, description: document.getElementById('description').value })
				});
				
				if (response.ok) {
					alert('创建成功！');
					location.reload();
				}
			}
			
			async function completeTodo(id) {
				const response = await fetch('/api/todos/' + id + '/complete', { method: 'PATCH' });
				if (response.ok) {
					alert('标记完成！');
					location.reload();
				}
			}
			
			async function deleteTodo(id) {
				if (!confirm('确定删除吗？')) return;
				const response = await fetch('/api/todos/' + id, { method: 'DELETE' });
				if (response.ok) {
					alert('删除成功！');
					location.reload();
				}
			}
		</script>
	</body>
	</html>
	`

	tmpl, err := template.New("todos").Parse(tmplStr)
	if err != nil {
		sendError(w, "模板错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, todos)
}

// APIDocsPage API 文档页面
func (h *Handler) APIDocsPage(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>API 文档</title>
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
			.endpoint { background: #f8f9fa; padding: 15px; margin: 15px 0; border-radius: 5px; }
			.method { display: inline-block; padding: 5px 10px; background: #6c757d; color: white; border-radius: 3px; }
			.path { font-family: monospace; background: #e9ecef; padding: 5px; border-radius: 3px; }
		</style>
	</head>
	<body>
		<h1>📚 API 文档</h1>
		<div class="endpoint">
			<span class="method">GET</span> <span class="path">/api/todos</span>
			<p>获取所有待办事项</p>
		</div>
		<div class="endpoint">
			<span class="method">POST</span> <span class="path">/api/todos</span>
			<p>创建待办事项</p>
			<pre>{
  "title": "任务标题",
  "description": "任务描述"
}</pre>
		</div>
		<div class="endpoint">
			<span class="method">GET</span> <span class="path">/api/todos/{id}</span>
			<p>获取单个待办事项</p>
		</div>
		<div class="endpoint">
			<span class="method">PUT</span> <span class="path">/api/todos/{id}</span>
			<p>更新待办事项</p>
		</div>
		<div class="endpoint">
			<span class="method">DELETE</span> <span class="path">/api/todos/{id}</span>
			<p>删除待办事项</p>
		</div>
		<div class="endpoint">
			<span class="method">PATCH</span> <span class="path">/api/todos/{id}/complete</span>
			<p>标记待办事项为完成</p>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// GetTodos 获取所有待办事项
func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetAllTodos()
	if err != nil {
		sendError(w, "获取失败", http.StatusInternalServerError)
		return
	}

	responses := make([]models.TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = todo.ToResponse()
	}

	sendJSON(w, responses, http.StatusOK)
}

// GetTodo 获取单个待办事项
func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, "无效ID", http.StatusBadRequest)
		return
	}

	todo, err := h.store.GetTodoByID(id)
	if err != nil {
		sendError(w, "未找到", http.StatusNotFound)
		return
	}

	sendJSON(w, todo.ToResponse(), http.StatusOK)
}

// CreateTodo 创建待办事项
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "无效数据", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		sendError(w, "标题必填", http.StatusBadRequest)
		return
	}

	todo, err := h.store.CreateTodo(&req)
	if err != nil {
		sendError(w, "创建失败", http.StatusInternalServerError)
		return
	}

	sendJSON(w, todo.ToResponse(), http.StatusCreated)
}

// UpdateTodo 更新待办事项
func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, "无效ID", http.StatusBadRequest)
		return
	}

	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "无效数据", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		sendError(w, "标题必填", http.StatusBadRequest)
		return
	}

	todo, err := h.store.UpdateTodo(id, &req)
	if err != nil {
		sendError(w, "更新失败", http.StatusNotFound)
		return
	}

	sendJSON(w, todo.ToResponse(), http.StatusOK)
}

// DeleteTodo 删除待办事项
func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, "无效ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteTodo(id); err != nil {
		sendError(w, "删除失败", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]string{"message": "删除成功"}, http.StatusOK)
}

// CompleteTodo 标记完成
func (h *Handler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, "无效ID", http.StatusBadRequest)
		return
	}

	todo, err := h.store.GetTodoByID(id)
	if err != nil {
		sendError(w, "未找到", http.StatusNotFound)
		return
	}

	req := &models.TodoRequest{
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   true,
		Priority:    todo.Priority,
		Category:    todo.Category,
		DueDate:     todo.DueDate,
	}

	updatedTodo, err := h.store.UpdateTodo(id, req)
	if err != nil {
		sendError(w, "更新失败", http.StatusInternalServerError)
		return
	}

	sendJSON(w, updatedTodo.ToResponse(), http.StatusOK)
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().Unix(),
		"service": "xstreamtool-go",
		"version": "1.0.0",
	}
	sendJSON(w, response, http.StatusOK)
}

// 辅助函数
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON编码错误: %v", err)
	}
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	sendJSON(w, map[string]string{"error": message}, statusCode)
}

// 中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
