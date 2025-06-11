USE forum_y;

-- Ajouter le champ love_count à la table threads
ALTER TABLE threads 
ADD COLUMN love_count INT DEFAULT 0 AFTER dislike_count;

-- Ajouter le champ love_count à la table messages aussi
ALTER TABLE messages 
ADD COLUMN love_count INT DEFAULT 0 AFTER dislike_count;

-- Mettre à jour les comptes de love existants pour les threads
UPDATE threads t 
SET love_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.thread_id = t.id_thread AND r.reaction_type = 'love'
);

-- Mettre à jour les comptes de love existants pour les messages
UPDATE messages m 
SET love_count = (
    SELECT COUNT(*) 
    FROM reactions r 
    WHERE r.message_id = m.id_message AND r.reaction_type = 'love'
);

SELECT 'Champ love_count ajouté et mis à jour avec succès' as Status; 