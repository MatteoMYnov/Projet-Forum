-- Script de mise à jour pour ajouter les nouvelles réactions
USE forum_y;

-- Sauvegarder les réactions existantes
CREATE TEMPORARY TABLE temp_reactions AS
SELECT user_id, thread_id, message_id, reaction_type, created_at
FROM reactions
WHERE reaction_type IN ('like', 'dislike', 'love', 'repost');

-- Supprimer l'ancienne table
DROP TABLE reactions;

-- Créer la nouvelle table avec tous les types de réactions
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
    -- Contraintes d'unicité pour éviter les doublons
    UNIQUE KEY unique_user_thread_reaction (user_id, thread_id, reaction_type),
    UNIQUE KEY unique_user_message_reaction (user_id, message_id, reaction_type),
    -- Index pour optimiser les recherches
    INDEX idx_user (user_id),
    INDEX idx_thread (thread_id),
    INDEX idx_message (message_id),
    INDEX idx_reaction_type (reaction_type),
    INDEX idx_created_at (created_at),
    -- Contrainte : une réaction doit être soit sur un thread, soit sur un message
    CHECK ((thread_id IS NOT NULL AND message_id IS NULL) OR (thread_id IS NULL AND message_id IS NOT NULL))
);

-- Restaurer les données sauvegardées
INSERT INTO reactions (user_id, thread_id, message_id, reaction_type, created_at)
SELECT user_id, thread_id, message_id, reaction_type, created_at
FROM temp_reactions;

-- Supprimer la table temporaire
DROP TEMPORARY TABLE temp_reactions;

-- Ajouter quelques réactions de test si aucune donnée n'existait
INSERT IGNORE INTO reactions (user_id, thread_id, reaction_type) VALUES
(1, 1, 'like'),
(2, 1, 'love'),
(1, 2, 'wow'),
(3, 1, 'laugh');

-- Mettre à jour les comptes dans la table threads
UPDATE threads t 
SET like_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.thread_id = t.id_thread AND r.reaction_type = 'like'
),
dislike_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.thread_id = t.id_thread AND r.reaction_type = 'dislike'
);

-- Mettre à jour les comptes dans la table messages
UPDATE messages m 
SET like_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.message_id = m.id_message AND r.reaction_type = 'like'
),
dislike_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.message_id = m.id_message AND r.reaction_type = 'dislike'
);

SELECT 'Mise à jour des réactions terminée!' as status; 