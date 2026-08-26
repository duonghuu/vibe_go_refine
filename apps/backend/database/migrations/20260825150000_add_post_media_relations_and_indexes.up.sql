-- Remove duplicate links before enforcing the composite uniqueness constraint.
DELETE pm_old
FROM post_media AS pm_old
INNER JOIN post_media AS pm_keep
  ON pm_old.post_id = pm_keep.post_id
 AND pm_old.media_id = pm_keep.media_id
 AND pm_old.collection = pm_keep.collection
 AND pm_old.id > pm_keep.id;

ALTER TABLE post_media
    ADD CONSTRAINT fk_post_media_post
        FOREIGN KEY (post_id) REFERENCES posts(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    ADD CONSTRAINT fk_post_media_media
        FOREIGN KEY (media_id) REFERENCES media(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    ADD UNIQUE KEY uq_post_media_link (post_id, media_id, collection),
    ADD KEY idx_post_media_post_collection_sort
        (post_id, collection, sort_order, deleted_at);
