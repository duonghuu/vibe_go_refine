# Ý TƯỞNG: Public API cho trang tĩnh (Public Page API)

## 1. Thông tin chung

* **Dự án:** TechBite.
* **Tính năng:** Public API đọc Page theo slug.
* **Mục đích:** Cung cấp contract ổn định để Webview/Next.js render các trang tĩnh trong giai đoạn sau.
* **Ngoài phạm vi:** Hạng mục này không dựng dynamic route hoặc giao diện Page trên Webview.

---

## 2. Hành vi API

```text
GET /api/v1/pages/:slug
```

* Không yêu cầu Access Token.
* Chỉ trả bản ghi `PUBLISHED` chưa bị soft-delete.
* `DRAFT`, Page đã xóa hoặc slug không tồn tại đều trả `404`, không làm lộ sự tồn tại của nội dung nháp.
* Response bao gồm nội dung, thumbnail/gallery đã sắp xếp và SEO đã resolve.

Ví dụ response:

```json
{
  "data": {
    "id": 1,
    "title": "Giới thiệu",
    "slug": "gioi-thieu",
    "content": "<p>Nội dung trang</p>",
    "updatedAt": "2026-08-31T10:30:00Z",
    "thumbnail": null,
    "gallery": [],
    "seo": {
      "resolvedTitle": "Giới thiệu",
      "resolvedDescription": "Nội dung trang",
      "resolvedCanonicalUrl": "https://example.com/gioi-thieu",
      "robots": "index,follow"
    }
  }
}
```

---

## 3. SEO & cache

* Dùng module `seo_meta` hiện có với `entityType=page`; Entity Registry phải xác thực Page tồn tại và chưa bị xóa.
* Fallback SEO lấy từ title/content/thumbnail của Page, không ghi giá trị fallback ngược vào database.
* Cache response public theo slug trong Redis.
* Tạo, sửa, xóa Page; thay gallery/thumbnail; hoặc upsert/xóa SEO đều phải xóa cache slug tương ứng.

---

## 4. Xử lý lỗi

* `400 Bad Request`: slug không đúng định dạng.
* `404 Not Found`: không có Page công khai tương ứng.
* `500 Internal Server Error`: lỗi dữ liệu hoặc hạ tầng; không trả nội dung nội bộ cho client.

