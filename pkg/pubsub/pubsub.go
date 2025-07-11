package pubsub

import (
	"sync"
)

// Константы конфигурации по умолчанию
const (
	// DefaultChannelBufferSize - размер буфера канала подписчика по умолчанию
	DefaultChannelBufferSize = 100
)

// Message представляет сообщение в системе Pub/Sub.
// Содержит топик и данные для передачи подписчикам.
type Message struct {
	Topic string      `json:"topic"`          // Название топика
	Data  interface{} `json:"data,omitempty"` // Данные сообщения
}

// Subscriber представляет подписчика на топик.
// Каждый подписчик имеет уникальный ID и канал для получения сообщений.
type Subscriber struct {
	ID      string       `json:"id"` // Уникальный идентификатор подписчика
	Channel chan Message `json:"-"`  // Канал для получения сообщений (не сериализуется)
}

// PubSub представляет простую in-memory систему публикации/подписки.
// Поддерживает множественные топики и подписчиков с thread-safe операциями.
//
// Основные возможности:
// - Thread-safe операции подписки/отписки
// - Буферизированные каналы для предотвращения блокировок
// - Автоматическая очистка пустых топиков
// - Неблокирующая публикация сообщений
type PubSub struct {
	mu                sync.RWMutex                      // Мьютекс для thread-safe операций
	subscribers       map[string]map[string]*Subscriber // topic -> subscriberID -> subscriber
	channelBufferSize int                               // Размер буфера для каналов подписчиков
}

// New создает новый экземпляр PubSub с буфером по умолчанию.
func New() *PubSub {
	return NewWithConfig(DefaultChannelBufferSize)
}

// NewWithConfig создает новый экземпляр PubSub с указанным размером буфера канала.
// channelBufferSize определяет размер буфера для каналов подписчиков.
// Больший буфер снижает вероятность потери сообщений при медленных подписчиках.
func NewWithConfig(channelBufferSize int) *PubSub {
	if channelBufferSize <= 0 {
		channelBufferSize = DefaultChannelBufferSize
	}

	return &PubSub{
		subscribers:       make(map[string]map[string]*Subscriber),
		channelBufferSize: channelBufferSize,
	}
}

// Subscribe подписывает клиента на топик.
// Создает новый канал для подписчика с настроенным размером буфера.
// Если топик не существует, он создается автоматически.
//
// Параметры:
//   - topic: название топика для подписки
//   - subscriberID: уникальный идентификатор подписчика
//
// Возвращает Subscriber с каналом для получения сообщений.
func (ps *PubSub) Subscribe(topic string, subscriberID string) *Subscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Создаем топик если он не существует
	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[string]*Subscriber)
	}

	subscriber := &Subscriber{
		ID:      subscriberID,
		Channel: make(chan Message, ps.channelBufferSize),
	}

	ps.subscribers[topic][subscriberID] = subscriber
	return subscriber
}

// Unsubscribe отписывает клиента от топика.
// Закрывает канал подписчика и удаляет его из списка.
// Если топик остается без подписчиков, он удаляется для экономии памяти.
//
// Параметры:
//   - topic: название топика
//   - subscriberID: идентификатор подписчика для отписки
func (ps *PubSub) Unsubscribe(topic string, subscriberID string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	topicSubs, topicExists := ps.subscribers[topic]
	if !topicExists {
		return
	}

	sub, subExists := topicSubs[subscriberID]
	if subExists {
		close(sub.Channel)
		delete(topicSubs, subscriberID)
	}

	// Удаляем топик если нет подписчиков (экономия памяти)
	if len(topicSubs) == 0 {
		delete(ps.subscribers, topic)
	}
}

// Publish публикует сообщение в топик всем подписчикам.
// Использует неблокирующую отправку для предотвращения deadlock'ов.
// Если канал подписчика переполнен, сообщение пропускается.
//
// Параметры:
//   - topic: название топика для публикации
//   - data: данные для отправки подписчикам
//
// Операция thread-safe и не блокирует при переполненных каналах.
func (ps *PubSub) Publish(topic string, data interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	message := Message{
		Topic: topic,
		Data:  data,
	}

	topicSubs, exists := ps.subscribers[topic]
	if !exists {
		return // Топик не существует или нет подписчиков
	}

	// Отправляем сообщение всем подписчикам топика
	for _, subscriber := range topicSubs {
		// Неблокирующая отправка - защита от медленных подписчиков
		select {
		case subscriber.Channel <- message:
			// Сообщение успешно отправлено
		default:
			// Канал переполнен, пропускаем сообщение
			// В production можно добавить логирование или метрики
		}
	}
}

// GetSubscribersCount возвращает количество подписчиков на топик.
// Используется для мониторинга и health check.
//
// Параметры:
//   - topic: название топика
//
// Возвращает количество активных подписчиков.
func (ps *PubSub) GetSubscribersCount(topic string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	topicSubs, exists := ps.subscribers[topic]
	if !exists {
		return 0
	}
	return len(topicSubs)
}

// Close закрывает все каналы подписчиков и очищает внутренние структуры.
// Используется при завершении работы приложения для корректной очистки ресурсов.
func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Закрываем все каналы подписчиков
	for _, topicSubs := range ps.subscribers {
		for _, subscriber := range topicSubs {
			close(subscriber.Channel)
		}
	}

	// Очищаем все структуры данных
	ps.subscribers = make(map[string]map[string]*Subscriber)
}
