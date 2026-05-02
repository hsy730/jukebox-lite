package repository

import (
	"fmt"
	"jukebox-lite/model"
	"sort"
	"sync"
	"time"
)

type MemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
	seq    int64
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (r *MemoryOrderRepository) Create(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	order.ID = fmt.Sprintf("ORD%08d", r.seq)
	order.CreatedAt = time.Now().Unix()
	order.Status = "pending"
	r.orders[order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) GetByID(id string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found: %s", id)
	}
	return order, nil
}

func (r *MemoryOrderRepository) ListBySingerID(singerID string, page, limit int) ([]*model.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.Order
	for _, o := range r.orders {
		if o.SingerID == singerID {
			result = append(result, o)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	total := int64(len(result))
	start := (page - 1) * limit
	if start >= len(result) {
		return []*model.Order{}, total, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (r *MemoryOrderRepository) ListByUserOpenID(openID string, page, limit int) ([]*model.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.Order
	for _, o := range r.orders {
		if o.UserOpenID == openID {
			result = append(result, o)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	total := int64(len(result))
	start := (page - 1) * limit
	if start >= len(result) {
		return []*model.Order{}, total, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (r *MemoryOrderRepository) UpdateStatus(id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[id]
	if !ok {
		return fmt.Errorf("order not found: %s", id)
	}
	order.Status = status
	return nil
}

type MemorySingerRepository struct {
	mu      sync.RWMutex
	singers map[string]*model.Singer
}

func NewMemorySingerRepository() *MemorySingerRepository {
	repo := &MemorySingerRepository{
		singers: make(map[string]*model.Singer),
	}
	repo.initSeed()
	return repo
}

func (r *MemorySingerRepository) initSeed() {
	seed := []*model.Singer{
		{ID: "singer_001", Name: "小明", Avatar: "https://img.zcool.cn/community/01786557e4a6fa0000018c1bf080ca.png", OpenID: "", Status: "online"},
		{ID: "singer_002", Name: "小红", Avatar: "https://img.zcool.cn/community/01786557e4a6fa0000018c1bf080ca.png", OpenID: "", Status: "online"},
		{ID: "singer_003", Name: "阿杰", Avatar: "https://img.zcool.cn/community/01786557e4a6fa0000018c1bf080ca.png", OpenID: "", Status: "offline"},
	}
	for _, s := range seed {
		r.singers[s.ID] = s
	}
}

func (r *MemorySingerRepository) Create(singer *model.Singer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.singers[singer.ID] = singer
	return nil
}

func (r *MemorySingerRepository) GetByID(id string) (*model.Singer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	singer, ok := r.singers[id]
	if !ok {
		return nil, fmt.Errorf("singer not found: %s", id)
	}
	return singer, nil
}

func (r *MemorySingerRepository) List(page, limit int) ([]*model.Singer, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.Singer
	for _, s := range r.singers {
		result = append(result, s)
	}
	total := int64(len(result))
	start := (page - 1) * limit
	if start >= len(result) {
		return []*model.Singer{}, total, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (r *MemorySingerRepository) UpdateStatus(id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	singer, ok := r.singers[id]
	if !ok {
		return fmt.Errorf("singer not found: %s", id)
	}
	singer.Status = status
	return nil
}

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *MemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.OpenID] = user
	return nil
}

func (r *MemoryUserRepository) GetByOpenID(openID string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[openID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", openID)
	}
	return user, nil
}

func (r *MemoryUserRepository) UpdateRole(openID string, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[openID]
	if !ok {
		return fmt.Errorf("user not found: %s", openID)
	}
	user.Role = role
	return nil
}
