# FRONTEND PLAN: Quản trị danh mục sản phẩm (Admin Category List)

## 1. Mục tiêu
- Xây dựng màn hình Quản lý danh mục sản phẩm (Category List) dựa trên ý tưởng tại `02-category-list-idea.md` và thiết kế bảng (Table Card) từ mockup `dash-product/index.html`.
- Chuyển đổi mã HTML tĩnh thành các component React sử dụng Refine.js và Material UI (MUI), tận dụng tối đa các hook của Refine (`useTable`, `useForm`, `useShow`).
- Xử lý cấu trúc phân cấp (Danh mục cha/con) và hỗ trợ thao tác kéo thả (Drag & Drop) để sắp xếp thứ tự hiển thị.

## 2. Phạm vi chức năng
- **Hiển thị danh sách:** Hiển thị danh mục theo dạng bảng (Table), hỗ trợ Tree-view (mở rộng danh mục con).
- **Tìm kiếm & Lọc:** Tìm kiếm theo tên/slug (có debounce), lọc theo trạng thái (Hoạt động/Ẩn) và cấp bậc (Cha/Con).
- **Thao tác nhanh:** Thêm, Sửa (qua Drawer), Xóa (kèm cảnh báo cascade delete), Đổi trạng thái.
- **Sắp xếp:** Cho phép kéo thả trực tiếp trên bảng để đổi vị trí (Sort Order).
- **Thống kê:** Hiển thị tóm tắt tổng số danh mục, số hoạt động, số ẩn, số sản phẩm chưa phân loại.

## 3. Luồng màn hình (Main Content)
Chỉ thay thế phần `Main Content` trong layout, kế thừa layout đã có của dự án.

### 3.1. Header & Page Title
- Tiêu đề: `Category List` (kế thừa class `text-[32px] font-bold tracking-tight text-dark`).
- Breadcrumb: Trang chủ > Cấu hình hệ thống > Danh mục.
- Nút hành động chính:
  - `Thêm danh mục mới` (Màu cam thương hiệu hoặc primary của MUI theme).
  - `Sắp xếp nhanh (Reorder Mode)` (Chuyển bảng sang chế độ kéo thả).

### 3.2. 

### 3.3. Khu vực Tìm kiếm và Bộ lọc (Table Header Tools)
Tận dụng lại khối header của Card trong mockup:
- Khung tìm kiếm: Dùng `TextField` MUI kết hợp icon search, style bo tròn (`rounded-full`, `bg-lightbg`).
- Bộ lọc Cấp độ & Trạng thái: Dùng MUI `Select` hoặc `ToggleButtonGroup`.
- Nút Reset Filter.

### 3.4. Bảng danh sách (Table Card)
Kế thừa style bảng từ HTML mockup (`bg-white rounded-[14px] border border-bordercolor`):
- Sử dụng `DataGrid` của MUI (Premium/Pro nếu có hỗ trợ TreeData, hoặc tự custom Row với độ thụt lề thụt lùi `paddingLeft`).
- Các cột:
  - **Kéo thả:** Icon drag (chỉ hiện khi bật chế độ sắp xếp).
  - **Hình ảnh:** Tương tự cột Image trong mockup (kích thước 60x60, bo góc).
  - **Tên danh mục:** Nếu là danh mục con thì có icon/khoảng trắng thụt lề.
  - **Slug:** Đường dẫn tĩnh.
  - **Số lượng sản phẩm:** Cột hiển thị số sản phẩm kèm link điều hướng sang trang Product.
  - **Thứ tự:** Sort Order.
  - **Trạng thái:** Status Badge (Màu xanh/Xám tương tự).
  - **Thao tác (Action):** Tái sử dụng style các nút sửa/xóa trong mockup (Nút vuông bo góc `w-8 h-8 bg-gray-100`).

## 4. Component Breakdown
### 4.1. Page Container
- `CategoryListPage`: Bọc `List` của Refine, quản lý state cho table, filter, và chế độ kéo thả.

### 4.2. UI Blocks
- `CategorySummaryCards`: Component hiển thị 4 thẻ thống kê.
- `CategoryFilterBar`: Component tìm kiếm & lọc dữ liệu.
- `CategoryTreeTable`: Bảng MUI DataGrid custom hiển thị dữ liệu phân cấp.
- `CategoryFormDrawer`: Drawer (Slide-over) hiện từ bên phải chứa Form thêm/sửa danh mục (tận dụng `useDrawerForm`).
- `ImageUploader`: Upload ảnh đại diện, hỗ trợ hiển thị preview dạng crop vuông.
- `CategoryConfirmDelete`: Dialog xác nhận xóa thông minh (cảnh báo nếu có sản phẩm/danh mục con liên kết).

### 4.3. Shared Primitives (Tái sử dụng)
- `StatusBadge` (Tái sử dụng từ Product).
- `TableActionButtons` (Các nút Edit/Delete lấy từ mockup).

## 5. Data & State Management
### 5.1. Refine Resources
- Sử dụng resource `categories`.

### 5.2. State quản lý
- `searchKeyword`: Từ khóa tìm kiếm (Debounce 300ms).
- `filters`: Lọc trạng thái / Cấp độ.
- `isReorderMode`: Boolean bật tắt chế độ kéo thả.
- `drawerState`: Trạng thái mở form Thêm/Sửa.

## 6. Mapping Mockup sang Material UI (MUI)
- **Container / Card:** Sử dụng `Paper` hoặc `Box` với `sx={{ borderRadius: '14px', border: '1px solid #D5D5D5' }}`.
- **Search Input:** `TextField` với `InputProps={{ startAdornment: <SearchIcon /> }}`.
- **Table Header (`th`):** Tùy chỉnh `headerClassName` của DataGrid để in đậm, font-size nhỏ.
- **Table Cell (`td`):** Tùy chỉnh `cellClassName` hoặc `renderCell` để hiện text font-semibold.
- **Hình ảnh:** Dùng `Box` hoặc `Avatar` kích thước 60x60, `borderRadius: '8px'`.
- **Nút Actions:** Dùng `IconButton` của MUI, set màu nền xám nhạt (`#F3F4F6`), hover tối hơn.

## 7. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Tự động sinh Slug:** Khi nhập tên ở `CategoryFormDrawer`, field Slug sẽ tự cập nhật bằng cách gọi hàm `slugify`.
- **Ràng buộc xóa:** Gọi API kiểm tra trước khi xóa. Nếu báo lỗi ràng buộc thì hiển thị thông báo thay vì xóa trực tiếp.
- **Tính năng kéo thả (Reorder):** Sử dụng thư viện ngoài như `@dnd-kit/core` hoặc tính năng Row Reorder của DataGrid (nếu dùng bản Pro) để cập nhật trường `sort_order` và gọi API update.

## 8. Checklist Triển khai
- [ ] Khai báo resource `categories` trong Refine.
- [ ] Dựng UI Page layout và Summary Cards.
- [ ] Dựng Bảng DataGrid hiển thị danh sách phẳng (Flat list).
- [ ] Xử lý logic Tree-view (Nhóm danh mục con theo danh mục cha).
- [ ] Áp dụng style (UI mockup `dash-product`) vào DataGrid của MUI.
- [ ] Thêm Toolbar và Filter (Debounce Search).
- [ ] Dựng component `CategoryFormDrawer` cho Thêm/Sửa.
- [ ] Tích hợp tính năng xóa (Cascade warning).
- [ ] Tích hợp kéo thả đổi vị trí (Tùy chọn nâng cao).
