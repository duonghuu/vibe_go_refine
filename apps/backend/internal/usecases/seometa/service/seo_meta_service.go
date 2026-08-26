package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"go_refine_dashboard_be/internal/domain/seometa/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/seometa/dto"
)

var ErrValidation = errors.New("seo validation failed")

type SEOService interface {
	Get(ctx context.Context, entityType string, entityID uint) (*dto.Response, error)
	Upsert(ctx context.Context, entityType string, entityID uint, req *dto.UpsertRequest) (*dto.UpsertResponse, error)
	Delete(ctx context.Context, entityType string, entityID uint) error
}
type seoService struct {
	repo     repository.SEOMetaRepository
	registry repository.EntityRegistry
	redis    *redis.Client
}

func NewSEOService(repo repository.SEOMetaRepository, registry repository.EntityRegistry, redisClient *redis.Client) SEOService {
	return &seoService{repo: repo, registry: registry, redis: redisClient}
}
func cacheKey(t string, id uint) string { return fmt.Sprintf("seo:%s:%d", t, id) }

func (s *seoService) Get(ctx context.Context, entityType string, entityID uint) (*dto.Response, error) {
	t := repository.EntityType(strings.ToLower(strings.TrimSpace(entityType)))
	ref, err := s.registry.Resolve(ctx, t, entityID)
	if err != nil {
		return nil, err
	}
	key := cacheKey(string(t), entityID)
	if s.redis != nil {
		if raw, e := s.redis.Get(ctx, key).Bytes(); e == nil {
			var cached dto.Response
			if json.Unmarshal(raw, &cached) == nil {
				return &cached, nil
			}
		} else if !errors.Is(e, redis.Nil) {
			log.Printf("seo cache read failed: %v", e)
		}
	}
	item, err := s.repo.FindByEntity(ctx, string(t), entityID)
	if err != nil {
		return nil, err
	}
	res := resolve(item, ref)
	out := &dto.Response{Data: res}
	if s.redis != nil {
		if raw, e := json.Marshal(out); e == nil {
			if e = s.redis.Set(ctx, key, raw, time.Hour).Err(); e != nil {
				log.Printf("seo cache write failed: %v", e)
			}
		}
	}
	return out, nil
}

func (s *seoService) Upsert(ctx context.Context, entityType string, entityID uint, req *dto.UpsertRequest) (*dto.UpsertResponse, error) {
	t := repository.EntityType(strings.ToLower(strings.TrimSpace(entityType)))
	ref, err := s.registry.Resolve(ctx, t, entityID)
	if err != nil {
		return nil, err
	}
	clean, err := validate(req)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.Upsert(ctx, &entity.SEOMeta{EntityType: string(t), EntityID: entityID, MetaTitle: clean.MetaTitle, MetaDescription: clean.MetaDescription, MetaKeywords: clean.MetaKeywords, CanonicalURL: clean.CanonicalURL, OGTitle: clean.OGTitle, OGDescription: clean.OGDescription, OGImage: clean.OGImage, OGType: clean.OGType, TwitterTitle: clean.TwitterTitle, TwitterDescription: clean.TwitterDescription, TwitterImage: clean.TwitterImage, TwitterCard: clean.TwitterCard, Robots: clean.Robots, SchemaJSON: clean.SchemaJSON})
	if err != nil {
		return nil, err
	}
	if s.redis != nil {
		if err := s.redis.Del(ctx, cacheKey(string(t), entityID)).Err(); err != nil {
			log.Printf("seo cache invalidation failed: %v", err)
		}
	}
	return &dto.UpsertResponse{Data: resolve(item, ref), Message: "Cập nhật SEO thành công"}, nil
}

func (s *seoService) Delete(ctx context.Context, entityType string, entityID uint) error {
	t := repository.EntityType(strings.ToLower(strings.TrimSpace(entityType)))
	if _, err := s.registry.Resolve(ctx, t, entityID); err != nil {
		return err
	}
	if err := s.repo.DeleteByEntity(ctx, string(t), entityID); err != nil {
		return err
	}
	if s.redis != nil {
		if err := s.redis.Del(ctx, cacheKey(string(t), entityID)).Err(); err != nil {
			log.Printf("seo cache invalidation failed: %v", err)
		}
	}
	return nil
}

