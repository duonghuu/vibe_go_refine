# Ý TƯỞNG: Quản trị Thêm mới sản phẩm (Admin Create Product)

## 1. Thông tin chung (Meta Info)

* **Dự án:** TechBite.
* **Tính năng:** Màn hình Thêm mới sản phẩm (Admin Create Product).
* **Mục đích:** Cho phép Admin/Operator nhập thông tin chi tiết để tạo mới một sản phẩm vào hệ thống.

---

## 2. Đối tượng & Trải nghiệm (Target & UX)

* **Người dùng chính:** Quản trị viên (Admin), Nhân viên vận hành (Operator), Quản lý cửa hàng.
* **Hành động chính:**
  * Nhập thông tin cơ bản: Tên sản phẩm, mô tả, chọn danh mục.
  * Tải lên hình ảnh đại diện của sản phẩm.
  * Nhập thông tin định giá: Giá gốc, giá khuyến mãi.
  * Nhập thông tin kho hàng: SKU (Mã sản phẩm), số lượng tồn kho.
  * Thiết lập trạng thái hiển thị của sản phẩm (Đang bán, Ẩn).
* **Cảm xúc mang lại:** Form nhập liệu sạch sẽ, rõ ràng, được nhóm thành các khu vực logic giúp người dùng không bị quá tải. Thao tác lưu nhanh chóng và có phản hồi lỗi chính xác.

---

## 3. Đặc tả Thiết kế (Design Specs)

### Phong cách UI

* Bố cục dạng Grid/Card, chia thành 2 cột (Main Column và Side Column) để tận dụng không gian màn hình Desktop.
* Thiết kế nhất quán với danh sách sản phẩm: bo góc nhẹ (`rounded-xl` hoặc `rounded-2xl`), khoảng trắng hợp lý.
* Sticky Action Bar (Thanh công cụ cố định trên cùng hoặc dưới cùng) để người dùng luôn nhìn thấy nút "Lưu" mà không cần cuộn trang.
* Responsive: Layout 2 cột sẽ chuyển thành 1 cột trên màn hình Tablet/Mobile.

### Màu sắc chủ đạo (Brand Colors)

* **Cam thương hiệu:** Nút "Lưu", "Tạo mới" (Primary Action): `bg-[#ff8c42]`
* **Đỏ:** Báo lỗi nhập liệu (Validation Error Text/Border): `text-red-500`
* **Xám nhạt/Trắng:** Màu nền cho các Card chứa form.

### Cấu trúc màn hình (Top to Bottom)

#### Header & Sticky Actions
* Tiêu đề: "Thêm sản phẩm mới" hoặc "Chỉnh sửa sản phẩm: [Tên sản phẩm]"
* Breadcrumb: Bảng điều khiển > Sản phẩm > Thêm mới
* Nút hành động:
  * Nút **Hủy bỏ** (Quay về danh sách)
  * Nút **Lưu & Ẩn** (Lưu nháp)
  * Nút **Lưu & Xuất bản** (Primary Action)

#### Bố cục 2 Cột (2-Column Layout)

**Cột chính (Main Column - Chiếm khoảng 2/3 không gian):**

1. **Thông tin cơ bản (Basic Info Card)**
   * Tên sản phẩm (Input text - Bắt buộc).
   * Mô tả (Rich Text Editor để có thể in đậm, thêm bullet points...).

2. **Giá cả (Pricing Card)**
   * Giá bán gốc (Number Input - Bắt buộc).
   * Giá khuyến mãi (Number Input - Tùy chọn).

3. **Quản lý kho (Inventory Card)**
   * Mã sản phẩm / SKU (Input Text).
   * Số lượng tồn kho (Number Input - Mặc định 0).

**Cột phụ (Side Column - Chiếm khoảng 1/3 không gian):**

1. **Hình ảnh (Media Card)**
   * Khu vực tải ảnh chính lên (Kéo thả hoặc Click để chọn file).
   * Hiển thị ảnh xem trước (Preview) sau khi tải.

2. **Phân loại (Organization Card)**
   * Danh mục (Select Dropdown - Burger, Pizza, Gà rán...).

3. **Trạng thái (Status Card)**
   * Trạng thái hiển thị (Radio hoặc Switch/Toggle: Đang bán / Ẩn).

---

## 4. Kiểm tra dữ liệu (Validation & Error Handling)

* **Tên sản phẩm:** Bắt buộc nhập, tối đa 255 ký tự.
* **Danh mục:** Bắt buộc chọn.
* **Giá gốc:** Bắt buộc nhập, phải lớn hơn 0.
* **Giá khuyến mãi:** Nếu có nhập, bắt buộc phải nhỏ hơn **Giá gốc**.
* **SKU:** Phải duy nhất. Nếu để trống, hệ thống có thể tự sinh mã.
* **Tồn kho:** Phải lớn hơn hoặc bằng 0.
* **Hình ảnh:** Kiểm tra định dạng (chỉ nhận JPG, PNG, WEBP), dung lượng tối đa (ví dụ 5MB).

---

## 5. Chức năng nổi bật

* **Real-time Validation:** Báo lỗi ngay lập tức khi người dùng nhập sai hoặc bỏ trống trường bắt buộc và rời khỏi ô nhập đó (onBlur/onChange).
* **Drag & Drop Image:** Tải ảnh lên dễ dàng bằng thao tác kéo thả.
* **Unsaved Changes Warning:** Cảnh báo (Confirm Dialog) nếu người dùng cố gắng điều hướng sang trang khác khi đang nhập dở dữ liệu mà chưa lưu.
* **Toast Notification:** Hiển thị thông báo ở góc màn hình khi lưu thành công hoặc gặp lỗi kết nối.

---

## 6. Định hướng UI Component

* **Page Layout / Header:** Chứa Breadcrumb và các nút lưu ở vị trí cố định.
* **Card Component:** Bao bọc từng nhóm field (Basic Info, Pricing, v.v.).
* **Form Elements:**
  * `TextField` / `Input`
  * `Select` (dành cho Danh mục)
  * `Switch` / `RadioGroup` (dành cho Trạng thái)
* **Upload Component / Dropzone:** Giao diện tải ảnh.
* **Rich Text Editor:** Trình soạn thảo văn bản cho mô tả.
* **Toast / Snackbar:** Phản hồi kết quả thao tác.
* **Dialog:** Hộp thoại xác nhận khi rời trang.
