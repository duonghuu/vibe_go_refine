# QUY HOẠCH KIẾN TRÚC BACK-END: CHỈNH SỬA SẢN PHẨM

**Dự án:** TechBite
**Module:** Product Management
**Tài liệu tham chiếu:** `.docs/ideas/dashboard/09-product-edit-idea.md`

*(Ghi chú: Bản thiết kế này tập trung vào luồng cập nhật dữ liệu (Update) cho sản phẩm, đồng bộ với kiến trúc phân tầng DDD/Clean Architecture đang có).*

---

## TRỤ CỘT 1: THIẾT KẾ DỮ LIỆU (DATABASE SCHEMA STRUCT)

Cấu trúc Entity `Product` giữ nguyên như thiết kế trong phần tạo mới sản phẩm, sử dụng GORM.

---

## TRỤ CỘT 2: GIAO KÈO API (API CONTRACT & CONTEXT AUTH)

### 1. Lấy chi tiết sản phẩm (Dành cho form Edit fill data)
- **Method & Route:** `GET /api/v1/admin/products/:id`
- **Auth:** Yêu cầu đi qua Middleware Auth (Kiểm tra token hợp lệ).

**Response Payload (Success - 200 OK):** Trả về struct `ProductResponse` tương tự tạo mới.

### 2. Cập nhật sản phẩm
- **Method & Route:** `PUT /api/v1/admin/products/:id` (hoặc `PATCH`)
- **Auth:** Yêu cầu đi qua Middleware Auth (Kiểm tra token hợp lệ, lấy `userId`, phân quyền Admin/Operator).

**Request Binding Struct:**
```go
package dto

type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"required,max=255"`
	Description string   `json:"description"`
	CategoryID  uint     `json:"categoryId" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	SalePrice   *float64 `json:"salePrice" binding:"omitempty,gte=0"`
	SKU         string   `json:"sku"`
	Stock       int      `json:"stock" binding:"omitempty,gte=0"`
	ImageURL    string   `json:"image"`
	Status      string   `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}
```

### 3. Các Quy tắc Nghiệp vụ (Business Rules) tại tầng UseCase / Service

Khi cập nhật một product, Service layer (UseCase) cần xử lý các luồng logic (Business validations) sau:
1. **Kiểm tra tồn tại:**
   Gọi Repository kiểm tra xem Product ID có tồn tại hay không. Báo lỗi 404 (Not Found) nếu không thấy.
2. **Kiểm tra trùng lặp SKU:**
   Nếu SKU bị thay đổi, phải kiểm tra xem SKU mới có bị trùng với sản phẩm KHÁC hay không.
3. **Kiểm tra Validate Category:**
   Gọi Repository kiểm tra xem `CategoryID` gửi lên có thật sự tồn tại trong DB không. Báo lỗi 404/400 nếu danh mục không hợp lệ.
4. **Kiểm tra logic Giá khuyến mãi:**
   Nếu client có truyền lên `SalePrice`, bắt buộc check: `SalePrice < Price`.
5. **Giữ nguyên các trường hệ thống:**
   Không cập nhật `SoldCount`, `CreatedAt`, `DeletedAt` từ payload của người dùng.

---

## TRỤ CỘT 3: XỬ LÝ CACHE & QUẢN LÝ TRẠNG THÁI (REDIS INTEGRATION)

### Xóa Cache (Cache Invalidation)
Mỗi khi có thông tin sản phẩm bị cập nhật, cần xóa cache để đảm bảo dữ liệu mới nhất được hiển thị trên Frontend End-user.

- **Quy trình xử lý tại UseCase:**
  ```go
  // 1. Cập nhật sản phẩm trong DB
  err := u.productRepo.Update(ctx, productID, productEntity)
  if err != nil {
      return err
  }
  
  // 2. Xóa cache
  // Xóa cache chi tiết sản phẩm
  u.redisClient.Del(ctx, fmt.Sprintf("techbite:products:detail:%d", productID))
  // Xóa cache danh sách tổng
  u.redisClient.Del(ctx, "techbite:products:home_list")
  // Xóa cache danh sách theo CategoryID cũ và mới (nếu thay đổi Category)
  u.redisClient.Del(ctx, fmt.Sprintf("techbite:categories:products:%d", oldCategoryID))
  u.redisClient.Del(ctx, fmt.Sprintf("techbite:categories:products:%d", productEntity.CategoryID))
  ```
