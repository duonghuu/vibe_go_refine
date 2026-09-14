# Kế hoạch quản lý nội dung hình ảnh theo Type

  ## 1. Tóm tắt giải pháp

  - Xây domain mới ImageContent, tổ chức tương tự PostType + Post nhưng chỉ dành cho nội
    dung hình ảnh toàn site.

  - Không dùng trực tiếp posts vì các trường slug, category, content và SEO không phù hợp;
    không mở rộng page_section_items vì dữ liệu không thuộc từng Page và thiếu metadata
    theo ngữ cảnh.

  - Tái sử dụng hoàn toàn bảng và API media; không tạo luồng upload mới.
  - Phạm vi gồm Database, Backend API, Admin Refine và Public API. Chưa thay đổi cách
    render trong apps/webview.

  ## 2. Database và kiến trúc Backend

  Tạo hai bảng:

  - image_content_types:
      - id, code duy nhất và không đổi, name, status, sort_order, max_items.
      - field_config kiểu JSON, chỉ cấu hình bốn field được hỗ trợ: name, description,
        secondaryDescription, url; mỗi field có enabled và required.

      - Soft delete và timestamps.

  - image_contents:
      - id, type_code, media_id, name, description, secondary_description, target_url,
        sort_order, status, created_by.

      - FK tới image_content_types.code, media.id và users.id.
      - Index chính (type_code, status, sort_order, deleted_at), cùng index cho media_id,
        created_by.

      - Soft delete; cùng một media được phép tái sử dụng ở nhiều item.

  Seeder idempotent khởi tạo:

   Type       Quy tắc
  ━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   LOGO       Chỉ ảnh, maxItems = 1
  ─────────  ───────────────────────────────────────────────────
   SLIDER     Ảnh, tên, mô tả và URL bắt buộc; mô tả 2 tùy chọn
  ─────────  ───────────────────────────────────────────────────
   PARTNER    Ảnh, tên và URL bắt buộc

  Backend tuân theo DDD/CQRS hiện tại: Entity → Repository → Command/Query Service →
  Controller, đăng ký bằng Wire. Mutation dùng transaction và khóa row Type để bảo đảm
  maxItems và thứ tự không bị race condition.

  ## 3. API và Admin Refine

  API quản trị:

  - CRUD /api/v1/admin/image-content-types; chỉ ADMIN được thay đổi Type.
  - CRUD /api/v1/admin/image-contents; ADMIN và STAFF được quản lý nội dung.
  - PUT /api/v1/admin/image-contents/order nhận typeCode và toàn bộ snapshot {id,
    sortOrder} để reorder nguyên tử.

  - API danh sách trả {data, total}, hỗ trợ phân trang, tìm kiếm, lọc typeCode/status và
    sort tương thích Refine.

  - Cập nhật Type gây xung đột nếu maxItems mới nhỏ hơn số item hiện có hoặc cấu hình field
    bắt buộc không tương thích dữ liệu đang lưu.

  - URL chỉ chấp nhận đường dẫn nội bộ hoặc http/https; từ chối javascript:, data: và URL
    sai định dạng.

  Public API:

  GET /api/v1/public/image-contents?typeCodes=LOGO,SLIDER,PARTNER

  {
    "data": [
      {
        "typeCode": "SLIDER",
        "name": "Slider",
        "items": [
          {
            "id": 1,
            "name": "Tên slider",
            "description": "Mô tả",
            "secondaryDescription": null,
            "url": "/gioi-thieu",
            "sortOrder": 0,
            "media": {
              "id": 10,
              "originalUrl": "/uploads/example.webp",
              "thumbnailUrl": null,
              "mediumUrl": null,
              "mimeType": "image/webp"
            }
          }
        ]
      }
    ]
  }

  Public API chỉ trả Type và item ACTIVE, theo sortOrder ASC, id ASC; field bị vô hiệu hóa
  hoặc rỗng không được dùng làm dữ liệu hiển thị. Cập nhật OpenAPI và error contract gồm
  VALIDATION_ERROR, TYPE_NOT_FOUND, TYPE_CODE_CONFLICT, TYPE_CONFIG_CONFLICT,
  MAX_ITEMS_EXCEEDED, MEDIA_NOT_FOUND, MEDIA_TYPE_INVALID, FORBIDDEN, ORDER_CONFLICT.

  Admin bổ sung hai resource:

  - “Loại nội dung hình ảnh”: DataGrid và trang cấu hình Type; code khóa khi edit, cấu hình
    field bằng checkbox Enabled/Required và giới hạn item.

  - “Nội dung hình ảnh”: DataGrid theo Type, preview ảnh, trạng thái, thứ tự; form tự sinh
    field theo fieldConfig, dùng Media Library/upload hiện có.

  - Có đầy đủ Loading/Error/Empty, debounce tìm kiếm 300 ms, xác nhận xóa, cảnh báo unsaved
    changes và snackbar sau mutation.

  - Reorder lưu bằng snapshot; nút tạo bị vô hiệu hóa khi Type đạt maxItems.
  - Dùng type/interface TypeScript tường minh, không dùng any.

  ## 4. Quy tắc nghiệp vụ và cache

  - Media phải tồn tại, chưa bị xóa và có MIME image/*; STAFF chỉ được gắn media do mình sở
    hữu, ADMIN được dùng mọi media.

  - Khi gắn ảnh, chuyển media sang attached; khi xóa item không chuyển ngược về temporary
    vì media có thể đang được tái sử dụng.

  - required luôn kéo theo enabled; từ chối key lạ trong fieldConfig.
  - maxItems đếm mọi item chưa soft-delete, kể cả item INACTIVE; LOGO vì vậy luôn có tối đa
    một bản ghi.

  - Sau create/delete/reorder, chuẩn hóa thứ tự liên tục từ 0.
  - Cache toàn bộ public projection tại public:image-contents:v1, TTL 15 phút; mọi mutation
    xóa cache sau khi transaction commit.

  - Redis lỗi thì đọc trực tiếp DB và ghi log, không làm Public API mất khả dụng.
  - Tái sử dụng JWT middleware và Redis token hiện tại; không thay đổi refresh/logout.
  - Migration chỉ thêm hai bảng mới, không sửa dữ liệu media, posts, pages hoặc
    page_sections.

  ## 5. Kiểm thử và tiêu chí nghiệm thu

  - Kiểm thử migration up/down, FK/index và seeder chạy lặp không tạo bản ghi trùng.
  - Unit test Service cho field config, URL, quyền media, MIME, maxItems, cập nhật Type
    không tương thích và reorder stale/duplicate.

  - Integration test CRUD, RBAC, response Refine, public filtering/order, loại bỏ item
  - Kiểm tra transaction rollback khi một item/media trong request không hợp lệ.
  - Chạy go test ./..., build backend và npm run build cho Admin.
  - Nghiệm thu khi quản trị được Logo/Slider/Partner, Logo không vượt quá một item, thứ tự
    public chính xác, và có thể tạo một Type mới với tổ hợp field hiện hữu mà không cần đổi
    schema hoặc backend.

  - Không seed nội dung mẫu và không sửa apps/webview; việc render Logo/Hero/Partner từ
    Public API được tách sang task tích hợp webview sau.