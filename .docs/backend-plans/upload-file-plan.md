# KẾ HOẠCH HẠ TẦNG BACKEND - CHỨC NĂNG UPLOAD FILE

## 1. Thiết kế Dữ liệu (Database Schema Struct)

### 1.1 GORM Model: `Media` (Bảng `media`)
Thực thể `Media` dùng để quản lý trạng thái, thông tin các file upload lên hệ thống (hình ảnh sản phẩm, avatar,...).

```go
package domain

import (
	"time"
	"gorm.io/gorm"
)

type MediaStatus string

const (
	MediaStatusTemporary MediaStatus = "temporary"
	MediaStatusAttached  MediaStatus = "attached"
)

type Media struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"originalName"`
	FileName     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"fileName"` // Tên generate (UUID) để tránh trùng/path traversal
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mimeType"`
	Size         int64          `gorm:"not null" json:"size"`
	OriginalUrl  string         `gorm:"type:varchar(500);not null" json:"originalUrl"`
	ThumbnailUrl string         `gorm:"type:varchar(500)" json:"thumbnailUrl"`
	MediumUrl    string         `gorm:"type:varchar(500)" json:"mediumUrl"`
	Status       MediaStatus    `gorm:"type:enum('temporary', 'attached');default:'temporary';index:idx_status_created" json:"status"`
	
	OwnerID      uint           `gorm:"not null;index" json:"ownerId"` // ID của user upload
	ProductID    *uint          `gorm:"index" json:"productId"` // Nullable, null nếu file mồ côi (chưa attach vào sản phẩm)
	
	CreatedAt    time.Time      `gorm:"index:idx_status_created" json:"createdAt"` // Index composite với status phục vụ cron job dọn rác
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

* **Constraints & Indexing**: 
  - Tạo Composite Index `idx_status_created` trên `(status, created_at)` để tối ưu hóa truy vấn dọn rác (Background Job).
  - Sử dụng Soft Delete với `DeletedAt`.
  - Khóa ngoại `OwnerID` trỏ tới bảng `User`. Khóa ngoại `ProductID` trỏ tới bảng `Product` (cho phép NULL).

## 2. Giao kèo API (API Contract & Context Auth)

### 2.1 API Upload File (Single File)
- **Method & Route:** `POST /api/v1/media/upload`
- **Auth:** Yêu cầu đăng nhập. Middleware xác thực JWT, lấy `userId` gán vào context.
- **Request Binding:** Dạng `multipart/form-data`.
  - Không định nghĩa Go struct cho thân Request JSON, thay vào đó lấy file qua `c.FormFile("file")`.
  - Giới hạn kích thước payload tại Gin Middleware bằng cách cấu hình `MaxMultipartMemory`.
- **Response Payload:**
```go
// Success Response (200 OK)
type UploadMediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalUrl  string `json:"originalUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	MediumUrl    string `json:"mediumUrl"`
	Status       string `json:"status"`
}
```

### 2.2 API Cập nhật Product (Bổ sung gắn Media)
- Khác với việc tạo API attach rời rạc, thiết kế tốt nhất là tích hợp ngay vào API Tạo/Cập nhật của Product: `POST /api/v1/products` và `PUT /api/v1/products/:id`.
- Yêu cầu Client gửi mảng `mediaIds` đã upload thành công (status=temporary) vào body JSON. Backend sau khi lưu Product thành công sẽ Update lại tất cả các bản ghi bảng `media` có ID trùng khớp sang trạng thái `status = 'attached'` và gắn `product_id = mới_tạo`.

### 2.3 API Xóa File Tạm Thời
- **Method & Route:** `DELETE /api/v1/media/:id`
- **Auth:** Yêu cầu đăng nhập. Check phân quyền sở hữu file (`media.owner_id == userId`) hoặc Role Admin.
- **Logic:** Cho phép người dùng xóa nhanh file upload nhầm trên giao diện lúc chưa lưu sản phẩm, tiến hành Hard Delete bản ghi khỏi DB và xóa file vật lý tương ứng.
- **Response:** `200 OK` với thông báo thành công.

### 2.4 Kiểm soát Bảo Mật tại API (Security Core)
- Không lưu nguyên tên file gửi lên, BẮT BUỘC đổi tên file thành định dạng `UUID_timestamp.ext` ở Server.
- CẤM việc tin vào phần mở rộng file, phải kiểm tra Magic Bytes ở header của file upload (thông qua thư viện `github.com/h2non/filetype` trong Go).
- Thư mục lưu file không được phân quyền thực thi (No executable permission) tại Server Nginx.

## 3. Xử lý Cache & Quản lý Trạng thái Token (Redis Integration)

### 3.1 Rate Limiting (Chống Spam Endpoint Upload)
- API upload là nơi có nguy cơ bị tấn công DDOS cao gây tốn tài nguyên ổ cứng.
- Triển khai Rate Limiting trên API Upload sử dụng Redis Counter.
- **Key Pattern:** `rate_limit:upload:{ip}` hoặc `rate_limit:upload:{user_id}`.
- **Giới hạn:** Cấu hình mức hợp lý, ví dụ tối đa `20 requests / 1 phút`. Nếu vượt quá trả về HTTP Status `429 Too Many Requests`.

### 3.2 Tái tích hợp luồng Auth & Blacklist
- Các API thao tác trên Media (Upload, Delete) sẽ được bọc bởi JwtAuthMiddleware.
- **Logic Xác thực:** 
  - Token hợp lệ -> Kiểm tra Key Access Token tồn tại trong tập Blacklist của Redis hay chưa. Nếu có (do user logout hoặc phát hiện rủi ro) -> `401 Unauthorized`.
  - Việc thu hồi Token tuân thủ: Xác định TTL còn lại của Access Token lúc bị thu hồi và lưu Key vào Redis dùng lệnh `SETEX` để đảm bảo hệ thống Redis tự động dọn dẹp RAM đúng theo vòng đời thực của Token.

## 4. Xử lý File mồ côi (Cron Background Job)

- Tầng 2 phòng thủ theo yêu cầu của hội đồng hệ thống: Xử lý rác.
- Backend triển khai một Cron Job định kỳ (VD: mỗi tiếng 1 lần) để quét dọn các file có trạng thái tạm thời nhưng mãi chưa được đính kèm vào Product.
- Logic truy vấn Go Gorm:
```go
db.Where("status = ? AND created_at < ?", domain.MediaStatusTemporary, time.Now().Add(-24 * time.Hour)).Find(&garbageMedias)
```
- Duyệt vòng lặp:
  - Xóa file gốc + file thumbnail/medium tại Storage/ổ cứng.
  - Xóa bản ghi trong DB.
