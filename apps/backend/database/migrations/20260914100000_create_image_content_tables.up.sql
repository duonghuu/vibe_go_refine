CREATE TABLE `image_content_types` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(50) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `sort_order` INT NOT NULL DEFAULT 0,
  `max_items` INT UNSIGNED NULL,
  `field_config` JSON NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_image_content_types_code` (`code`),
  KEY `idx_image_content_types_status_sort` (`status`, `sort_order`, `deleted_at`),
  KEY `idx_image_content_types_deleted_at` (`deleted_at`),
  CONSTRAINT `chk_image_content_types_status`
    CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_image_content_types_sort_order`
    CHECK (`sort_order` >= 0),
  CONSTRAINT `chk_image_content_types_max_items`
    CHECK (`max_items` IS NULL OR (`max_items` >= 1 AND `max_items` <= 1000))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `image_contents` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_code` VARCHAR(50) NOT NULL,
  `media_id` INT UNSIGNED NOT NULL,
  `name` VARCHAR(255) NULL,
  `description` TEXT NULL,
  `secondary_description` TEXT NULL,
  `target_url` VARCHAR(500) NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_by` INT NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_image_contents_type_status_sort`
    (`type_code`, `status`, `sort_order`, `deleted_at`),
  KEY `idx_image_contents_media_id` (`media_id`),
  KEY `idx_image_contents_created_by` (`created_by`),
  KEY `idx_image_contents_deleted_at` (`deleted_at`),
  CONSTRAINT `chk_image_contents_status`
    CHECK (`status` IN ('ACTIVE', 'INACTIVE')),
  CONSTRAINT `chk_image_contents_sort_order`
    CHECK (`sort_order` >= 0),
  CONSTRAINT `fk_image_contents_type`
    FOREIGN KEY (`type_code`) REFERENCES `image_content_types` (`code`)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_image_contents_media`
    FOREIGN KEY (`media_id`) REFERENCES `media` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_image_contents_created_by`
    FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
