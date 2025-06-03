USE forum_y;

CREATE TABLE messages (
    id_message INT AUTO_INCREMENT PRIMARY KEY,
    thread_id INT NOT NULL,
    author_id INT NOT NULL,
    content TEXT NOT NULL,
    parent_message_id INT DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    like_count INT DEFAULT 0,
    dislike_count INT DEFAULT 0,
    is_edited BOOLEAN DEFAULT FALSE,
    -- Clés étrangères
    FOREIGN KEY (thread_id) REFERENCES threads(id_thread) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id_user) ON DELETE CASCADE,
    FOREIGN KEY (parent_message_id) REFERENCES messages(id_message) ON DELETE CASCADE,
    -- Index pour optimiser les recherches
    INDEX idx_thread (thread_id),
    INDEX idx_author (author_id),
    INDEX idx_parent (parent_message_id),
    INDEX idx_created_at (created_at)
);