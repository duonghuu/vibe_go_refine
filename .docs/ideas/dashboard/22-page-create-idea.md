# Ý TƯỞNG: Tạo trang tĩnh (Admin Create Page)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Tạo mới Trang (Page).
* **Mục đích:** Cho phép Admin hoặc Staff tạo nội dung cho một trang tĩnh độc lập với bài viết.

---

## 2. Đối tượng & trải nghiệm

* **Người dùng chính:** Admin, Staff.
* **Hành động chính:** Nhập tiêu đề, slug, nội dung, trạng thái xuất bản; sau khi tạo thành công có thể gắn gallery và cấu hình SEO.
* **Trải nghiệm mong muốn:** Form hai cột, phần nội dung tập trung, các thiết lập xuất bản/media/SEO nằm ở cột phụ.

---

## 3. Đặc tả màn hình

### Header & actions

* Tiêu đề: **Tạo trang mới**.
* Breadcrumb: `Trang chủ > Quản trị nội dung > Trang > Tạo mới`.
* Nút **Hủy** quay lại `/pages`.
* Nút **Lưu nháp** đặt `status=DRAFT`.
* Nút **Xuất bản** đặt `status=PUBLISHED`.

### Cột nội dung

* **Tiêu đề (`title`):** bắt buộc, tối đa 255 ký tự.
* **Slug (`slug`):** bắt buộc; tự sinh từ tiêu đề cho đến khi người dùng chỉnh thủ công.
* **Nội dung (`content`):** Rich Text Editor, bắt buộc.

### Cột thiết lập

* **Trạng thái (`status`):** `DRAFT` hoặc `PUBLISHED`.
* **Hình ảnh:** dùng Page Media Card cho thumbnail và gallery sau khi Page đã được tạo.
* **SEO:** dùng SEO Meta Card chung với `entityType=page` sau khi Page đã được tạo.

* Khi chưa lưu được Page lần đầu, khu vực Media và SEO phải giải thích rằng cần lưu Page trước để có ID.
* Có Unsaved Changes Warning khi rời trang khi form đang thay đổi.

---

## 4. Validation & lỗi

| Trường | Quy tắc |
| --- | --- |
| `title` | Bắt buộc, sau khi trim không rỗng, tối đa 255 ký tự. |
| `slug` | Bắt buộc, unique toàn hệ thống, chỉ gồm chữ thường, số và dấu `-`. |
| `content` | Bắt buộc, không được rỗng sau khi loại bỏ nội dung trống của editor. |
| `status` | Chỉ nhận `DRAFT` hoặc `PUBLISHED`. |

* Lỗi slug trùng trả về theo field để form hiển thị đúng vị trí.
* `authorId` không hiển thị và không được gửi từ frontend; Backend lấy từ JWT.

---

## 5. API dự kiến

```text
POST /api/v1/admin/pages
```

```json
{
  "title": "Giới thiệu",
  "slug": "gioi-thieu",
  "content": "<p>Nội dung trang</p>",
  "status": "DRAFT"
}
```

Response thành công trả `201` cùng Page vừa tạo trong trường `data`.

---

## 6. Sau khi thao tác

* Thành công: toast phù hợp với trạng thái đã chọn; chuyển tới `/pages/edit/:id` để tiếp tục quản lý media và SEO.
* Thất bại: giữ lại dữ liệu form, hiển thị lỗi API và không tạo liên kết media/SEO mồ côi.

