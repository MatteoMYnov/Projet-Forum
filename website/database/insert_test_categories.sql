USE forum_y;

-- Insérer des catégories de test si elles n'existent pas déjà
INSERT IGNORE INTO categories (name, color, description) VALUES
('Discussion générale', '#1DA1F2', 'Pour toutes vos discussions générales'),
('Technologie', '#17BF63', 'Discussions sur la technologie, programmation, etc.'),
('Gaming', '#794BC4', 'Tout ce qui concerne les jeux vidéo'),
('Sport', '#E1306C', 'Discussions sportives'),
('Culture', '#F56040', 'Livres, films, musique et culture'),
('Aide & Support', '#FFAD1F', 'Demandes d''aide et support technique'); 