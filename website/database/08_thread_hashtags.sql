USE forum_y;

CREATE TABLE thread_hashtags (
    thread_id INT NOT NULL,
    hashtag_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- Clé primaire composée
    PRIMARY KEY (thread_id, hashtag_id),
    -- Clés étrangères
    FOREIGN KEY (thread_id) REFERENCES threads(id_thread) ON DELETE CASCADE,
    FOREIGN KEY (hashtag_id) REFERENCES hashtags(id_hashtag) ON DELETE CASCADE,
    -- Index pour optimiser les recherches
    INDEX idx_thread (thread_id),
    INDEX idx_hashtag (hashtag_id),
    INDEX idx_created_at (created_at)
);