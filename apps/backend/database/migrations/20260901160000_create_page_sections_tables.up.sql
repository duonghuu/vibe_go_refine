CREATE TABLE `page_sections` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `page_id` INT UNSIGNED NOT NULL,
  `key` VARCHAR(100) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `title` VARCHAR(255) NOT NULL DEFAULT '',
  `description` TEXT NOT NULL,
  `background_color` VARCHAR(9) NULL,
  `background_media_id` INT UNSIGNED NULL,
  `feature_media_id` INT UNSIGNED NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_page_sections_page_key` (`page_id`, `key`),
  KEY `idx_page_sections_page_status_sort` (`page_id`, `status`, `sort_order`, `deleted_at`),
  KEY `idx_page_sections_background_media` (`background_media_id`),
  KEY `idx_page_sections_feature_media` (`feature_media_id`),
  CONSTRAINT `chk_page_sections_status` CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_page_sections_sort_order` CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_page_sections_page` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_page_sections_background_media` FOREIGN KEY (`background_media_id`) REFERENCES `media` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_page_sections_feature_media` FOREIGN KEY (`feature_media_id`) REFERENCES `media` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `page_section_items` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `section_id` INT UNSIGNED NOT NULL,
  `item_type` VARCHAR(20) NOT NULL,
  `item_id` INT UNSIGNED NOT NULL,
  `collection` VARCHAR(100) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_page_section_item_link` (`section_id`, `collection`, `item_type`, `item_id`),
  KEY `idx_page_section_items_section_collection_sort` (`section_id`, `collection`, `sort_order`, `deleted_at`),
  KEY `idx_page_section_items_source` (`item_type`, `item_id`, `deleted_at`),
  CONSTRAINT `chk_page_section_items_type` CHECK (`item_type` IN ('CATEGORY', 'POST', 'MEDIA')),
  CONSTRAINT `chk_page_section_items_sort_order` CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_page_section_items_section` FOREIGN KEY (`section_id`) REFERENCES `page_sections` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
