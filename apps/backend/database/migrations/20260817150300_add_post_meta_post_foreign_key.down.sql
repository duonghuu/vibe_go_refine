ALTER TABLE `post_meta`
  DROP FOREIGN KEY `fk_post_meta_post`;

ALTER TABLE `post_meta`
  MODIFY COLUMN `post_id` BIGINT UNSIGNED NOT NULL;
