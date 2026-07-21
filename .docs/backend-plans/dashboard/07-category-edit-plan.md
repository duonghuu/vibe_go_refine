# QUY HOẠCH KIẾN TRÚC BACK-END: CHỈNH SỬA DANH MỤC SẢN PHẨM

**Dự án:** TechBite
**Module:** Category Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/07-category-edit-idea.md`

*(Ghi chú: Bản thiết kế này kế thừa và bổ sung các logic nghiệp vụ (Business Rules) từ `02-category-list-plan.md` tập trung vào luồng cập nhật dữ liệu (Update/Edit). Tạm thời giả định hệ thống Auth Middleware đã được xử lý chung).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Sử dụng lại GORM Entity `Category` đã định nghĩa. Trong luồng chỉnh sửa, chúng ta quan tâm đặc biệt đến các logic kiểm tra dữ liệu trước khi thực hiện UPDATE.

```go
package entity

// Kế thừa struct Category từ các plan trước
// Các trường cần lưu ý đặc biệt khi cập nhật:
// - Slug: Có thể bị đổi, cần kiểm tra trùng lặp (uniqueIndex)
// - ParentID: Cần có logic chống vòng lặp đệ quy (Circular Dependency)
// - Status: Có logic kiểm tra khóa ngoại (ví dụ: đang chứa sản phẩm hay không)
```

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Cập nhật thông tin danh mục
- **Method & Route:** `PUT /api/v1/admin/categories/:id`
- **Auth:** Yêu cầu đi qua Middleware Auth (Admin/Menu Planner). Lấy `userId` từ context (nếu cần track người cập nhật `UpdatedBy`).

**Request Binding Struct:**
```go
package dto

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Slug        *string `json:"slug" binding:"omitempty,max=255"` // Cần check unique dưới DB (trừ chính nó)
	ParentID    *uint   `json:"parentId"` // Nullable
	Description *string `json:"description"`
	ImageURL    *string `json:"imageUrl"`
	SortOrder   *int    `json:"sortOrder" binding:"omitempty,min=0"`
	Status      *string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

**Response Payload (Success - 200 OK):**
```go
type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
```

### 2. Các Quy tắc Nghiệp vụ (Business Rules) tại tầng UseCase / Service

Trước khi gọi `db.Save()` hoặc `db.Updates()`, tầng Service cần thực hiện các validate sau:

1. **Kiểm tra tồn tại:** Truy vấn `ID` có tồn tại trong hệ thống không.
2. **Kiểm tra trùng lặp Slug:** Nếu `Slug` được truyền lên và khác `Slug` hiện tại, tiến hành check duplicate trong Database:
   ```go
   var count int64
   db.Model(&Category{}).Where("slug = ? AND id != ?", req.Slug, id).Count(&count)
   if count > 0 {
       return errors.New("Slug đã tồn tại") // Trả về HTTP 409 Conflict hoặc 400 Bad Request
   }
   ```
3. **Chống vòng lặp danh mục (Circular Dependency):**
   - Không được phép set `ParentID` bằng chính `ID` của danh mục (`req.ParentID == id`).
   - Nếu `ParentID` khác null, phải query đệ quy (hoặc dùng các kỹ thuật query nested sets / lấy toàn bộ tree) để đảm bảo `ParentID` mới KHÔNG nằm trong danh sách các con cháu (Descendants) của danh mục hiện tại.
4. **Kiểm tra ràng buộc trạng thái Ẩn (HIDDEN):**
   - Nếu Client truyền lên `Status = HIDDEN` trong khi DB đang là `ACTIVE`:
     - Kiểm tra xem có Product nào đang link với Category này không.
     - Kiểm tra xem có Category con nào đang link với Category này không.
     *(Logic này Frontend đã cảnh báo, nhưng Backend vẫn nên có option để chặn cứng hoặc cho phép ẩn theo nghiệp vụ chốt cuối cùng. Ở MVP có thể cho phép ẩn và Frontend hiển thị cảnh báo, hoặc Backend chặn cứng nếu yêu cầu chặt chẽ).*

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache (Cache Invalidation)
Giao diện khách hàng phụ thuộc rất nhiều vào danh sách Category hiển thị dạng cây. Do đó, ngay sau khi lệnh UPDATE thành công, bắt buộc phải xóa cache để dữ liệu mới nhất được phản ánh ngay.

- **Key Pattern:** `techbite:categories:active_tree` (và các key list khác liên quan).
- **Logic tại UseCase:**
  ```go
  // Thực hiện update trên MySQL
  err := u.categoryRepo.Update(ctx, id, updateData)
  if err != nil {
      return err
  }
  
  // Xóa cache (Fire and forget hoặc Sync tùy độ quan trọng)
  u.redisClient.Del(ctx, "techbite:categories:active_tree")
  ```

---
*Ghi chú: API upload ảnh (`POST /api/v1/media/upload`) là một module độc lập xử lý Multipart-Form data, trả về URL của ảnh. Frontend sẽ gọi API này trước, lấy URL và gắn vào payload của API PUT.*
