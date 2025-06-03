-- Insertion de données de test
INSERT INTO utilisateurs (nom_utilisateur, email, mot_de_passe, bio) VALUES
('john_doe', 'john@example.com', SHA2('password123!', 512), 'Développeur passionné'),
('jane_smith', 'jane@example.com', SHA2('password456!', 512), 'Designer UI/UX'),
('tech_guru', 'tech@example.com', SHA2('password789!', 512), 'Expert en technologie');

-- Exemples de tweets
INSERT INTO tweets (contenu, id_utilisateur) VALUES
('Mon premier tweet! #hello #world', 1),
('Je découvre cette plateforme, c''est génial! #nouveau', 2),
('Les dernières tendances en développement web #webdev #coding', 3);

-- Exemples d'abonnements
INSERT INTO abonnements (id_suiveur, id_suivi) VALUES
(1, 2), -- John suit Jane
(1, 3), -- John suit Tech Guru
(2, 1), -- Jane suit John
(3, 1); -- Tech Guru suit John