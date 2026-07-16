DROP INDEX idx_categories_slug ON categories;
DROP INDEX idx_categories_sort_order ON categories;
DROP INDEX idx_categories_status ON categories;

ALTER TABLE categories
DROP COLUMN slug,
DROP COLUMN description,
DROP COLUMN image_url,
DROP COLUMN sort_order,
DROP COLUMN status;