func ptrClean(v *string, max int) (*string, error) {
	if v == nil {
		return nil, nil
	}
	x := strings.TrimSpace(*v)
	if x == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(x) > max {
		return nil, ErrValidation
	}
	return &x, nil
}
func validURL(v *string, max int) (*string, error) {
	x, err := ptrClean(v, max)
	if err != nil || x == nil {
		return x, err
	}
	u, err := url.ParseRequestURI(*x)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, ErrValidation
	}
	return x, nil
}
func validate(r *dto.UpsertRequest) (*dto.UpsertRequest, error) {
	var err error
	for _, p := range []*struct {
		v **string
		n int
	}{{&r.MetaTitle, 255}, {&r.MetaDescription, 500}, {&r.MetaKeywords, 500}, {&r.OGTitle, 255}, {&r.OGDescription, 500}, {&r.OGType, 100}, {&r.TwitterTitle, 255}, {&r.TwitterDescription, 500}, {&r.TwitterCard, 100}, {&r.Robots, 100}} {
		*p.v, err = ptrClean(*p.v, p.n)
		if err != nil {
			return nil, ErrValidation
		}
	}
	if r.Robots == nil {
		x := "index,follow"
		r.Robots = &x
	}
	if *r.Robots != "index,follow" && *r.Robots != "noindex,follow" && *r.Robots != "noindex,nofollow" {
		return nil, ErrValidation
	}
	if r.CanonicalURL, err = validURL(r.CanonicalURL, 500); err != nil {
		return nil, ErrValidation
	}
	if r.OGImage, err = validURL(r.OGImage, 500); err != nil {
		return nil, ErrValidation
	}
	if r.TwitterImage, err = validURL(r.TwitterImage, 500); err != nil {
		return nil, ErrValidation
	}
	if len(r.SchemaJSON) > 65536 || (len(r.SchemaJSON) > 0 && string(r.SchemaJSON) != "null" && !json.Valid(r.SchemaJSON)) {
		return nil, ErrValidation
	}
	if string(r.SchemaJSON) == "null" {
		r.SchemaJSON = nil
	}
	return r, nil
}
func resolve(item *entity.SEOMeta, ref *repository.EntityReference) dto.SEOMetaResponse {
	r := dto.SEOMetaResponse{EntityType: string(ref.Type), EntityID: ref.ID, ResolvedTitle: ref.Title, ResolvedDescription: ref.Description, ResolvedCanonicalURL: "", Robots: "index,follow"}
	if ref.Content != "" {
		r.ResolvedDescription = ref.Content
	}
	if item != nil {
		r.MetaTitle, r.MetaDescription, r.MetaKeywords, r.CanonicalURL = item.MetaTitle, item.MetaDescription, item.MetaKeywords, item.CanonicalURL
		r.OGTitle, r.OGDescription, r.OGImage, r.OGType = item.OGTitle, item.OGDescription, item.OGImage, item.OGType
		r.TwitterTitle, r.TwitterDescription, r.TwitterImage, r.TwitterCard = item.TwitterTitle, item.TwitterDescription, item.TwitterImage, item.TwitterCard
		if item.Robots != nil {
			r.Robots = *item.Robots
		}
		r.SchemaJSON = json.RawMessage(item.SchemaJSON)
	}
	if r.MetaTitle != nil {
		r.ResolvedTitle = *r.MetaTitle
	}
	if r.MetaDescription != nil {
		r.ResolvedDescription = *r.MetaDescription
	}
	if r.CanonicalURL != nil {
		r.ResolvedCanonicalURL = *r.CanonicalURL
	}
	if r.ResolvedCanonicalURL == "" {
		r.ResolvedCanonicalURL = ref.Slug
	}
	if r.OGTitle == nil {
		x := r.ResolvedTitle
		r.OGTitle = &x
	}
	if r.OGDescription == nil {
		x := r.ResolvedDescription
		r.OGDescription = &x
	}
	if r.OGImage == nil && ref.ThumbnailURL != "" {
		x := ref.ThumbnailURL
		r.OGImage = &x
	}
	if r.TwitterTitle == nil {
		r.TwitterTitle = r.OGTitle
	}
	if r.TwitterDescription == nil {
		r.TwitterDescription = r.OGDescription
	}
	if r.TwitterImage == nil {
		r.TwitterImage = r.OGImage
	}
	return r
}
