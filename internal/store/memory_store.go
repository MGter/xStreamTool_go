package store

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/MGter/xStreamTool_go/internal/models"
)

// 定义错误变量
var (
	ErrTodoNotFound = errors.New("待办事项不存在") // 当根据ID找不到待办事项时返回的错误
	ErrInvalidID    = errors.New("无效的ID")   // 当ID格式无效时返回的错误
)

// TodoStore 待办事项存储接口
// 定义了一组操作待办事项数据的接口方法
// 通过接口可以实现不同的存储后端（如内存、数据库等）
type TodoStore interface {
	GetAllTodos() ([]*models.Todo, error)                                               // 获取所有待办事项
	GetTodoByID(id int) (*models.Todo, error)                                           // 根据ID获取单个待办事项
	CreateTodo(req *models.TodoRequest) (*models.Todo, error)                           // 创建新的待办事项
	UpdateTodo(id int, req *models.TodoRequest) (*models.Todo, error)                   // 更新待办事项
	DeleteTodo(id int) error                                                            // 删除待办事项
	SearchTodos(query string, category string, completed *bool) ([]*models.Todo, error) // 搜索待办事项
	GetStats() (map[string]interface{}, error)                                          // 获取待办事项统计信息
}

// MemoryStore 内存存储实现
// 基于内存的待办事项存储实现，使用map存储数据
type MemoryStore struct {
	mu     sync.RWMutex         // 读写锁，用于保证并发安全
	todos  map[int]*models.Todo // 存储待办事项的map，key为ID，value为待办事项对象
	nextID int                  // 下一个可用的ID
}

// NewMemoryStore 创建新的内存存储
func NewMemoryStore() *MemoryStore {
	// 创建MemoryStore实例
	store := &MemoryStore{
		todos:  make(map[int]*models.Todo), // 初始化空的待办事项map
		nextID: 1,                          // 从ID 1开始
	}

	// 初始化示例数据
	store.Seed()
	return store
}

// GetAllTodos 获取所有待办事项
func (s *MemoryStore) GetAllTodos() ([]*models.Todo, error) {
	s.mu.RLock()         // 获取读锁
	defer s.mu.RUnlock() // 函数返回时释放读锁

	// 将map中的所有待办事项转换为切片
	// make： 创建一个切片，长度为当前待办事项数量
	todos := make([]*models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	// 按创建时间倒序排序（最新的在前）
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.After(todos[j].CreatedAt)
	})

	return todos, nil
}

// GetTodoByID 根据ID获取待办事项
func (s *MemoryStore) GetTodoByID(id int) (*models.Todo, error) {
	s.mu.RLock()         // 获取读锁
	defer s.mu.RUnlock() // 函数返回时释放读锁

	// 从map中查找指定ID的待办事项
	todo, exists := s.todos[id]
	if !exists {
		return nil, ErrTodoNotFound // 如果不存在，返回错误
	}

	return todo, nil
}

// CreateTodo 创建新的待办事项
func (s *MemoryStore) CreateTodo(req *models.TodoRequest) (*models.Todo, error) {
	s.mu.Lock()         // 获取写锁
	defer s.mu.Unlock() // 函数返回时释放写锁

	// 获取当前时间
	now := time.Now()

	// 创建新的待办事项对象
	todo := &models.Todo{
		ID:          s.nextID,        // 使用下一个可用的ID
		Title:       req.Title,       // 标题
		Description: req.Description, // 描述
		Completed:   req.Completed,   // 完成状态
		Priority:    req.Priority,    // 优先级
		Category:    req.Category,    // 分类
		DueDate:     req.DueDate,     // 截止日期
		CreatedAt:   now,             // 创建时间
		UpdatedAt:   now,             // 更新时间
	}

	// 将待办事项添加到map中
	s.todos[todo.ID] = todo
	s.nextID++ // ID自增，为下一个待办事项准备

	return todo, nil
}

// UpdateTodo 更新待办事项
func (s *MemoryStore) UpdateTodo(id int, req *models.TodoRequest) (*models.Todo, error) {
	s.mu.Lock()         // 获取写锁
	defer s.mu.Unlock() // 函数返回时释放写锁

	// 查找要更新的待办事项
	todo, exists := s.todos[id]
	if !exists {
		return nil, ErrTodoNotFound // 如果不存在，返回错误
	}

	// 更新待办事项的字段
	todo.FromRequest(req)
	return todo, nil
}

// DeleteTodo 删除待办事项
func (s *MemoryStore) DeleteTodo(id int) error {
	s.mu.Lock()         // 获取写锁
	defer s.mu.Unlock() // 函数返回时释放写锁

	// 检查待办事项是否存在
	if _, exists := s.todos[id]; !exists {
		return ErrTodoNotFound // 如果不存在，返回错误
	}

	// 从map中删除待办事项
	delete(s.todos, id)
	return nil
}

