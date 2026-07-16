ALTER TABLE categories
ADD COLUMN slug VARCHAR(255) NOT NULL AFTER name,
ADD COLUMN description TEXT NULL AFTER parent_id,
ADD COLUMN image_url VARCHAR(500) NULL AFTER description,
ADD COLUMN sort_order INT DEFAULT 0 AFTER image_url,
ADD COLUMN status VARCHAR(20) DEFAULT 'ACTIVE' AFTER sort_order;

CREATE UNIQUE INDEX idx_categories_slug ON categories (slug);
CREATE INDEX idx_categories_sort_order ON categories (sort_order);
CREATE INDEX idx_categories_status ON categories (status);
