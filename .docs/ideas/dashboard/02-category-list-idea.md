# Ý TƯỞNG: Quản trị danh mục sản phẩm (Admin Category List)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Quản trị danh mục sản phẩm (Admin Category List).
* **Mục đích:** Quản lý toàn bộ danh mục món ăn/đồ uống trong hệ thống: cấu trúc phân cấp (Danh mục cha/con), thứ tự hiển thị trên menu, thêm mới, chỉnh sửa, xóa, thay đổi trạng thái và theo dõi số lượng sản phẩm thuộc từng danh mục.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Nhân viên thiết lập thực đơn (Menu Planner), Quản lý cửa hàng.
* **Hành động chính:**
* Tìm kiếm nhanh danh mục theo tên hoặc mã định danh.
* Lọc theo trạng thái hoạt động hoặc cấp bậc danh mục (Danh mục cha/con).
* Sắp xếp, kéo thả đổi vị trí (Sort Order) hiển thị trên giao diện người dùng (App/Web Order).
* Thêm mới, chỉnh sửa thông tin (Tên, Slug, Ảnh đại diện, Mô tả, Danh mục cha).
* Bật/Tắt trạng thái hoạt động (ẩn danh mục sẽ ẩn toàn bộ sản phẩm thuộc danh mục đó trên app khách hàng).
* Theo dõi nhanh số lượng sản phẩm đang liên kết trong danh mục.


* **Cảm xúc mang lại:** Gọn gàng, trực quan, cấu trúc phân cấp rõ ràng, dễ dàng thao tác kéo thả và điều chỉnh menu nhanh chóng khi có sự thay đổi về mùa hoặc chiến dịch bán hàng.

---

## 3. Đặc tả Thiết kế (Design Specs)

### Phong cách UI

* Dashboard hiện đại, tối giản, tối ưu hóa hiển thị dạng cây (Tree-view) hoặc bảng lồng nhau (Nested Table) nếu có danh mục con.
* Sử dụng Card và Table có khả năng mở rộng/thu gọn (Expand/Collapse).
* Bo góc nhẹ (`rounded-xl` hoặc `rounded-2xl`), đường viền mảnh tinh tế.

### Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu (Nút thêm mới, Lưu):** `bg-[#ff8c42]`
* **Đỏ (Xóa danh mục):** `bg-red-500`
* **Xanh lá (Đang hoạt động - Active):** `bg-green-500`
* **Xám (Đang ẩn - Hidden):** `bg-gray-400`
* **Xanh dương (Số lượng sản phẩm):** `bg-blue-500`

### Cấu trúc màn hình (Top to Bottom)

#### Header

* Tiêu đề: "Quản lý danh mục"
* Breadcrumb: `Trang chủ > Cấu hình hệ thống > Danh mục`
* Nhóm nút hành động:
* Nút **Thêm danh mục mới** `bg-[#ff8c42]`
* Nút **Sắp xếp nhanh (Reorder Mode)** (Bật chế độ kéo thả để đổi vị trí menu)



#### Thống kê nhanh (Summary Cards)

Hiển thị 4 Card:

* Tổng số danh mục
* Danh mục đang hoạt động (Active)
* Danh mục đang ẩn (Hidden)
* Tổng số sản phẩm đã được phân loại (Uncategorized check)

#### Thanh tìm kiếm & Bộ lọc

* Ô tìm kiếm danh mục theo tên hoặc Slug.
* Bộ lọc Cấp độ danh mục (Danh mục cha gốc / Danh mục con).
* Bộ lọc Trạng thái (Hoạt động / Ẩn).
* Nút Reset Filter.

#### Danh sách danh mục (Dạng Table hỗ trợ Tree-view)

Các cột hiển thị:

