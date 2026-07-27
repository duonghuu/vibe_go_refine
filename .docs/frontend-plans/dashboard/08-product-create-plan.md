# FRONTEND PLAN: Thêm mới sản phẩm (Create Product)

## 1. Mục tiêu
- Xây dựng màn hình Thêm mới sản phẩm dựa trên tài liệu ý tưởng `08-product-create-idea.md`.
- Sử dụng các hook của Refine.js (như `useForm`) và các component Material UI (MUI) kết hợp `react-hook-form` để quản lý state và validation.
- Xây dựng layout 2 cột rõ ràng, nhóm thông tin theo card (Thông tin cơ bản, Giá cả, Quản lý kho, Hình ảnh, Phân loại, Trạng thái) giúp người dùng nhập liệu hiệu quả.

## 2. Phạm vi chức năng
- **Thông tin cơ bản:** Tên sản phẩm, Mô tả chi tiết (Rich Text Editor).
- **Giá cả:** Giá bán gốc, Giá khuyến mãi.
- **Quản lý kho:** Mã sản phẩm (SKU), Số lượng tồn kho.
- **Hình ảnh:** Tải lên ảnh sản phẩm (Kéo thả, xem trước).
- **Phân loại & Trạng thái:** Chọn danh mục sản phẩm từ dropdown, Cài đặt trạng thái hiển thị (Đang bán / Ẩn).
- **Hành động form:** Hủy bỏ (quay về danh sách), Lưu & Ẩn (Lưu nháp), Lưu & Xuất bản (Tạo mới và active).

## 3. Luồng màn hình (Main Content)

### 3.1. Header & Page Title
- **Tiêu đề:** `Thêm sản phẩm mới`
- **Breadcrumb:** `Bảng điều khiển > Sản phẩm > Thêm mới`

### 3.2. Form Layout (Sử dụng CSS Grid / Flexbox)
Bố cục 2 cột (Main Column chiếm khoảng 2/3, Side Column chiếm 1/3) trên màn hình Desktop, sẽ chuyển thành 1 cột trên Tablet/Mobile:

- **Cột chính (Main Column):**
  - Card "Thông tin cơ bản":
    - **Tên sản phẩm (*):** MUI `TextField`, bắt buộc.
    - **Mô tả:** Rich Text Editor (như MUI RTE, React-Quill, hoặc TipTap) cho phép in đậm, bullet, v.v.
  - Card "Giá cả":
    - **Giá bán gốc (*):** MUI `TextField` dạng Number, bắt buộc > 0.
    - **Giá khuyến mãi:** MUI `TextField` dạng Number, có thể trống, nếu nhập phải < Giá bán gốc.
  - Card "Quản lý kho":
    - **Mã sản phẩm / SKU:** MUI `TextField`. Nếu để trống, hệ thống sẽ tự sinh mã.
    - **Số lượng tồn kho:** MUI `TextField` dạng Number, mặc định 0.

- **Cột phụ (Side Column):**
  - Card "Hình ảnh":
    - Khu vực Dropzone upload ảnh. Có validate dung lượng (≤ 5MB) và định dạng. Có chức năng xem trước ảnh.
  - Card "Phân loại":
    - **Danh mục (*):** MUI `Select` hoặc `Autocomplete`. Sử dụng API để fetch danh sách các danh mục khả dụng.
  - Card "Trạng thái":
    - **Trạng thái hiển thị:** MUI `Switch` hoặc `RadioGroup` (Mặc định: Đang bán).

### 3.3. Sticky Footer Actions (Cố định ở dưới hoặc trên cùng)
- Nút **Hủy bỏ:** Hủy bỏ việc tạo, hiển thị dialog confirm nếu form đã bị thay đổi, điều hướng về `/products`.
- Nút **Lưu & Ẩn (Lưu nháp):** Trigger form submit và gán giá trị biến trạng thái gửi lên backend là "Ẩn" (Inactive).
- Nút **Lưu & Xuất bản (Primary):** Trigger form submit với biến trạng thái "Đang bán" (Active). Thành công điều hướng về list sản phẩm.

