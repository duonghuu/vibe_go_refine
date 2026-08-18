# KẾ HOẠCH FRONTEND: Quản lý danh mục bài viết (Post Category List)

## 1. Thông tin chung
- **Tính năng:** Màn hình danh sách Danh mục bài viết.
- **Vị trí file dự kiến:** `apps/frontend/src/pages/post-categories/list.tsx`
- **Resource Refine:** `post-categories`
- **Mục đích:** Hiển thị, lọc, tìm kiếm và quản lý các danh mục bài viết theo từng loại bài viết (Post Type).

## 2. Bố cục màn hình (Layout)
- Sử dụng layout cơ bản với Header (Title + Breadcrumb) và Container chứa thẻ Card bọc `DataGrid` của Material UI.
- **Tiêu đề:** "Danh sách danh mục bài viết"

### 2.1. Thanh công cụ (Toolbar / Header Actions)
- **Search box:** Ô input tìm kiếm theo tên danh mục (có icon Search, debounce 300ms).
- **Filter Dropdown (Post Type):** Dropdown bắt buộc chọn loại bài viết (`typeCode`). Mặc định chọn một loại (vd: NEWS).
- **Filter Dropdown (Status):** Dropdown lọc theo trạng thái (`ACTIVE`, `INACTIVE`).
- **Nút "Thêm mới":** Liên kết sang trang `/post-categories/create?typeCode={typeCode}`.

### 2.2. Bảng dữ liệu (DataGrid)
Các cột hiển thị:
1. **Hình ảnh (Image):** Avatar vuông/bo góc (60x60).
2. **Tên danh mục (Name):** Hiển thị dạng text, nếu là danh mục con (có `parentId`) thì lùi đầu dòng và thêm tiền tố (vd: `— Tin thể thao`).
3. **Mã (Slug):** Text.
4. **Thứ tự (Sort Order):** Text/Số.
5. **Trạng thái (Status):** MUI Chip (`ACTIVE` màu xanh lá, `INACTIVE` màu xám).
6. **Hành động (Actions):** 
   - Nút Chỉnh sửa (điều hướng sang `/post-categories/edit/:id`).
   - Nút Xóa (Hiển thị popup xác nhận).

## 3. Cấu hình Refine & Hooks
- Sử dụng hook `useDataGrid<IPostCategory, HttpError>({ resource: "post-categories", syncWithLocation: true })`.
- Do đặc thù bảng danh mục cần hiển thị theo `typeCode`, bộ lọc `filters` khởi tạo ban đầu phải đi kèm điều kiện `typeCode` hợp lệ.
- Sử dụng `useSelect` để lấy danh sách `post-types` phục vụ cho Filter Dropdown (Post Type).

## 4. Ràng buộc & Tối ưu (Constraints & Optimization)
- **Tối ưu gọi API:** Xử lý debounce 300ms bằng `lodash/debounce` cho ô Search.
- **Type Safety:** Định nghĩa interface `IPostCategory` cho Refine tương ứng với cấu trúc dữ liệu Backend trả về.
- **Design System:** Tuân thủ `STYLEGUIDE.md` (không dùng hex colors tự định nghĩa cứng, sử dụng shadow và border radius theo MUI Theme).
- **Bảo mật & UI:** Render skeleton hoặc loading progress bar trên DataGrid khi API đang xử lý.
