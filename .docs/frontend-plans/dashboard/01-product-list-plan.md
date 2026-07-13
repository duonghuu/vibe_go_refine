# FRONTEND PLAN: Quản trị danh sách sản phẩm

## 1. Mục tiêu
- Xây dựng màn hình quản trị sản phẩm cho TechBite với khả năng tìm kiếm, lọc, sắp xếp và thao tác nhanh trên danh sách lớn.
- Ưu tiên trải nghiệm vận hành cho Admin/Operator/Store Manager: giảm số click, hiển thị rõ trạng thái, thao tác hàng loạt thuận tiện.
- Bám theo kiến trúc hiện tại của frontend: Refine.js + Material UI + DataGrid.

## 2. Phạm vi chức năng
### 2.1. Danh sách sản phẩm
- Hiển thị danh sách sản phẩm dạng bảng.
- Hỗ trợ phân trang, sort, filter, search.
- Hiển thị thông tin cốt lõi: ảnh, tên, SKU, danh mục, giá, giá khuyến mãi, tồn kho, đã bán, trạng thái, ngày cập nhật.

### 2.2. Thao tác trên từng sản phẩm
- Xem chi tiết.
- Chỉnh sửa.
- Sao chép.
- Ẩn/Hiện.
- Xóa.

### 2.3. Thao tác hàng loạt
- Chọn nhiều sản phẩm bằng checkbox.
- Xóa hàng loạt.
- Đổi trạng thái hiển thị.
- Chuyển danh mục.
- Xuất Excel.

### 2.4. Tạo mới và chỉnh sửa
- Tạo sản phẩm mới từ màn hình danh sách hoặc trang riêng tùy độ dài form.
- Chỉnh sửa thông tin sản phẩm.
- Sử dụng Drawer cho form ngắn, chuyển sang trang riêng nếu form dài hoặc có nhiều trường nâng cao.

## 3. Luồng màn hình
### 3.1. Header
- Tiêu đề: `Quản lý sản phẩm`.
- Breadcrumb điều hướng.
- Nút hành động chính:
  - `Thêm sản phẩm`
  - `Import`
  - `Export`

### 3.2. Khối thống kê nhanh
- Tổng số sản phẩm.
- Đang bán.
- Ngừng bán.
- Sắp hết hàng.

### 3.3. Khu vực tìm kiếm và lọc
- Ô tìm kiếm theo tên, SKU, mã sản phẩm.
- Filter theo danh mục.
- Filter theo trạng thái.
- Filter theo khoảng giá.
- Filter theo tồn kho.
- Filter theo khoảng thời gian tạo/cập nhật.
- Nút `Reset Filter`.

### 3.4. Bảng danh sách
- Sử dụng `DataGrid` của MUI.
- Có cột checkbox chọn hàng.
- Có cột ảnh sản phẩm.
- Các cột dữ liệu chính và cột thao tác.
- Cột thao tác ở cuối bảng.

### 3.5. Phân trang
- Chọn số dòng: 10 / 20 / 50 / 100.
- Điều hướng trang.
- Hiển thị tổng số bản ghi.

## 4. Component Breakdown
### 4.1. Page container
- `ProductListPage`
- Chịu trách nhiệm lấy dữ liệu, điều phối filter/sort/pagination, mở drawer, confirm dialog và bulk actions.

### 4.2. UI blocks
- `ProductSummaryCards`
- `ProductFilterBar`
- `ProductTable`
- `ProductRowActions`
- `BulkActionToolbar`
- `ProductFormDrawer`
- `ConfirmDeleteDialog`
- `EmptyState`
- `LoadingSkeleton`
- `ToastNotification`

### 4.3. Shared primitives
- `StatusBadge`
- `PriceCell`
- `InventoryCell`
- `ProductImageCell`

## 5. Data & State
### 5.1. Resource chính
- `products`
- `categories`

### 5.2. Dữ liệu cần lấy
- `id`
- `name`
- `sku`
- `category`
- `price`
- `salePrice`
- `stock`
- `soldCount`
- `status`
- `image`
- `createdAt`
- `updatedAt`

