USE forum_y;

-- 1. Supprimer d'abord les anciennes contraintes d'unicité
ALTER TABLE reactions 
DROP INDEX unique_user_thread_reaction,
DROP INDEX unique_user_message_reaction;

-- 2. Nettoyer les doublons existants en gardant seulement la réaction la plus récente
DELETE r1 FROM reactions r1
INNER JOIN reactions r2 
WHERE r1.user_id = r2.user_id 
  AND r1.thread_id = r2.thread_id 
  AND r1.id_reaction < r2.id_reaction;

DELETE r1 FROM reactions r1
INNER JOIN reactions r2 
WHERE r1.user_id = r2.user_id 
  AND r1.message_id = r2.message_id 
  AND r1.id_reaction < r2.id_reaction;

-- 3. Ajouter les nouvelles contraintes d'unicité
ALTER TABLE reactions 
ADD CONSTRAINT unique_user_thread UNIQUE (user_id, thread_id),
ADD CONSTRAINT unique_user_message UNIQUE (user_id, message_id);

-- 4. Mettre à jour les comptes dans les threads
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

-- 5. Mettre à jour les comptes dans les messages
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

SELECT 'Contraintes de réactions corrigées et doublons supprimés' as Status; 