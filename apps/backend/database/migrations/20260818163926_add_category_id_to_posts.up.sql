ALTER TABLE posts 
ADD COLUMN category_id INT UNSIGNED NULL;

CREATE INDEX idx_posts_category_id ON posts(category_id);

ALTER TABLE posts 
ADD CONSTRAINT fk_posts_category 
FOREIGN KEY (category_id) 
REFERENCES post_categories(id) 
ON DELETE SET NULL;
