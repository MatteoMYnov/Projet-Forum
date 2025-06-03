USE forum_y;

CREATE TABLE threads (
    id_thread INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(280) NOT NULL,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    category_id INT,
    status ENUM('open', 'closed', 'archived') DEFAULT 'open',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_pinned BOOLEAN DEFAULT FALSE,
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    dislike_count INT DEFAULT 0,
    message_count INT DEFAULT 0,
    last_activity DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- Clés étrangères
    FOREIGN KEY (author_id) REFERENCES users(id_user) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id_category) ON DELETE SET NULL,
    -- Index pour optimiser les recherches
    INDEX idx_author (author_id),
    INDEX idx_category (category_id),
    INDEX idx_created_at (created_at),
    INDEX idx_last_activity (last_activity),
    INDEX idx_status (status),
    INDEX idx_pinned (is_pinned)
);