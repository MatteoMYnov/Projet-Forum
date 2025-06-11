USE forum_y;

-- ===================================
-- ÉTAPE 1: Ajouter les colonnes love_count si elles n'existent pas
-- ===================================

-- Vérifier et ajouter love_count à threads
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_SCHEMA = 'forum_y' 
                   AND TABLE_NAME = 'threads' 
                   AND COLUMN_NAME = 'love_count');

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE threads ADD COLUMN love_count INT DEFAULT 0 AFTER dislike_count',
    'SELECT "Column threads.love_count already exists" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Vérifier et ajouter love_count à messages
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_SCHEMA = 'forum_y' 
                   AND TABLE_NAME = 'messages' 
                   AND COLUMN_NAME = 'love_count');

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE messages ADD COLUMN love_count INT DEFAULT 0 AFTER dislike_count',
    'SELECT "Column messages.love_count already exists" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ===================================
-- ÉTAPE 2: Corriger les contraintes d'unicité des réactions
-- ===================================

-- Supprimer les anciennes contraintes (si elles existent)
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
                     WHERE TABLE_SCHEMA = 'forum_y' 
                     AND TABLE_NAME = 'reactions' 
                     AND INDEX_NAME = 'unique_user_thread_reaction');

SET @sql = IF(@index_exists > 0,
    'ALTER TABLE reactions DROP INDEX unique_user_thread_reaction',
    'SELECT "Index unique_user_thread_reaction does not exist" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
                     WHERE TABLE_SCHEMA = 'forum_y' 
                     AND TABLE_NAME = 'reactions' 
                     AND INDEX_NAME = 'unique_user_message_reaction');

SET @sql = IF(@index_exists > 0,
    'ALTER TABLE reactions DROP INDEX unique_user_message_reaction',
    'SELECT "Index unique_user_message_reaction does not exist" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Nettoyer les doublons pour les threads (garder le plus récent)
DELETE r1 FROM reactions r1
INNER JOIN reactions r2 
WHERE r1.user_id = r2.user_id 
  AND r1.thread_id = r2.thread_id 
  AND r1.thread_id IS NOT NULL
  AND r1.id_reaction < r2.id_reaction;

-- Nettoyer les doublons pour les messages (garder le plus récent)
DELETE r1 FROM reactions r1
INNER JOIN reactions r2 
WHERE r1.user_id = r2.user_id 
  AND r1.message_id = r2.message_id 
  AND r1.message_id IS NOT NULL
  AND r1.id_reaction < r2.id_reaction;

-- Ajouter les nouvelles contraintes d'unicité (si elles n'existent pas)
SET @constraint_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
                          WHERE TABLE_SCHEMA = 'forum_y' 
                          AND TABLE_NAME = 'reactions' 
                          AND INDEX_NAME = 'unique_user_thread');

SET @sql = IF(@constraint_exists = 0,
    'ALTER TABLE reactions ADD CONSTRAINT unique_user_thread UNIQUE (user_id, thread_id)',
    'SELECT "Constraint unique_user_thread already exists" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @constraint_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
                          WHERE TABLE_SCHEMA = 'forum_y' 
                          AND TABLE_NAME = 'reactions' 
                          AND INDEX_NAME = 'unique_user_message');

SET @sql = IF(@constraint_exists = 0,
    'ALTER TABLE reactions ADD CONSTRAINT unique_user_message UNIQUE (user_id, message_id)',
    'SELECT "Constraint unique_user_message already exists" as Info');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ===================================
-- ÉTAPE 3: Recalculer tous les compteurs
-- ===================================

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
),
love_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.thread_id = t.id_thread AND r.reaction_type = 'love'
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
),
love_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.message_id = m.id_message AND r.reaction_type = 'love'
);

-- ===================================
-- ÉTAPE 4: Rapport final
-- ===================================

SELECT 
    'Migration des réactions terminée avec succès!' as Status,
    (SELECT COUNT(*) FROM reactions) as Total_Reactions,
    (SELECT COUNT(*) FROM reactions WHERE reaction_type = 'like') as Like_Count,
    (SELECT COUNT(*) FROM reactions WHERE reaction_type = 'dislike') as Dislike_Count,
    (SELECT COUNT(*) FROM reactions WHERE reaction_type = 'love') as Love_Count,
    (SELECT COUNT(*) FROM threads WHERE love_count > 0) as Threads_With_Love,
    (SELECT COUNT(*) FROM messages WHERE love_count > 0) as Messages_With_Love;

SELECT 'Vérification des contraintes d\'unicité...' as Info;

-- Vérifier qu'il n'y a plus de doublons
SELECT 
    CASE 
        WHEN (SELECT COUNT(*) FROM (
            SELECT user_id, thread_id, COUNT(*) as cnt 
            FROM reactions 
            WHERE thread_id IS NOT NULL 
            GROUP BY user_id, thread_id 
            HAVING cnt > 1
        ) as duplicates) = 0 
        THEN '✅ Aucun doublon thread détecté'
        ELSE '❌ Des doublons thread subsistent'
    END as Thread_Duplicates_Check;

SELECT 
    CASE 
        WHEN (SELECT COUNT(*) FROM (
            SELECT user_id, message_id, COUNT(*) as cnt 
            FROM reactions 
            WHERE message_id IS NOT NULL 
            GROUP BY user_id, message_id 
            HAVING cnt > 1
        ) as duplicates) = 0 
        THEN '✅ Aucun doublon message détecté'
        ELSE '❌ Des doublons message subsistent'
    END as Message_Duplicates_Check; 