ALTER TABLE `post_meta`
  MODIFY COLUMN `post_id` INT UNSIGNED NOT NULL;

ALTER TABLE `post_meta`
  ADD CONSTRAINT `fk_post_meta_post`
  FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
  ON DELETE CASCADE ON UPDATE CASCADE;
