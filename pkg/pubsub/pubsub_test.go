package pubsub

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ps := New()

	if ps == nil {
		t.Fatal("New() returned nil")
	}

	if ps.channelBufferSize != DefaultChannelBufferSize {
		t.Errorf("Expected buffer size %d, got %d", DefaultChannelBufferSize, ps.channelBufferSize)
	}

	if ps.subscribers == nil {
		t.Error("subscribers map should be initialized")
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name               string
		bufferSize         int
		expectedBufferSize int
	}{
		{
			name:               "позитивный размер буфера",
			bufferSize:         50,
			expectedBufferSize: 50,
		},
		{
			name:               "нулевой размер буфера",
			bufferSize:         0,
			expectedBufferSize: DefaultChannelBufferSize,
		},
		{
			name:               "отрицательный размер буфера",
			bufferSize:         -10,
			expectedBufferSize: DefaultChannelBufferSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewWithConfig(tt.bufferSize)

			if ps.channelBufferSize != tt.expectedBufferSize {
				t.Errorf("Expected buffer size %d, got %d", tt.expectedBufferSize, ps.channelBufferSize)
			}
		})
	}
}

func TestPubSub_Subscribe(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"
	subscriberID := "subscriber-1"

	subscriber := ps.Subscribe(topic, subscriberID)

	if subscriber == nil {
		t.Fatal("Subscribe returned nil subscriber")
	}

	if subscriber.ID != subscriberID {
		t.Errorf("Expected subscriber ID %s, got %s", subscriberID, subscriber.ID)
	}

	if subscriber.Channel == nil {
		t.Error("Subscriber channel should not be nil")
	}

	if cap(subscriber.Channel) != 10 {
		t.Errorf("Expected channel buffer size 10, got %d", cap(subscriber.Channel))
	}

	// Проверяем, что подписчик добавлен в топик
	count := ps.GetSubscribersCount(topic)
	if count != 1 {
		t.Errorf("Expected 1 subscriber, got %d", count)
	}
}

func TestPubSub_Subscribe_MultipleSubscribers(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"

	// Подписываем нескольких подписчиков
	subscriber1 := ps.Subscribe(topic, "subscriber-1")
	subscriber2 := ps.Subscribe(topic, "subscriber-2")

	if subscriber1.ID == subscriber2.ID {
		t.Error("Subscribers should have different IDs")
	}

	count := ps.GetSubscribersCount(topic)
	if count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"
	subscriberID := "subscriber-1"

	// Подписываемся
	subscriber := ps.Subscribe(topic, subscriberID)

	// Проверяем, что подписчик добавлен
	if ps.GetSubscribersCount(topic) != 1 {
		t.Error("Subscriber should be added")
	}

	// Отписываемся
	ps.Unsubscribe(topic, subscriberID)

	// Проверяем, что подписчик удален
	if ps.GetSubscribersCount(topic) != 0 {
		t.Error("Subscriber should be removed")
	}

	// Проверяем, что канал закрыт
	select {
	case _, ok := <-subscriber.Channel:
		if ok {
			t.Error("Channel should be closed")
		}
	default:
		// Канал может быть закрыт без доступных сообщений
	}
}

func TestPubSub_Unsubscribe_NonExistentTopic(t *testing.T) {
	ps := NewWithConfig(10)

	// Попытка отписаться от несуществующего топика не должна вызвать панику
	ps.Unsubscribe("non-existent-topic", "subscriber-1")

	// Тест проходит, если не было паники
}

func TestPubSub_Unsubscribe_NonExistentSubscriber(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"

	// Создаем топик с одним подписчиком
	ps.Subscribe(topic, "subscriber-1")

	// Попытка отписать несуществующего подписчика
	ps.Unsubscribe(topic, "non-existent-subscriber")

	// Проверяем, что существующий подписчик остался
	if ps.GetSubscribersCount(topic) != 1 {
		t.Error("Existing subscriber should remain")
	}
}

