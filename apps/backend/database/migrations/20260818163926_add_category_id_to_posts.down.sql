ALTER TABLE posts DROP FOREIGN KEY fk_posts_category;
DROP INDEX idx_posts_category_id ON posts;
ALTER TABLE posts DROP COLUMN category_id;
