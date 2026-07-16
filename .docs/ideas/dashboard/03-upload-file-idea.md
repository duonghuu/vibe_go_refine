# Ý TƯỞNG: Chức năng upload file

## 1. UI/UX — Góc nhìn UI/UX Designer

**Về khu vực upload:**

**UI/UX Designer:** Với sản phẩm, tôi khuyên kết hợp cả 3 kiểu chứ không chọn 1:
- **Drag & drop** làm vùng chính (khung lớn, dashed border, đổi màu khi hover file vào)
- **Click to upload** làm fallback (vẫn phải có, vì không phải ai cũng biết kéo thả, đặc biệt trên mobile)

**Trạng thái chờ & preview:**

**UI/UX Designer:** 
- Progress bar **theo từng file** là bắt buộc nếu cho phép multi-upload, không nên gộp chung 1 thanh progress cho cả batch — người dùng sẽ không biết file nào bị lỗi.
- Preview ảnh nên hiện ngay lập tức bằng `URL.createObjectURL()` ở client, **trước khi** file kịp upload lên server — cho cảm giác phản hồi tức thì (perceived performance).
- Crop/xoay ảnh trực tiếp trên client (dùng thư viện như `react-easy-crop` hoặc `cropperjs`) là best practice của các sàn TMĐT lớn, tránh phải upload rồi mới sửa lại (tốn băng thông 2 lần).
- Nên có trạng thái rõ ràng theo từng ảnh: `pending` → `uploading (%)` → `success` / `error` (kèm nút retry riêng cho từng ảnh lỗi, không bắt user upload lại cả batch).

**Phân biệt ảnh cũ/mới khi Edit:**

**UI/UX Designer:** Đây là bài toán UX dễ gây rối nhất. Gợi ý dùng **visual convention thống nhất**:
- Ảnh cũ giữ nguyên: hiển thị bình thường, có icon "X" nhỏ góc để xóa
- Ảnh cũ đã đánh dấu xóa: làm mờ (opacity 40%) + gạch chéo, kèm nút "Hoàn tác" (undo) — **không xóa khỏi UI ngay**, chỉ đánh dấu, để user có thể đổi ý
- Ảnh mới vừa chọn: viền màu nổi bật (VD viền xanh) + badge "Mới" nhỏ ở góc

Về mặt data, tôi đề xuất Solution Architect thiết kế state ở client dạng 3 mảng riêng: `existingImages[]`, `imagesToDelete[]` (chỉ chứa ID), `newImages[]` (File objects) — khi submit gộp lại thành 1 payload.

---

## 2. Kiến trúc & Lưu trữ — Solution Architect & Product Owner

**Solution Architect:** Về nơi lưu trữ, xếp hạng theo mức độ phổ biến trong dự án thực tế:

- Giải pháp: Lưu ổ đĩa server
- Ưu điểm: Đơn giản, không cần setup gì
- Nhược điểm: Không scale khi có nhiều server (load balancer), mất dữ liệu nếu server chết, backup thủ công

**Về resize/nén — nên làm ở đâu:**

**Solution Architect:** Tuyệt đối nên làm ở **Backend** (hoặc dùng dịch vụ như Cloudinary tự động làm), không phải Frontend. Lý do:
- Frontend resize chỉ tối ưu **băng thông upload**, nhưng nếu không kiểm soát ở backend, ai đó có thể bypass frontend (gọi API trực tiếp) và đẩy file gốc siêu nặng lên
- Cần backend luôn tạo lại 3-4 size chuẩn: `thumbnail` (~150px, cho list/grid), `medium` (~600px, cho trang chi tiết), `original` (lưu để archive, có thể không serve trực tiếp)

**Product Owner:** Đồng ý với queue-based resize. Về chi phí, tôi muốn nhấn thêm: luôn set giới hạn kích thước file tối đa (VD 5MB/ảnh) và giới hạn số lượng ảnh/sản phẩm (VD tối đa 10 ảnh) ngay từ đầu — không giới hạn sẽ khiến storage cost khó dự đoán khi scale.

---

## 3. Bảo mật — Security Engineer

**Security Engineer:** Đây là phần tôi lo nhất vì upload file là một trong những vector tấn công phổ biến nhất. Checklist bắt buộc:

**Chặn file độc hại:**
- **Không bao giờ tin đuôi file hay MIME type do client gửi lên.** Client có thể đổi tên `shell.php` thành `shell.php.jpg` hoặc set header `Content-Type: image/jpeg` giả.
- Bắt buộc kiểm tra **magic bytes** (file signature thực sự) ở backend — dùng thư viện như `file-type` (Node) hoặc `python-magic` để đọc vài byte đầu và xác nhận đúng là JPEG/PNG/WebP thật.
- Nếu server nằm trong hạ tầng cho phép, dùng thư viện ảnh (Sharp, Pillow) để **re-encode lại ảnh** (decode rồi encode lại thành file mới) — cách này triệt để loại bỏ payload ẩn trong metadata hoặc polyglot file (file vừa là ảnh vừa là script hợp lệ).
- **Không bao giờ** lưu file upload trong thư mục có quyền thực thi (executable), và **không bao giờ** đặt tên file theo tên gốc do user gửi — luôn generate tên mới (UUID) để tránh path traversal (`../../etc/passwd`) và tránh overwrite file người khác.

**Chống spam upload:**
- Rate limiting theo user/IP cho endpoint upload (VD tối đa 20 request/phút)
- Giới hạn kích thước request ở tầng reverse proxy (Nginx `client_max_body_size`) trước khi request chạm tới code backend — chặn từ sớm, đỡ tốn tài nguyên xử lý
- Validate số lượng file và tổng dung lượng trước khi bắt đầu nhận file (không đợi upload xong mới báo lỗi)

---

## 4. Instant Upload

---

## 5. Giải quyết File mồ côi — Cả hội đồng

**Solution Architect:** Cơ chế chuẩn mà các dự án lớn dùng, theo tôi thấy có 3 tầng phòng thủ:

**Tầng 1 — Đánh dấu trạng thái ngay khi upload:**
Khi file upload thành công (theo hướng Instant Upload), backend lưu 1 record vào bảng `media` hoặc `files` với trạng thái `status = 'temporary'` và `owner_id = user_id`, `product_id = NULL`. Khi user bấm Save Product thành công, backend update các file được chọn sang `status = 'attached'` và gắn `product_id`.

**Tầng 2 — Background job dọn rác định kỳ:**
Cron job (chạy mỗi giờ hoặc mỗi ngày tùy volume) quét bảng `media` where `status = 'temporary' AND created_at < now() - interval '24 hours'` → xóa file khỏi storage + xóa record DB. Đây là cơ chế phổ biến nhất, đơn giản và đủ tin cậy.

