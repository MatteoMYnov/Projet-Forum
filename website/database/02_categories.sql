USE forum_y;

CREATE TABLE categories (
    id_category INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#1DA1F2',
    description TEXT,
    thread_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    -- Index pour optimiser les recherches
    INDEX idx_name (name),
    INDEX idx_active (is_active)
);
