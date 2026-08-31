CREATE TABLE `page_media` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `page_id` INT UNSIGNED NOT NULL,
  `media_id` INT UNSIGNED NOT NULL,
  `collection` VARCHAR(20) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  CONSTRAINT `chk_page_media_collection` CHECK (`collection` IN ('thumbnail', 'gallery')),
  CONSTRAINT `chk_page_media_sort_order` CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_page_media_page` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_page_media_media` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE KEY `uq_page_media_link` (`page_id`, `media_id`, `collection`),
  KEY `idx_page_media_page_collection_sort` (`page_id`, `collection`, `sort_order`, `deleted_at`),
  KEY `idx_page_media_media_id` (`media_id`),
  KEY `idx_page_media_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
