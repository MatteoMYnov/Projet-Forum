USE forum_y;

CREATE TABLE follows (
    id_follow INT AUTO_INCREMENT PRIMARY KEY,
    follower_id INT NOT NULL,
    following_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- Clés étrangères
    FOREIGN KEY (follower_id) REFERENCES users(id_user) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id_user) ON DELETE CASCADE,
    -- Contrainte d'unicité pour éviter les doublons
    UNIQUE KEY unique_follow (follower_id, following_id),
    -- Index pour optimiser les recherches
    INDEX idx_follower (follower_id),
    INDEX idx_following (following_id),
    INDEX idx_created_at (created_at),
    -- Contrainte : un utilisateur ne peut pas se suivre lui-même
    CHECK (follower_id != following_id)
);