## 4. Component Breakdown
- `ProductCreatePage`: Component bọc màn hình chính, sử dụng thẻ `<Create>` từ `@refinedev/mui` và dùng hook `useForm`.
- `ProductBasicInfoCard`: Component chứa Tên và Mô tả sản phẩm.
- `ProductPricingCard`: Component chứa form Giá gốc và Giá khuyến mãi.
- `ProductInventoryCard`: Component chứa SKU và Số lượng.
- `ProductMediaCard`: Component chứa Uploader kéo thả ảnh.
- `ProductOrganizationCard`: Component chứa Dropdown Chọn Danh mục (dùng `useSelect`).
- `ProductStatusCard`: Component cấu hình Trạng thái hiển thị.

## 5. Data & State Management

### 5.1. Refine Hooks & Form
- Sử dụng `@refinedev/react-hook-form` với hook `useForm` (hoặc cấu trúc form có sẵn) tương tác với resource `products`.
- Sử dụng hook `useSelect` từ Refine để lấy danh sách danh mục phục vụ cho thẻ Category Dropdown.
- Tích hợp logic xử lý file/ảnh: Phụ thuộc vào config hệ thống, frontend upload ảnh lên cloud/S3 trước nhận URL, hoặc gửi file trong đối tượng FormData.

### 5.2. Form Validation (Sử dụng Yup / Zod / React Hook Form rules)
- **Tên sản phẩm:** `required` (Bắt buộc), max 255 ký tự.
- **Giá bán gốc:** `required`, `min: 1`.
- **Giá khuyến mãi:** Custom validate (nếu không trống, thì `value < gia_ban_goc`).
- **Danh mục:** `required`.
- **Số lượng tồn kho:** integer, `min: 0`.
- **Hình ảnh:** Báo lỗi nếu quá dung lượng (VD: > 5MB) hoặc sai định dạng.

## 6. Các tính năng nâng cao (Nghiệp vụ cốt lõi)
- **Real-time Validation:** Khai báo chế độ `mode: "onBlur"` hoặc `"onChange"` trong `react-hook-form` để hiển thị lỗi ngay khi user nhập sai.
- **Multi-Submit Action:** Ghi đè hàm callback submit. Tùy vào việc bấm nút "Lưu nháp" hay "Lưu & Xuất bản", mutate giá trị thuộc tính `status` trong dữ liệu gửi lên API trước khi gọi hàm create.
- **Unsaved Changes Warning:** Bật tính năng `warnWhenUnsavedChanges` của Refine để kích hoạt một popup cảnh báo nếu người dùng đang nhập dở nhưng lại back trình duyệt hoặc bấm link chuyển sang màn khác.
- **Drag & Drop Image:** Tái sử dụng component Upload hoặc sử dụng `react-dropzone` để thao tác chọn file dễ dàng.

## 7. Checklist Triển khai
- [ ] Khởi tạo page component `ProductCreatePage` với layout chia 2 cột Grid.
- [ ] Setup `useForm` (react-hook-form + Refine).
- [ ] Bố cục giao diện các Card (Basic Info, Pricing, Inventory, Media, etc.).
- [ ] Tích hợp Rich Text Editor Component cho trường "Mô tả".
- [ ] Sử dụng `useSelect` để call data render Danh mục Dropdown.
- [ ] Component Upload Hình ảnh hỗ trợ preview ảnh.
- [ ] Định nghĩa Validation Schema (Ràng buộc giá, tên, v.v.).
- [ ] Xử lý logic 2 nút lưu "Lưu nháp" và "Lưu xuất bản".
- [ ] Tích hợp popup cảnh báo thoát form khi đang nhập dở (`warnWhenUnsavedChanges`).
- [ ] Kiểm tra submit luồng tạo và hiển thị thông báo toast thành công/thất bại.
