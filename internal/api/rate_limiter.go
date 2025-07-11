package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter реализует token bucket алгоритм для ограничения частоты запросов
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     time.Duration // Интервал между пополнениями токенов
	capacity int           // Максимальное количество токенов
	cleanup  time.Duration // Интервал очистки неактивных посетителей
}

// Visitor представляет посетителя с его bucket токенов
type Visitor struct {
	tokens     int
	lastSeen   time.Time
	capacity   int
	refillRate time.Duration
}

// NewRateLimiter создает новый rate limiter
// rate - сколько запросов в секунду разрешено
// capacity - максимальный burst размер
func NewRateLimiter(requestsPerSecond int, burstCapacity int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     time.Second / time.Duration(requestsPerSecond),
		capacity: burstCapacity,
		cleanup:  time.Minute * 5, // Очищаем неактивных посетителей каждые 5 минут
	}

	// Запускаем фоновую очистку
	go rl.cleanupRoutine()

	return rl
}

// Allow проверяет, разрешен ли запрос для данного IP
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor := rl.getVisitor(ip)
	return visitor.allowRequest()
}

// getVisitor получает или создает посетителя для IP
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			tokens:     rl.capacity,
			lastSeen:   time.Now(),
			capacity:   rl.capacity,
			refillRate: rl.rate,
		}
		rl.visitors[ip] = visitor
	}

	visitor.lastSeen = time.Now()
	return visitor
}

// allowRequest проверяет и потребляет токен если возможно
func (v *Visitor) allowRequest() bool {
	now := time.Now()

	// Пополняем токены на основе прошедшего времени
	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed / v.refillRate)

	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > v.capacity {
			v.tokens = v.capacity
		}
		v.lastSeen = now
	}

	// Проверяем, есть ли доступные токены
	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// cleanupRoutine периодически удаляет неактивных посетителей
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupInactive()
	}
}

// cleanupInactive удаляет посетителей, которые не активны более часа
func (rl *RateLimiter) cleanupInactive() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-time.Hour)
	for ip, visitor := range rl.visitors {
		if visitor.lastSeen.Before(cutoff) {
			delete(rl.visitors, ip)
		}
	}
}

// Middleware возвращает HTTP middleware для rate limiting
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := rl.getClientIP(r)

			if !rl.Allow(ip) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP извлекает IP адрес клиента из запроса
func (rl *RateLimiter) getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Используем RemoteAddr как fallback
	return r.RemoteAddr
}

// GetStats возвращает статистику rate limiter для мониторинга
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"active_visitors": len(rl.visitors),
		"rate_per_second": int(time.Second / rl.rate),
		"burst_capacity":  rl.capacity,
	}
}

// CommentRateLimiter специальный rate limiter для создания комментариев
type CommentRateLimiter struct {
	*RateLimiter
	perPostLimits map[string]*PostLimiter
	mu            sync.RWMutex
}

// PostLimiter ограничивает комментарии к конкретному посту
type PostLimiter struct {
	postID   string
	lastSeen time.Time
	count    int
	limit    int
	window   time.Duration
}

// NewCommentRateLimiter создает rate limiter специально для комментариев
func NewCommentRateLimiter() *CommentRateLimiter {
	return &CommentRateLimiter{
		RateLimiter:   NewRateLimiter(10, 20), // 10 запросов в секунду, burst 20
		perPostLimits: make(map[string]*PostLimiter),
	}
}

// AllowComment проверяет, можно ли создать комментарий
func (crl *CommentRateLimiter) AllowComment(ip, postID string) bool {
	// Сначала проверяем общий rate limit
	if !crl.Allow(ip) {
		return false
	}

	// Затем проверяем лимит на пост
	return crl.allowCommentToPost(ip, postID)
}

// allowCommentToPost проверяет лимит комментариев к конкретному посту
func (crl *CommentRateLimiter) allowCommentToPost(ip, postID string) bool {
	crl.mu.Lock()
	defer crl.mu.Unlock()

	key := fmt.Sprintf("%s:%s", ip, postID)
	limiter, exists := crl.perPostLimits[key]

	if !exists {
		limiter = &PostLimiter{
			postID:   postID,
			lastSeen: time.Now(),
			count:    0,
			limit:    5,                // Максимум 5 комментариев к одному посту
			window:   time.Minute * 10, // За 10 минут
		}
		crl.perPostLimits[key] = limiter
	}

	now := time.Now()

	// Сбрасываем счетчик если прошло время окна
	if now.Sub(limiter.lastSeen) > limiter.window {
		limiter.count = 0
		limiter.lastSeen = now
	}

	// Проверяем лимит
	if limiter.count >= limiter.limit {
		return false
	}

	limiter.count++
	limiter.lastSeen = now
	return true
}

// GraphQLRateLimiter специальный rate limiter для GraphQL запросов
type GraphQLRateLimiter struct {
	*RateLimiter
	complexityLimits map[string]int // Лимиты по сложности запросов
}

// NewGraphQLRateLimiter создает rate limiter для GraphQL
func NewGraphQLRateLimiter() *GraphQLRateLimiter {
	return &GraphQLRateLimiter{
		RateLimiter: NewRateLimiter(30, 50), // 30 запросов в секунду, burst 50
		complexityLimits: map[string]int{
			"Query":        1000, // Максимальная сложность Query
			"Mutation":     500,  // Максимальная сложность Mutation
			"Subscription": 100,  // Максимальная сложность Subscription
		},
	}
}

// AllowGraphQLRequest проверяет GraphQL запрос с учетом сложности
func (grl *GraphQLRateLimiter) AllowGraphQLRequest(ip, operationType string, complexity int) bool {
	// Проверяем общий rate limit
	if !grl.Allow(ip) {
		return false
	}

	// Проверяем лимит сложности
	maxComplexity, exists := grl.complexityLimits[operationType]
	if exists && complexity > maxComplexity {
		return false
	}

	return true
}
