---
name: code-ui-webview
description: Xây dựng giao diện public trong apps/webview bằng Next.js App Router, TypeScript và Tailwind CSS để hiển thị Page, Post và Section được quản trị từ apps/frontend. Dùng khi có lệnh /code-ui-webview hoặc yêu cầu code giao diện webview; không dùng cho giao diện admin Refine/MUI hay chỉ riêng việc nối API.
---

# Code UI Webview

Thi công giao diện public trong `apps/webview`. Áp dụng toàn bộ quy tắc nền tại `.agents/AGENTS.md` và quy tắc Next.js tại `apps/webview/AGENTS.md`; file này chỉ bổ sung workflow đặc thù cho content renderer.

## Phạm vi

- Chuyển HTML/Tailwind đã export trong `.docs/ui-mockups/` thành Next.js App Router.
- Xây page, layout và section renderer cho Page/Post/Section được quản trị từ admin.
- Chuẩn bị contract dữ liệu typed và fixture để UI có thể nối API mà không phải viết lại component.
- Không sửa `apps/frontend`, backend, API contract hoặc database nếu người dùng không yêu cầu.
- Nếu nhiệm vụ chủ yếu là thay mock bằng API thật, dùng `/integrate-webview`.

## Chuẩn bị

1. Xác định route, loại nội dung admin cần hiển thị, mockup và contract liên quan.
2. Đọc frontend plan tương ứng nếu có; không tự tạo giả định thay cho plan hoặc API contract đang tồn tại.
3. Đọc đúng tài liệu trong `apps/webview/node_modules/next/dist/docs/` trước khi dùng API hoặc convention Next.js có liên quan.
4. Quét `apps/webview/app`, `components`, `data`, `public` và `lib` nếu có. Tái sử dụng component và asset hiện tại trước khi tạo mới.

## Nguồn sự thật và thứ tự ưu tiên

Khi có khác biệt, ưu tiên theo thứ tự:

1. Yêu cầu hiện tại của người dùng.
2. API contract, backend plan và mô hình nội dung admin hiện có.
3. Mockup HTML/Tailwind và asset được export.
4. Frontend plan liên quan.
5. Pattern đã tồn tại trong `apps/webview`.

Có thể chuyển giá trị lặp lại từ bản export thành CSS variable/token dùng chung, nhưng không làm thay đổi kết quả hiển thị.

## Content renderer

- Dùng một registry/mapping tường minh từ `section.key` hoặc loại section ổn định sang component tương ứng; không suy luận component từ title hoặc nội dung hiển thị.
- Chỉ render Page/Post đã publish và Section đang active theo đúng `sortOrder`; dữ liệu ẩn/nháp phải được loại ở nguồn server và có kiểm tra phòng thủ tại boundary.
- Section không được hỗ trợ hoặc payload lỗi không được làm sập toàn trang: bỏ qua có chủ đích hoặc dùng fallback đã được thiết kế.
- Phân biệt rõ DTO từ API, model đã normalize và props của UI. Không import type phụ thuộc Refine/MUI từ `apps/frontend` vào webview.
- Nội dung rich text từ admin là dữ liệu không tin cậy. Không dùng `dangerouslySetInnerHTML` với HTML chưa được sanitize. Tái sử dụng sanitizer/renderer đã được phê duyệt trong repo; nếu chưa có, nêu rõ dependency hoặc contract còn thiếu thay vì render không an toàn.
- Giữ việc lấy và chuẩn hóa dữ liệu ở page, server-only module hoặc data layer. Component section chỉ nhận props typed.

## Next.js public page

- Dùng `next/image` cho ảnh nội dung khi phù hợp, khai báo kích thước hoặc `fill`/`sizes` đúng layout, và chỉ mở `remotePatterns` cho host ảnh thực sự cần.
- Dùng `next/link` cho điều hướng nội bộ; dùng thẻ `a` cho URL ngoài, email, số điện thoại và anchor cùng trang khi phù hợp.
- Sinh metadata từ SEO/Page/Post data bằng API metadata của App Router khi màn hình có route nội dung; có fallback hợp lý cho title, description, canonical và ảnh chia sẻ.
- Tận dụng `not-found`, error/loading boundary, cache và revalidation theo yêu cầu xuất bản nội dung; không hard-code chiến lược cache khi contract chưa xác định.
- Dùng HTML semantic, heading đúng cấp, landmark, label, alt text và keyboard interaction. Giữ responsive behavior theo mockup.
- Tránh client-side JavaScript cho nội dung tĩnh. Không hydrate cả section chỉ vì một control nhỏ cần tương tác.

## Dữ liệu trong giai đoạn code UI

- Không gọi API thật khi yêu cầu chỉ là code giao diện.
- Tạo fixture JSON có cấu trúc sát public API contract và đặt trong `apps/webview/data/`, trừ khi repo đã quy định vị trí khác.
- Khai báo type/model riêng trong webview và đưa dữ liệu qua cùng normalizer/props mà integration thật sẽ dùng.
- Dùng asset từ mockup hoặc `apps/webview/public`; không tự chọn ảnh ngẫu nhiên từ Internet.
- Thể hiện các trạng thái loading, empty, not found và error nếu chúng xuất hiện trong phạm vi màn hình hoặc plan.

Khi người dùng đồng thời yêu cầu nối API thật, kết hợp với `/integrate-webview` và tuân thủ data-fetching/debounce rules tại `.agents/AGENTS.md`.

## Trình tự thực hiện

1. Đối chiếu mockup với component/asset hiện có và lập mapping section → component.
2. Xác định ranh giới Server/Client cùng DTO → normalized model → UI props.
3. Cài đặt vào đúng thư mục trong `apps/webview`, giữ thay đổi ngoài phạm vi ở mức tối thiểu.
4. Chạy `npm run lint` trong `apps/webview`; chạy thêm `npm run build` khi thay đổi route, metadata, config, data fetching hoặc ranh giới Server/Client.
5. Báo cáo ngắn gọn file đã đổi, dữ liệu đang là fixture hay API thật, và kết quả kiểm tra.
