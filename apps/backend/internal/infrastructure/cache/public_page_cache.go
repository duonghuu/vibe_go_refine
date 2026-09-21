package cache

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"

	"github.com/redis/go-redis/v9"
)

const publicPageCachePrefix = "public:pages:slug:"

type PublicPageCacheInvalidator interface {
	InvalidateSlugs(ctx context.Context, slugs ...string)
	InvalidatePage(ctx context.Context, pageID uint)
	InvalidateSource(ctx context.Context, itemType pageEntity.PageSectionItemType, itemIDs ...uint)
	InvalidateMedia(ctx context.Context, mediaIDs ...uint)
	FindMediaSlugs(ctx context.Context, mediaIDs ...uint) ([]string, error)
}

type publicPageCacheInvalidator struct {
	redis *redis.Client
	repo  repository.PublicPageCacheRepository
}

func NewPublicPageCacheInvalidator(redisClient *redis.Client, repo repository.PublicPageCacheRepository) PublicPageCacheInvalidator {
	return &publicPageCacheInvalidator{redis: redisClient, repo: repo}
}

func (i *publicPageCacheInvalidator) InvalidateSlugs(ctx context.Context, slugs ...string) {
	keys := make([]string, 0, len(slugs))
	seen := make(map[string]struct{}, len(slugs))
	for _, slug := range slugs {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if slug == "" {
			continue
		}
		key := PublicPageCacheKey(slug)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	if len(keys) == 0 || i.redis == nil {
		return
	}

	backoffs := []time.Duration{0, 50 * time.Millisecond, 150 * time.Millisecond}
	var lastErr error
	for attempt, backoff := range backoffs {
		if backoff > 0 {
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}

		pipe := i.redis.Pipeline()
		for _, key := range keys {
			pipe.Del(ctx, key)
		}
		if _, err := pipe.Exec(ctx); err == nil {
			return
		} else {
			lastErr = err
			log.Printf("public page cache invalidation attempt %d failed for %d keys: %v", attempt+1, len(keys), err)
		}
		if ctx.Err() != nil {
			return
		}
	}
	log.Printf("public page cache invalidation exhausted retries: %v", lastErr)
}

func (i *publicPageCacheInvalidator) InvalidatePage(ctx context.Context, pageID uint) {
	if i.repo == nil {
		return
	}
	slug, err := i.repo.FindPageSlug(ctx, pageID)
	if err != nil {
		log.Printf("public page cache page lookup failed for page %d: %v", pageID, err)
		return
	}
	i.InvalidateSlugs(ctx, slug)
}

func (i *publicPageCacheInvalidator) InvalidateSource(ctx context.Context, itemType pageEntity.PageSectionItemType, itemIDs ...uint) {
	if i.repo == nil || len(itemIDs) == 0 {
		return
	}
	slugs, err := i.repo.FindPageSlugsBySource(ctx, itemType, uniqueIDs(itemIDs))
	if err != nil {
		log.Printf("public page cache source lookup failed for type %q: %v", itemType, err)
		return
	}
	i.InvalidateSlugs(ctx, slugs...)
}

func (i *publicPageCacheInvalidator) InvalidateMedia(ctx context.Context, mediaIDs ...uint) {
	slugs, err := i.FindMediaSlugs(ctx, mediaIDs...)
	if err != nil {
		log.Printf("public page cache media lookup failed: %v", err)
		return
	}
	i.InvalidateSlugs(ctx, slugs...)
}

func (i *publicPageCacheInvalidator) FindMediaSlugs(ctx context.Context, mediaIDs ...uint) ([]string, error) {
	if i.repo == nil || len(mediaIDs) == 0 {
		return []string{}, nil
	}
	return i.repo.FindPageSlugsByMedia(ctx, uniqueIDs(mediaIDs))
}

func uniqueIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func PublicPageCacheKey(slug string) string {
	return fmt.Sprintf("%s%s", publicPageCachePrefix, strings.ToLower(strings.TrimSpace(slug)))
}
