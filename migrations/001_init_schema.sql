-- migrations/001_init_schema.sql
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    comments_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE, -- NULL для корневых комментариев
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Базовые индексы для поиска комментариев
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);

-- ПРОИЗВОДИТЕЛЬНОСТЬ: Составные индексы для пагинации
-- Индекс для быстрого поиска корневых комментариев с сортировкой
CREATE INDEX idx_comments_post_root_created ON comments(post_id, created_at)
WHERE parent_id IS NULL;

-- Индекс для быстрого поиска дочерних комментариев с сортировкой
CREATE INDEX idx_comments_parent_created ON comments(parent_id, created_at)
WHERE parent_id IS NOT NULL;

-- Индекс для поиска постов по времени создания
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- Индекс для активных постов (с включенными комментариями)
CREATE INDEX idx_posts_comments_enabled ON posts(comments_enabled)
WHERE comments_enabled = true;