### 5.3. State quản lý trên frontend
- Search keyword.
- Bộ lọc danh mục/trạng thái/giá/tồn kho/thời gian.
- Sort model.
- Row selection cho bulk actions.
- Trạng thái mở/đóng drawer.
- Trạng thái mở/đóng confirm dialog.
- Trạng thái toast feedback.

## 6. Quy ước hiển thị
### 6.1. Trạng thái sản phẩm
- `Đang bán`: chip xanh.
- `Ẩn`: chip xám.
- `Hết hàng`: chip vàng hoặc đỏ nhạt tuỳ mức cảnh báo.
- `Sắp hết hàng`: chip vàng.

### 6.2. Giá và tồn kho
- Giá gốc và giá khuyến mãi cần tách biệt rõ.
- Nếu không có giá khuyến mãi thì hiển thị `-`.
- Tồn kho thấp cần có cảnh báo trực quan.
- Số đã bán hiển thị theo format ngắn gọn, dễ scan.

### 6.3. Hình ảnh
- Dùng thumbnail kích thước nhỏ, bo góc nhẹ.
- Nếu thiếu ảnh thì dùng placeholder có icon.

## 7. Tương tác chính
### 7.1. Tìm kiếm và lọc
- Search được debounce 300ms.
- Filter có thể kết hợp nhiều điều kiện.
- Reset filter đưa bảng về trạng thái mặc định.

### 7.2. Sắp xếp
- Mới cập nhật.
- Mới tạo.
- Giá tăng/giảm.
- Tên A-Z.
- Tồn kho.
- Số lượng bán.

### 7.3. Thao tác nguy hiểm
- Xóa sản phẩm phải mở confirm dialog.
- Ẩn/Hiện sản phẩm nên có confirm nếu là thay đổi trạng thái quan trọng.
- Sau khi thao tác thành công, hiển thị toast.

### 7.4. Bulk actions
- Chỉ hiện khi đã chọn ít nhất 1 dòng.
- Có thể hiển thị toolbar cố định phía trên bảng.

## 8. Trạng thái giao diện
- `Loading`: skeleton cho summary cards và table.
- `Empty`: không có dữ liệu hoặc không có kết quả lọc.
- `Error`: hiển thị message rõ ràng và nút thử lại.
- `Success`: toast xác nhận thao tác.

## 9. Gợi ý triển khai kỹ thuật
- Dùng `useDataGrid` của Refine cho danh sách.
- Dùng `DataGrid` để tận dụng pagination/sort/selection.
- Dùng `useMany` nếu cần resolve danh mục từ `categoryId`.
- Dùng `useDrawerForm` cho form ngắn tạo/sửa nhanh ngay tại trang list.
- Dùng `useModal` hoặc `ConfirmDialog` cho xóa và bulk actions nguy hiểm.
- Dùng `DateField` hoặc format utility để hiển thị thời gian nhất quán.
- Tách phần UI thuần ra khỏi phần data hook để dễ mở rộng.

## 10. Định tuyến đề xuất
- `/products`
- `/products/create`
- `/products/edit/:id`
- `/products/show/:id`

## 11. Tiêu chí hoàn thành
- Người dùng có thể xem, tìm, lọc, sort và phân trang danh sách sản phẩm.
- Người dùng có thể tạo/sửa/xóa/ẩn-hiện sản phẩm.
- Người dùng có thể chọn nhiều sản phẩm và thao tác hàng loạt.
- Bảng hiển thị rõ ràng trên desktop và tablet.
- Có loading, empty state, confirm dialog và toast feedback.

## 12. Checklist triển khai
- [ ] Tạo resource `products` trong cấu hình frontend.
- [ ] Tạo page `ProductList`.
- [ ] Tạo summary cards.
- [ ] Tạo filter bar có debounce.
- [ ] Tạo bảng DataGrid với row actions.
- [ ] Tạo drawer form cho create/edit.
- [ ] Tạo confirm dialog cho delete.
- [ ] Tạo bulk action toolbar.
- [ ] Tạo empty/loading states.
- [ ] Kiểm tra responsive và khoảng cách hiển thị.

## 13. Ghi chú mở rộng
- Có thể mở rộng thêm import/export theo file CSV/XLSX.
- Có thể thêm quản lý biến thể sản phẩm trong giai đoạn sau.
- Có thể bổ sung badge tồn kho theo ngưỡng cấu hình.
