-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Hôte : 127.0.0.1
-- Généré le : mar. 17 juin 2025 à 09:45
-- Version du serveur : 10.4.32-MariaDB
-- Version de PHP : 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de données : `forum_y`
--

-- --------------------------------------------------------

--
-- Structure de la table `categories`
--

CREATE TABLE `categories` (
  `id_category` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `color` varchar(7) DEFAULT '#1DA1F2',
  `description` text DEFAULT NULL,
  `thread_count` int(11) DEFAULT 0,
  `created_at` datetime DEFAULT current_timestamp(),
  `is_active` tinyint(1) DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `categories`
--

INSERT INTO `categories` (`id_category`, `name`, `color`, `description`, `thread_count`, `created_at`, `is_active`) VALUES
(1, 'Tech', '#1DA1F2', 'Discussions technologie', 0, '2025-06-03 09:43:11', 1),
(2, 'Gaming', '#FF6B6B', 'Jeux vidéo', 0, '2025-06-03 09:43:11', 1);

-- --------------------------------------------------------

--
-- Structure de la table `follows`
--

CREATE TABLE `follows` (
  `id_follow` int(11) NOT NULL,
  `follower_id` int(11) NOT NULL,
  `following_id` int(11) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp()
) ;

-- --------------------------------------------------------

--
-- Structure de la table `hashtags`
--

CREATE TABLE `hashtags` (
  `id_hashtag` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `usage_count` int(11) DEFAULT 0,
  `trending_score` decimal(10,2) DEFAULT 0.00,
  `created_at` datetime DEFAULT current_timestamp(),
  `last_used` datetime DEFAULT current_timestamp(),
  `is_trending` tinyint(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `messages`
--

CREATE TABLE `messages` (
  `id_message` int(11) NOT NULL,
  `thread_id` int(11) NOT NULL,
  `author_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `parent_message_id` int(11) DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `like_count` int(11) DEFAULT 0,
  `dislike_count` int(11) DEFAULT 0,
  `love_count` int(11) DEFAULT 0,
  `is_edited` tinyint(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `messages`
--

INSERT INTO `messages` (`id_message`, `thread_id`, `author_id`, `content`, `parent_message_id`, `created_at`, `updated_at`, `like_count`, `dislike_count`, `love_count`, `is_edited`) VALUES
(2, 8, 3, 'caca', NULL, '2025-06-11 14:11:01', '2025-06-11 14:11:25', 0, 0, 0, 0),
(3, 6, 3, 'caca', NULL, '2025-06-11 14:11:37', '2025-06-11 14:11:37', 0, 0, 0, 0),
(4, 8, 3, 'kaka', NULL, '2025-06-11 15:09:43', '2025-06-11 15:09:43', 0, 0, 0, 0),
(5, 6, 3, 'kkkk', NULL, '2025-06-11 17:07:18', '2025-06-11 17:07:18', 0, 0, 0, 0),
(6, 6, 3, '快去看看去、', NULL, '2025-06-11 17:07:22', '2025-06-11 17:07:22', 0, 0, 0, 0),
(7, 3, 3, 'csgo', NULL, '2025-06-11 17:07:34', '2025-06-11 17:07:34', 0, 0, 0, 0),
(9, 9, 6, 'bbbbbbbbbbb', NULL, '2025-06-11 20:55:03', '2025-06-11 20:55:03', 0, 0, 0, 0),
(10, 9, 3, 'bbbbbb', NULL, '2025-06-11 23:02:59', '2025-06-11 23:02:59', 0, 0, 0, 0),
(11, 10, 11, 'tghgrfc', NULL, '2025-06-17 09:23:38', '2025-06-17 09:23:38', 0, 0, 0, 0),
(12, 10, 11, 'fdfbdfrfgr', NULL, '2025-06-17 09:23:40', '2025-06-17 09:23:40', 0, 0, 0, 0);

-- --------------------------------------------------------

--
-- Structure de la table `reactions`
--

CREATE TABLE `reactions` (
  `id_reaction` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `thread_id` int(11) DEFAULT NULL,
  `message_id` int(11) DEFAULT NULL,
  `reaction_type` enum('like','dislike','love','repost') NOT NULL,
  `created_at` datetime DEFAULT current_timestamp()
) ;

--
-- Déchargement des données de la table `reactions`
--

INSERT INTO `reactions` (`id_reaction`, `user_id`, `thread_id`, `message_id`, `reaction_type`, `created_at`) VALUES
(55, 3, 7, NULL, 'dislike', '2025-06-11 12:09:23'),
(61, 3, 6, NULL, 'like', '2025-06-11 14:28:24'),
(63, 3, 8, NULL, 'like', '2025-06-11 14:39:41'),
(64, 3, 9, NULL, 'like', '2025-06-12 09:10:47'),
(65, 3, 3, NULL, 'like', '2025-06-12 09:10:49'),
(68, 3, 10, NULL, 'like', '2025-06-12 10:08:52'),
(70, 6, 10, NULL, 'like', '2025-06-13 02:28:49'),
(71, 11, 13, NULL, 'like', '2025-06-17 09:23:25');

-- --------------------------------------------------------

--
-- Structure de la table `threads`
--

CREATE TABLE `threads` (
  `id_thread` int(11) NOT NULL,
  `title` varchar(280) NOT NULL,
  `content` text NOT NULL,
  `author_id` int(11) NOT NULL,
  `category_id` int(11) DEFAULT NULL,
  `status` enum('open','closed','archived') DEFAULT 'open',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `is_pinned` tinyint(1) DEFAULT 0,
  `view_count` int(11) DEFAULT 0,
  `like_count` int(11) DEFAULT 0,
  `dislike_count` int(11) DEFAULT 0,
  `love_count` int(11) DEFAULT 0,
  `message_count` int(11) DEFAULT 0,
  `last_activity` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `threads`
--

INSERT INTO `threads` (`id_thread`, `title`, `content`, `author_id`, `category_id`, `status`, `created_at`, `updated_at`, `is_pinned`, `view_count`, `like_count`, `dislike_count`, `love_count`, `message_count`, `last_activity`) VALUES
(1, '🎉 Bienvenue sur Forum Y !', 'Premier post de test !', 1, 1, 'open', '2025-06-03 09:43:11', '2025-06-10 10:31:35', 0, 1, 0, 0, 0, 0, '2025-06-03 09:43:11'),
(2, 'Test post', 'Ceci est un test', 2, 2, 'open', '2025-06-03 09:43:11', '2025-06-11 20:35:28', 0, 3, 0, 0, 0, 0, '2025-06-03 09:43:11'),
(3, 'Gaming', 'Rien', 3, 2, 'open', '2025-06-10 10:31:17', '2025-06-12 09:54:55', 0, 15, 1, 0, 0, 1, '2025-06-12 09:54:54'),
(4, 'dzdzdzddzd', 'dz', 3, 1, 'open', '2025-06-10 10:44:25', '2025-06-10 20:35:07', 0, 3, 0, 0, 0, 0, '2025-06-10 10:44:25'),
(5, 'caca', 'rien', 3, 1, 'open', '2025-06-10 11:01:50', '2025-06-10 20:36:42', 0, 2, 0, 0, 0, 0, '2025-06-10 11:01:50'),
(6, 'caca', 'rien', 3, 2, 'open', '2025-06-10 11:08:51', '2025-06-15 16:44:46', 0, 11, 1, 0, 0, 3, '2025-06-11 17:07:22'),
(7, 'matteo', 'rien', 3, 2, 'archived', '2025-06-10 11:37:15', '2025-06-17 09:32:03', 0, 18, 0, 1, 0, 0, '2025-06-10 11:37:15'),
(8, 'jeu', 'jeu', 3, 2, 'open', '2025-06-10 11:50:12', '2025-06-17 09:16:35', 0, 90, 1, 0, 0, 2, '2025-06-17 09:16:32'),
(9, 'j\'aime Y forum', 'bebou', 6, 1, 'open', '2025-06-11 20:54:34', '2025-06-13 02:29:12', 0, 49, 1, 0, 0, 2, '2025-06-11 23:02:59'),
(10, 'ez ratio', 'blblblbl', 3, 2, 'open', '2025-06-12 10:08:46', '2025-06-17 09:31:56', 0, 19, 2, 0, 0, 2, '2025-06-17 09:23:40'),
(11, 'Ynov', 'Bachelor 1 Informatique', 3, NULL, 'open', '2025-06-17 08:57:18', '2025-06-17 09:25:07', 0, 9, 0, 0, 0, 0, '2025-06-17 08:57:18'),
(12, 'test', 'test', 3, NULL, 'open', '2025-06-17 09:12:27', '2025-06-17 09:24:47', 0, 6, 0, 0, 0, 0, '2025-06-17 09:12:27'),
(13, 'J\'adore les pizza', 'Pizza', 11, NULL, 'open', '2025-06-17 09:23:20', '2025-06-17 09:30:46', 0, 2, 1, 0, 0, 0, '2025-06-17 09:23:20'),
(14, 'xaxa', 'x', 3, NULL, 'closed', '2025-06-17 09:31:20', '2025-06-17 09:31:53', 0, 3, 0, 0, 0, 0, '2025-06-17 09:31:20');

-- --------------------------------------------------------

--
-- Structure de la table `thread_hashtags`
--

CREATE TABLE `thread_hashtags` (
  `thread_id` int(11) NOT NULL,
  `hashtag_id` int(11) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `users`
--

CREATE TABLE `users` (
  `id_user` int(11) NOT NULL,
  `username` varchar(50) NOT NULL,
  `display_name` varchar(100) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `profile_picture` varchar(500) DEFAULT '/img/default-avatar.png',
  `banner` varchar(255) DEFAULT '/img/banners/default-avatar.png',
  `banner_image` varchar(500) DEFAULT '/img/default-banner.png',
  `bio` text DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `last_login` datetime DEFAULT NULL,
  `is_verified` tinyint(1) DEFAULT 0,
  `is_banned` tinyint(1) DEFAULT 0,
  `role` enum('user','admin') DEFAULT 'user',
  `follower_count` int(11) DEFAULT 0,
  `following_count` int(11) DEFAULT 0,
  `thread_count` int(11) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `users`
--

INSERT INTO `users` (`id_user`, `username`, `display_name`, `email`, `password_hash`, `profile_picture`, `banner`, `banner_image`, `bio`, `created_at`, `last_login`, `is_verified`, `is_banned`, `role`, `follower_count`, `following_count`, `thread_count`) VALUES
(1, 'admin', 'admin', 'admin@forum-y.com', '7fcf4ba391c48784edde599889d6e3f1e47a27db36ecc050cc92f259bfac38afad2c68a1ae804d77075e8fb722503f3eca2b2c1006ee6f6c7b7628cb45fffd1d', '/img/avatars/default-avatar.png', '/img/banners/default-avatar.png', '/img/default-banner.png', NULL, '2025-06-03 09:43:11', NULL, 0, 0, 'admin', 0, 0, 0),
(2, 'alice', 'alice', 'alice@test.com', 'ff3db360a753794549736077fa2ac5da87c44de7265e03a1df4457045b1b54449cae2b1e22cf4915fc014f5a2a455f404444b3fff0361fa906f39d026016eb16', '/img/avatars/default-avatar.png', '/img/banners/default-avatar.png', '/img/default-banner.png', NULL, '2025-06-03 09:43:11', NULL, 0, 0, 'user', 0, 0, 0),
(3, 'caca', NULL, 'caca@ynov.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/1750108265_4b33f6881ed9b5b35fb19347690558e3.png', '/img/banners/1750108304_4981274cf51cde1de07570da7125285a.png', '/img/default-banner.png', NULL, '2025-06-06 10:01:11', '2025-06-17 09:24:37', 0, 0, 'user', 0, 0, 0),
(4, 'caca1', NULL, 'caca1@ynov.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/default-avatar.png', '/img/banners/default-avatar.png', '/img/default-banner.png', NULL, '2025-06-10 20:55:11', '2025-06-10 20:55:25', 0, 0, 'user', 0, 0, 0),
(5, 'caaca2', NULL, 'caca2@ynov.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/default-avatar.png', '/img/banners/default-avatar.png', '/img/default-banner.png', NULL, '2025-06-10 20:56:16', NULL, 0, 0, 'user', 0, 0, 0),
(6, 'romain', NULL, 'romain@ynov.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/1749581819_7363194bd43016225c90b1d43a85778b.png', '/img/banners/default-avatar.png', '/img/default-banner.png', NULL, '2025-06-10 20:56:59', '2025-06-17 09:40:56', 0, 0, 'user', 0, 0, 0),
(11, 'test', NULL, 'test@test.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/1750034653_381fbb8ca2428b9b0375be3bba810df5.png', '/img/banners/1750034653_3e6fc8693f82508502a048f81b74d81d.png', '/img/default-banner.png', NULL, '2025-06-16 02:44:13', '2025-06-17 09:22:34', 0, 0, 'user', 0, 0, 0),
(12, 'test1', NULL, 'test@gmai.com', 'c4d0e87a1d599766eeb2e8011786637b51194533017d289fbe2056a84fbc9bb80b4edeed72483f66da7e92bfbd64e667b90401f2b12839038c59b115ff17ee06', '/img/avatars/default-avatar.png', NULL, '/img/default-banner.png', NULL, '2025-06-16 03:14:18', '2025-06-16 03:14:37', 0, 0, 'user', 0, 0, 0);

-- --------------------------------------------------------

--
-- Structure de la table `wall_posts`
--

CREATE TABLE `wall_posts` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `author_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `wall_posts`
--

INSERT INTO `wall_posts` (`id`, `user_id`, `author_id`, `content`, `created_at`, `updated_at`) VALUES
(1, 6, 6, 'kaka', '2025-06-12 10:05:26', '2025-06-12 10:05:26'),
(2, 6, 6, 'kakakaka', '2025-06-12 10:05:30', '2025-06-12 10:05:30'),
(3, 6, 6, 'kaaka', '2025-06-12 10:08:05', '2025-06-12 10:08:05'),
(4, 6, 6, 'kaka', '2025-06-12 10:08:10', '2025-06-12 10:08:10'),
(5, 6, 6, 'kakakakaka', '2025-06-12 10:12:27', '2025-06-12 10:12:27'),
(6, 3, 3, 'kaka', '2025-06-15 14:31:10', '2025-06-15 14:31:10'),
(7, 3, 3, 'caca', '2025-06-16 11:33:10', '2025-06-16 11:33:10'),
(8, 3, 3, 'kakak', '2025-06-16 18:44:19', '2025-06-16 18:44:19'),
(9, 3, 3, 'bonjour je m\'appelle Jérémie', '2025-06-17 06:55:59', '2025-06-17 06:55:59'),
(10, 11, 11, 'j\'aime y formu', '2025-06-17 07:22:45', '2025-06-17 07:22:45');

--
-- Index pour les tables déchargées
--

--
-- Index pour la table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id_category`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `idx_name` (`name`),
  ADD KEY `idx_active` (`is_active`);

--
-- Index pour la table `follows`
--
ALTER TABLE `follows`
  ADD PRIMARY KEY (`id_follow`),
  ADD UNIQUE KEY `unique_follow` (`follower_id`,`following_id`),
  ADD KEY `idx_follower` (`follower_id`),
  ADD KEY `idx_following` (`following_id`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Index pour la table `hashtags`
--
ALTER TABLE `hashtags`
  ADD PRIMARY KEY (`id_hashtag`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `idx_name` (`name`),
  ADD KEY `idx_usage_count` (`usage_count`),
  ADD KEY `idx_trending_score` (`trending_score`),
  ADD KEY `idx_last_used` (`last_used`),
  ADD KEY `idx_trending` (`is_trending`);

--
-- Index pour la table `messages`
--
ALTER TABLE `messages`
  ADD PRIMARY KEY (`id_message`),
  ADD KEY `idx_thread` (`thread_id`),
  ADD KEY `idx_author` (`author_id`),
  ADD KEY `idx_parent` (`parent_message_id`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Index pour la table `reactions`
--
ALTER TABLE `reactions`
  ADD PRIMARY KEY (`id_reaction`),
  ADD UNIQUE KEY `unique_user_thread` (`user_id`,`thread_id`),
  ADD UNIQUE KEY `unique_user_message` (`user_id`,`message_id`),
  ADD KEY `idx_user` (`user_id`),
  ADD KEY `idx_thread` (`thread_id`),
  ADD KEY `idx_message` (`message_id`),
  ADD KEY `idx_reaction_type` (`reaction_type`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Index pour la table `threads`
--
ALTER TABLE `threads`
  ADD PRIMARY KEY (`id_thread`),
  ADD KEY `idx_author` (`author_id`),
  ADD KEY `idx_category` (`category_id`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_last_activity` (`last_activity`),
  ADD KEY `idx_status` (`status`),
  ADD KEY `idx_pinned` (`is_pinned`);

--
-- Index pour la table `thread_hashtags`
--
ALTER TABLE `thread_hashtags`
  ADD PRIMARY KEY (`thread_id`,`hashtag_id`),
  ADD KEY `idx_thread` (`thread_id`),
  ADD KEY `idx_hashtag` (`hashtag_id`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- Index pour la table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id_user`),
  ADD UNIQUE KEY `username` (`username`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_username` (`username`),
  ADD KEY `idx_email` (`email`),
  ADD KEY `idx_created_at` (`created_at`),
  ADD KEY `idx_role` (`role`),
  ADD KEY `idx_display_name` (`display_name`);

--
-- Index pour la table `wall_posts`
--
ALTER TABLE `wall_posts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_id` (`user_id`),
  ADD KEY `idx_created_at` (`created_at`);

--
-- AUTO_INCREMENT pour les tables déchargées
--

--
-- AUTO_INCREMENT pour la table `categories`
--
ALTER TABLE `categories`
  MODIFY `id_category` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT pour la table `follows`
--
ALTER TABLE `follows`
  MODIFY `id_follow` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT pour la table `hashtags`
--
ALTER TABLE `hashtags`
  MODIFY `id_hashtag` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT pour la table `messages`
--
ALTER TABLE `messages`
  MODIFY `id_message` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT pour la table `reactions`
--
ALTER TABLE `reactions`
  MODIFY `id_reaction` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT pour la table `threads`
--
ALTER TABLE `threads`
  MODIFY `id_thread` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT pour la table `users`
--
ALTER TABLE `users`
  MODIFY `id_user` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT pour la table `wall_posts`
--
ALTER TABLE `wall_posts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Contraintes pour les tables déchargées
--

--
-- Contraintes pour la table `follows`
--
ALTER TABLE `follows`
  ADD CONSTRAINT `follows_ibfk_1` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE,
  ADD CONSTRAINT `follows_ibfk_2` FOREIGN KEY (`following_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE;

--
-- Contraintes pour la table `messages`
--
ALTER TABLE `messages`
  ADD CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`thread_id`) REFERENCES `threads` (`id_thread`) ON DELETE CASCADE,
  ADD CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`author_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE,
  ADD CONSTRAINT `messages_ibfk_3` FOREIGN KEY (`parent_message_id`) REFERENCES `messages` (`id_message`) ON DELETE CASCADE;

--
-- Contraintes pour la table `reactions`
--
ALTER TABLE `reactions`
  ADD CONSTRAINT `reactions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE,
  ADD CONSTRAINT `reactions_ibfk_2` FOREIGN KEY (`thread_id`) REFERENCES `threads` (`id_thread`) ON DELETE CASCADE,
  ADD CONSTRAINT `reactions_ibfk_3` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id_message`) ON DELETE CASCADE;

--
-- Contraintes pour la table `threads`
--
ALTER TABLE `threads`
  ADD CONSTRAINT `threads_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE,
  ADD CONSTRAINT `threads_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id_category`) ON DELETE SET NULL;

--
-- Contraintes pour la table `thread_hashtags`
--
ALTER TABLE `thread_hashtags`
  ADD CONSTRAINT `thread_hashtags_ibfk_1` FOREIGN KEY (`thread_id`) REFERENCES `threads` (`id_thread`) ON DELETE CASCADE,
  ADD CONSTRAINT `thread_hashtags_ibfk_2` FOREIGN KEY (`hashtag_id`) REFERENCES `hashtags` (`id_hashtag`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
