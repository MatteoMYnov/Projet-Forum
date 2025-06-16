USE forum_y;

-- Modifier la table users pour définir une valeur par défaut pour le champ banner
ALTER TABLE users 
MODIFY COLUMN banner VARCHAR(255) DEFAULT '/img/banners/default-avatar.png';

-- Mettre à jour tous les utilisateurs existants qui n'ont pas de bannière
UPDATE users 
SET banner = '/img/banners/default-avatar.png' 
WHERE banner IS NULL OR banner = '';

-- Vérifier les modifications
SELECT id_user, username, banner FROM users LIMIT 10; 