// SearchTodos 搜索待办事项
func (s *MemoryStore) SearchTodos(query string, category string, completed *bool) ([]*models.Todo, error) {
	s.mu.RLock()         // 获取读锁
	defer s.mu.RUnlock() // 函数返回时释放读锁

	// 初始化结果切片
	results := make([]*models.Todo, 0)

	// 遍历所有待办事项，筛选符合条件的
	for _, todo := range s.todos {
		// 匹配查询条件
		matches := true

		// 如果查询字符串不为空，检查标题或描述是否包含该字符串
		if query != "" {
			matches = matches && (strings.Contains(todo.Title, query) || strings.Contains(todo.Description, query))
		}

		// 如果分类不为空，检查分类是否匹配
		if category != "" {
			matches = matches && todo.Category == category
		}

		// 如果completed不为nil，检查完成状态是否匹配
		if completed != nil {
			matches = matches && todo.Completed == *completed
		}

		// 如果所有条件都匹配，添加到结果中
		if matches {
			results = append(results, todo)
		}
	}

	// 按优先级（降序）和创建时间（倒序）排序
	sort.Slice(results, func(i, j int) bool {
		if results[i].Priority != results[j].Priority {
			return results[i].Priority > results[j].Priority // 优先级高的在前
		}
		return results[i].CreatedAt.After(results[j].CreatedAt) // 创建时间晚的在前
	})

	return results, nil
}

// GetStats 获取统计信息
func (s *MemoryStore) GetStats() (map[string]interface{}, error) {
	s.mu.RLock()         // 获取读锁
	defer s.mu.RUnlock() // 函数返回时释放读锁

	// 初始化统计信息map
	stats := map[string]interface{}{
		"total":       len(s.todos),         // 总数量
		"completed":   0,                    // 已完成数量
		"pending":     0,                    // 待完成数量
		"overdue":     0,                    // 已过期数量
		"by_priority": make(map[int]int),    // 按优先级统计
		"by_category": make(map[string]int), // 按分类统计
	}

	// 获取当前时间
	now := time.Now()

	// 遍历所有待办事项，进行统计
	for _, todo := range s.todos {
		if todo.Completed {
			// 已完成的任务
			stats["completed"] = stats["completed"].(int) + 1
		} else {
			// 未完成的任务
			stats["pending"] = stats["pending"].(int) + 1

			// 检查是否过期
			if !todo.DueDate.IsZero() && todo.DueDate.Before(now) {
				stats["overdue"] = stats["overdue"].(int) + 1
			}
		}

		// 按优先级统计
		stats["by_priority"].(map[int]int)[todo.Priority]++

		// 按分类统计
		if todo.Category != "" {
			stats["by_category"].(map[string]int)[todo.Category]++
		}
	}

	return stats, nil
}

// Seed 初始化示例数据
func (s *MemoryStore) Seed() {
	now := time.Now()

	// 创建第一个示例待办事项
	s.todos[1] = &models.Todo{
		ID:          1,
		Title:       "学习 Go 语言",
		Description: "掌握 Go 语言的基础语法和并发编程",
		Completed:   false,
		Priority:    3,
		Category:    "学习",
		DueDate:     now.Add(7 * 24 * time.Hour),
		CreatedAt:   now.Add(-2 * 24 * time.Hour),
		UpdatedAt:   now.Add(-2 * 24 * time.Hour),
	}

	// 创建第二个示例待办事项
	s.todos[2] = &models.Todo{
		ID:          2,
		Title:       "编写 HTTP 服务器",
		Description: "使用 Go 实现一个完整的 HTTP 服务器",
		Completed:   true,
		Priority:    4,
		Category:    "项目",
		DueDate:     now.Add(-1 * 24 * time.Hour),
		CreatedAt:   now.Add(-3 * 24 * time.Hour),
		UpdatedAt:   now.Add(-1 * 24 * time.Hour),
	}

	// 创建第三个示例待办事项
	s.todos[3] = &models.Todo{
		ID:          3,
		Title:       "部署到服务器",
		Description: "将应用部署到生产环境",
		Completed:   false,
		Priority:    2,
		Category:    "运维",
		DueDate:     now.Add(3 * 24 * time.Hour),
		CreatedAt:   now.Add(-1 * 24 * time.Hour),
		UpdatedAt:   now.Add(-1 * 24 * time.Hour),
	}

	// 设置下一个可用的ID为4
	s.nextID = 4
}
