USE forum_y;

CREATE TABLE hashtags (
    id_hashtag INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    usage_count INT DEFAULT 0,
    trending_score DECIMAL(10,2) DEFAULT 0.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_used DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_trending BOOLEAN DEFAULT FALSE,
    -- Index pour optimiser les recherches
    INDEX idx_name (name),
    INDEX idx_usage_count (usage_count),
    INDEX idx_trending_score (trending_score),
    INDEX idx_last_used (last_used),
    INDEX idx_trending (is_trending)
);