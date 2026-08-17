CREATE TABLE `posts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `type_code` varchar(50) NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `author_id` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_posts_slug` (`slug`),
  KEY `idx_posts_type_code` (`type_code`),
  KEY `idx_posts_author_id` (`author_id`),
  KEY `idx_posts_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_posts_type` FOREIGN KEY (`type_code`) REFERENCES `post_types` (`code`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_posts_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
