USE forum_y;

-- Supprimer la table si elle existe pour permettre la modification
DROP TABLE IF EXISTS reactions;

CREATE TABLE reactions (
    id_reaction INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    thread_id INT DEFAULT NULL,
    message_id INT DEFAULT NULL,
    reaction_type ENUM('like', 'dislike', 'love', 'laugh', 'wow', 'sad', 'angry', 'repost') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- Clés étrangères
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES threads(id_thread) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES messages(id_message) ON DELETE CASCADE,
    -- Contraintes d'unicité pour éviter les doublons - un utilisateur ne peut avoir qu'une seule réaction par thread/message
    UNIQUE KEY unique_user_thread (user_id, thread_id),
    UNIQUE KEY unique_user_message (user_id, message_id),
    -- Index pour optimiser les recherches
    INDEX idx_user (user_id),
    INDEX idx_thread (thread_id),
    INDEX idx_message (message_id),
    INDEX idx_reaction_type (reaction_type),
    INDEX idx_created_at (created_at),
    -- Contrainte : une réaction doit être soit sur un thread, soit sur un message
    CHECK ((thread_id IS NOT NULL AND message_id IS NULL) OR (thread_id IS NULL AND message_id IS NOT NULL))
);

-- Insérer quelques données de test pour les réactions
INSERT INTO reactions (user_id, thread_id, reaction_type) VALUES
(1, 1, 'like'),
(2, 1, 'love'),
(1, 2, 'wow');