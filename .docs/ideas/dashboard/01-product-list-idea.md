# Ý TƯỞNG: Quản trị danh sách sản phẩm (Admin Product List)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Quản trị danh sách sản phẩm (Admin Product List).
* **Mục đích:** Quản lý toàn bộ sản phẩm trong hệ thống: tìm kiếm, lọc, xem thông tin, thêm mới, chỉnh sửa, xóa, thay đổi trạng thái hiển thị và theo dõi tồn kho.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Nhân viên vận hành (Operator), Quản lý cửa hàng.
* **Hành động chính:**

  * Tìm kiếm sản phẩm theo tên, SKU hoặc mã sản phẩm.
  * Lọc theo danh mục, trạng thái, khoảng giá và tồn kho.
  * Sắp xếp theo thời gian tạo, cập nhật, giá hoặc số lượng bán.
  * Thêm mới, chỉnh sửa, sao chép hoặc xóa sản phẩm.
  * Bật/Tắt trạng thái hiển thị.
  * Xem nhanh thông tin và tồn kho.
* **Cảm xúc mang lại:** Chuyên nghiệp, dễ quan sát, thao tác nhanh, giảm số lần click, phù hợp làm việc với số lượng dữ liệu lớn.

---

## 3. Đặc tả Thiết kế (Design Specs)

### Phong cách UI

* Dashboard hiện đại (Modern Admin Dashboard).
* Thiết kế tối giản, ưu tiên khả năng đọc dữ liệu.
* Sử dụng Card và Table kết hợp.
* Khoảng trắng hợp lý, bo góc nhẹ (`rounded-xl` hoặc `rounded-2xl`).
* Responsive cho Desktop và Tablet.

### Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu:** Các nút hành động chính (Thêm sản phẩm, Lưu): `bg-[#ff8c42]`
* **Đỏ:** Xóa, cảnh báo: `bg-red-500`
* **Xanh lá:** Đang bán, còn hàng: `bg-green-500`
* **Xám:** Ngừng bán, ẩn sản phẩm: `bg-gray-400`
* **Vàng:** Sắp hết hàng: `bg-yellow-500`

### Cấu trúc màn hình (Top to Bottom)

#### Header

* Tiêu đề "Quản lý sản phẩm"
* Breadcrumb
* Nút **Thêm sản phẩm**
* Nút **Import**
* Nút **Export**

#### Thống kê nhanh (Summary Cards)

Hiển thị 4 Card:

* Tổng số sản phẩm
* Đang bán
* Ngừng bán
* Sắp hết hàng

#### Thanh tìm kiếm & Bộ lọc

Bao gồm:

* Ô tìm kiếm
* Danh mục
* Trạng thái
* Khoảng giá
* Kho hàng
* Khoảng thời gian tạo/cập nhật
* Nút Reset Filter

#### Danh sách sản phẩm

Hiển thị dạng **Table**.

Các cột gồm:

* Checkbox
* Hình ảnh
* Tên sản phẩm
* SKU
* Danh mục
* Giá
* Giá khuyến mãi
* Tồn kho
* Đã bán
* Trạng thái
* Ngày cập nhật
* Thao tác

#### Thao tác từng dòng

* Xem chi tiết
* Chỉnh sửa
* Sao chép
* Ẩn/Hiện
* Xóa

#### Thao tác hàng loạt (Bulk Action)

Khi chọn nhiều sản phẩm:

* Xóa
* Đổi trạng thái
* Chuyển danh mục
* Xuất Excel

#### Phân trang

* Chọn số dòng (10 / 20 / 50 / 100)
* Chuyển trang
* Hiển thị tổng số bản ghi

---

## 4. Dữ liệu cốt lõi (Mock Data)

### Danh mục

* Burger
* Pizza
* Gà rán
* Cơm
* Mì
* Đồ uống
* Tráng miệng

### Danh sách sản phẩm

Mỗi sản phẩm gồm:

* Hình ảnh lấy từ **unsplash.com**
* Tên món ăn/ngước uống
* SKU ngẫu nhiên (VD: `TB-000123`)
* Danh mục
* Giá gốc
* Giá khuyến mãi
* Số lượng tồn kho
* Số lượng đã bán
* Trạng thái (Đang bán / Ẩn / Hết hàng)
* Ngày tạo
* Ngày cập nhật

---

## 5. Chức năng nổi bật

* Tìm kiếm theo nhiều tiêu chí.
* Lọc nhiều điều kiện cùng lúc.
* Sắp xếp theo:

  * Mới cập nhật
  * Mới tạo
  * Giá tăng/giảm
  * Tên A-Z
  * Tồn kho
  * Số lượng bán
* Chọn nhiều sản phẩm để thao tác hàng loạt.
* Badge màu hiển thị trạng thái.
* Hover hàng để hiển thị nhanh các nút thao tác.
* Loading Skeleton khi tải dữ liệu.
* Empty State khi không có kết quả.
* Confirm Dialog trước khi xóa.
* Toast Notification sau khi thao tác thành công hoặc thất bại.

---

## 6. Định hướng UI Component

* Summary Card
* Search Input
* Filter Panel (Collapsible)
* Data Table
* Badge Status
* Pagination
* Dropdown Action Menu
* Bulk Action Toolbar
* Confirm Dialog
* Loading Skeleton
* Empty State
* Toast Notification

Thiết kế ưu tiên trải nghiệm quản trị với dữ liệu lớn, tốc độ thao tác nhanh và dễ mở rộng cho các chức năng như Import/Export, phân quyền và quản lý biến thể sản phẩm trong tương lai.