func TestPubSub_Publish(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"
	testData := "test message"

	// Подписываемся
	subscriber := ps.Subscribe(topic, "subscriber-1")

	// Публикуем сообщение
	ps.Publish(topic, testData)

	// Проверяем, что сообщение получено
	select {
	case msg := <-subscriber.Channel:
		if msg.Topic != topic {
			t.Errorf("Expected topic %s, got %s", topic, msg.Topic)
		}
		if msg.Data != testData {
			t.Errorf("Expected data %v, got %v", testData, msg.Data)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Message not received within timeout")
	}
}

func TestPubSub_Publish_MultipleSubscribers(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"
	testData := "test message"

	// Подписываем несколько подписчиков
	subscriber1 := ps.Subscribe(topic, "subscriber-1")
	subscriber2 := ps.Subscribe(topic, "subscriber-2")

	// Публикуем сообщение
	ps.Publish(topic, testData)

	// Проверяем, что оба подписчика получили сообщение
	for i, subscriber := range []*Subscriber{subscriber1, subscriber2} {
		select {
		case msg := <-subscriber.Channel:
			if msg.Data != testData {
				t.Errorf("Subscriber %d: expected data %v, got %v", i+1, testData, msg.Data)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Subscriber %d: message not received within timeout", i+1)
		}
	}
}

func TestPubSub_Publish_NonExistentTopic(t *testing.T) {
	ps := NewWithConfig(10)

	// Публикация в несуществующий топик не должна вызвать панику
	ps.Publish("non-existent-topic", "test data")

	// Тест проходит, если не было паники
}

func TestPubSub_GetSubscribersCount(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"

	// Проверяем нулевой счетчик для несуществующего топика
	if count := ps.GetSubscribersCount("non-existent-topic"); count != 0 {
		t.Errorf("Expected 0 subscribers for non-existent topic, got %d", count)
	}

	// Добавляем подписчиков
	ps.Subscribe(topic, "subscriber-1")
	ps.Subscribe(topic, "subscriber-2")

	if count := ps.GetSubscribersCount(topic); count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}

	// Удаляем одного подписчика
	ps.Unsubscribe(topic, "subscriber-1")

	if count := ps.GetSubscribersCount(topic); count != 1 {
		t.Errorf("Expected 1 subscriber after unsubscribe, got %d", count)
	}
}

func TestPubSub_Close(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"

	// Создаем нескольких подписчиков
	subscriber1 := ps.Subscribe(topic, "subscriber-1")
	subscriber2 := ps.Subscribe(topic, "subscriber-2")

	// Закрываем PubSub
	ps.Close()

	// Проверяем, что все каналы закрыты
	for i, subscriber := range []*Subscriber{subscriber1, subscriber2} {
		select {
		case _, ok := <-subscriber.Channel:
			if ok {
				t.Errorf("Subscriber %d: channel should be closed", i+1)
			}
		default:
			// Канал может быть закрыт без доступных сообщений
		}
	}

	// Проверяем, что структуры данных очищены
	if len(ps.subscribers) != 0 {
		t.Error("Subscribers map should be empty after close")
	}
}

func TestPubSub_ThreadSafety(t *testing.T) {
	ps := NewWithConfig(100)
	topic := "test-topic"

	var wg sync.WaitGroup
	subscriberCount := 5
	messageCount := 50

	// Создаем подписчиков сначала
	subscribers := make([]*Subscriber, subscriberCount)
	for i := 0; i < subscriberCount; i++ {
		subscriberID := fmt.Sprintf("subscriber-%d", i)
		subscribers[i] = ps.Subscribe(topic, subscriberID)
	}

	// Запускаем горутины для чтения
	for i := 0; i < subscriberCount; i++ {
		wg.Add(1)
		go func(id int, sub *Subscriber) {
			defer wg.Done()
			messagesReceived := 0

			for messagesReceived < messageCount {
				select {
				case <-sub.Channel:
					messagesReceived++
				case <-time.After(2 * time.Second):
					t.Errorf("Subscriber %d: timeout waiting for message %d", id, messagesReceived)
					return
				}
			}
		}(i, subscribers[i])
	}

	// Небольшая задержка для готовности подписчиков
	time.Sleep(10 * time.Millisecond)

	// Запускаем горутину для публикации
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < messageCount; i++ {
			ps.Publish(topic, fmt.Sprintf("message-%d", i))
			time.Sleep(500 * time.Microsecond) // Небольшая задержка
		}
	}()

	// Ждем завершения всех горутин
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Все горутины завершились успешно
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - possible deadlock or race condition")
	}
}

func TestPubSub_BufferOverflow(t *testing.T) {
	bufferSize := 2
	ps := NewWithConfig(bufferSize)
	topic := "test-topic"

	subscriber := ps.Subscribe(topic, "subscriber-1")

	// Публикуем больше сообщений, чем размер буфера
	for i := 0; i < bufferSize+5; i++ {
		ps.Publish(topic, fmt.Sprintf("message-%d", i))
	}

	// Проверяем, что система не зависла и сообщения в буфере доступны
	receivedCount := 0
	for {
		select {
		case <-subscriber.Channel:
			receivedCount++
		case <-time.After(10 * time.Millisecond):
			// Больше сообщений не поступает
			goto checkReceived
		}
	}

checkReceived:
	// Должны получить не больше, чем размер буфера (из-за неблокирующей отправки)
	if receivedCount > bufferSize {
		t.Errorf("Received %d messages, expected at most %d (buffer size)", receivedCount, bufferSize)
	}
}

func TestPubSub_TopicCleanup(t *testing.T) {
	ps := NewWithConfig(10)
	topic := "test-topic"

	// Подписываемся и сразу отписываемся
	ps.Subscribe(topic, "subscriber-1")
	ps.Unsubscribe(topic, "subscriber-1")

	// Проверяем, что топик удален (нет подписчиков)
	count := ps.GetSubscribersCount(topic)
	if count != 0 {
		t.Errorf("Expected 0 subscribers after cleanup, got %d", count)
	}

	// Проверяем, что топик действительно удален из внутренней структуры
	ps.mu.RLock()
	_, exists := ps.subscribers[topic]
	ps.mu.RUnlock()

	if exists {
		t.Error("Topic should be removed from internal structure when no subscribers remain")
	}
}

// Benchmark тесты для производительности
func BenchmarkPubSub_Subscribe(b *testing.B) {
	ps := NewWithConfig(100)
	topic := "benchmark-topic"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Subscribe(topic, fmt.Sprintf("subscriber-%d", i))
	}
}

func BenchmarkPubSub_Publish(b *testing.B) {
	ps := NewWithConfig(100)
	topic := "benchmark-topic"

	// Создаем подписчиков
	subscribers := make([]*Subscriber, 100)
	for i := 0; i < 100; i++ {
		subscribers[i] = ps.Subscribe(topic, fmt.Sprintf("subscriber-%d", i))
	}

	// Запускаем горутины для чтения сообщений
	for _, sub := range subscribers {
		go func(s *Subscriber) {
			for range s.Channel {
				// Просто читаем сообщения
			}
		}(sub)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Publish(topic, fmt.Sprintf("message-%d", i))
	}

	// Закрываем для завершения горутин
	ps.Close()
}

func BenchmarkPubSub_GetSubscribersCount(b *testing.B) {
	ps := NewWithConfig(100)
	topic := "benchmark-topic"

	// Создаем подписчиков
	for i := 0; i < 1000; i++ {
		ps.Subscribe(topic, fmt.Sprintf("subscriber-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.GetSubscribersCount(topic)
	}
}