* **Icon kéo thả:** Xuất hiện khi bật chế độ sắp xếp.
* **Ảnh đại diện (Thumbnail/Icon):** Ảnh minh họa nhỏ cho danh mục trên menu app.
* **Tên danh mục:** Thể hiện phân cấp thụt lề nếu là danh mục con (ví dụ: `Burger` -> `— Burger Bò`).
* **Slug:** Đường dẫn thân thiện (ví dụ: `burger-bo`).
* **Số lượng sản phẩm:** Số sản phẩm hiện có trong danh mục (Click vào sẽ link sang trang danh sách sản phẩm đã được lọc theo danh mục này).
* **Thứ tự hiển thị:** Số thứ tự ưu tiên xuất hiện trên Menu (Sort Order).
* **Trạng thái:** Badge màu (Hoạt động / Ẩn).
* **Ngày cập nhật:** Thời gian cập nhật gần nhất.
* **Thao tác:** Xem chi tiết, Chỉnh sửa, Thêm danh mục con, Xóa.

#### Thao tác hàng loạt (Bulk Action)

Khi chọn nhiều danh mục:

* Đổi trạng thái hàng loạt (Ẩn/Hiện).
* Xóa các danh mục đã chọn (kèm điều kiện ràng buộc nếu có sản phẩm bên trong).

---

## 4. Dữ liệu cốt lõi (Mock Data)

Dưới đây là bộ dữ liệu danh mục mẫu dành riêng cho dự án F&B TechBite:

| Ảnh đại diện (Unsplash) | Tên danh mục | Danh mục cha | Slug | Số sản phẩm | Thứ tự | Trạng thái |
| --- | --- | --- | --- | --- | --- | --- |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **Burger** | *Không (Gốc)* | `burger` | 15 | 1 | Hoạt động |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **Pizza** | *Không (Gốc)* | `pizza` | 22 | 2 | Hoạt động |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **Đồ uống** | *Không (Gốc)* | `do-uong` | 18 | 3 | Hoạt động |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **— Cà phê** | Đồ uống | `ca-phe` | 8 | 1 | Hoạt động |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **— Trà trái cây** | Đồ uống | `tra-trai-cay` | 10 | 2 | Hoạt động |
| `[https://images.unsplash.com/](https://images.unsplash.com/)...` | **Món chay** | *Không (Gốc)* | `mon-chay` | 0 | 4 | Ẩn |

---

## 5. Chức năng nổi bật cho Danh mục

* **Kéo thả sắp xếp (Drag & Drop Reordering):** Cho phép admin kéo thả trực tiếp các dòng để sắp xếp thứ tự hiển thị danh mục trên ứng dụng gọi món của khách hàng.
* **Ràng buộc xóa (Cascade Delete Protection):** Khi xóa một danh mục:
* Nếu có sản phẩm bên trong: Yêu cầu admin chọn chuyển sản phẩm sang danh mục khác hoặc đưa về trạng thái "Chưa phân loại" trước khi xóa.
* Nếu có danh mục con: Cảnh báo xóa luôn danh mục con hoặc đưa danh mục con lên làm danh mục cha gốc.


* **Tự động tạo Slug:** Khi admin nhập tên danh mục "Gà Rán Giòn", hệ thống tự động sinh ra slug `ga-ran-gion` (có thể chỉnh sửa thủ công).
* **Badge trạng thái thông minh:** Trạng thái ẩn của danh mục cha sẽ tự động áp dụng (inherited) xuống các danh mục con.

---

## 6. Định hướng UI Component

* **Category Tree Table:** Bảng dữ liệu hỗ trợ cấu trúc phân cấp (mở rộng/thu gọn danh mục con) kèm icon kéo thả.
* **Create/Edit Category Slide-over (hoặc Modal):** Form thêm/sửa danh mục hiển thị từ cạnh phải màn hình để không làm gián đoạn trải nghiệm của Admin.
* **Image Uploader:** Khu vực tải lên ảnh đại diện của danh mục dạng kéo thả, hỗ trợ crop ảnh vuông tối ưu cho app di động.
* **Confirm Dialog (Safe Delete):** Hộp thoại xác nhận xóa danh mục thông minh, đưa ra các tùy chọn giải quyết các sản phẩm liên kết bên trong.
* **Summary Cards & Filter Panel.**