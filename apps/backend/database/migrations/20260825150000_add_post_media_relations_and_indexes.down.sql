ALTER TABLE post_media
    DROP FOREIGN KEY fk_post_media_post,
    DROP FOREIGN KEY fk_post_media_media,
    DROP INDEX uq_post_media_link,
    DROP INDEX idx_post_media_post_collection_sort;
