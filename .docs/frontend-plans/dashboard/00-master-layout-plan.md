# KẾ HOẠCH KỸ THUẬT: BỐ CỤC TOÀN CỤC ADMIN DASHBOARD (MASTER LAYOUT)

## TỐI ƯU HÓA CHO REFINE.JS + MATERIAL UI (MUI)

### 1. PHÂN RÃ COMPONENT (COMPONENT TREE)

* **AdminMasterLayout** `[SMART]` *(Custom từ `<ThemedLayoutV2>` của Refine)*
* *Vai trò:* Nhận `children` làm nội dung chính. Tích hợp sẵn context để đồng bộ trạng thái đóng/mở Sidebar trên mobile và desktop qua `<ThemedLayoutContextProvider>`.
* **CustomSidebar** `[SMART]` *(Custom từ `<ThemedSiderV2>`)*
* *Vai trò:* Cố định bên trái. Tự động lắng nghe cấu hình `resources` truyền vào `<Refine>` để sinh ra menu tương ứng.
* **Logo** `[DUMB]` *(Custom từ `<ThemedTitleV2>`)*
* *Vai trò:* Hiển thị logo "DashStack", click sẽ điều hướng về trang chủ cấu hình trong `routerProvider` (thường là trang dashboard).


* **NavigationMenu** `[SMART]`
* *Vai trò:* Sử dụng hook `useMenu` của Refine để lấy danh sách các resource đã phân quyền kèm theo trạng thái active hiện tại.
* **NavItem** `[DUMB]`
* *Vai trò:* Sử dụng `<CanAccess>` component để kiểm tra quyền hiển thị menu + `<ActiveLink>` (hoặc `NavLink` từ router) để điều hướng.


* **NavGroupTitle** `[DUMB]`
* *Vai trò:* Render tiêu đề nhóm dựa trên thuộc tính `meta.parent` hoặc `meta.label` trong định nghĩa resource của Refine.






* **MainContentWrapper** `[DUMB]`
* **CustomTopBar (Header)** `[SMART]` *(Custom từ `<ThemedHeaderV2>`)*
* *Vai trò:* Header cố định ở trên cùng, chứa các action và tích hợp đồng bộ dữ liệu người dùng từ `authProvider`.
* **SidebarToggle** `[DUMB]` *(Sử dụng nút toggle từ hook `useThemedLayoutContext`)*
* **GlobalSearch** `[SMART]`
* *Vai trò:* Ô nhập từ khóa. Có thể tích hợp với tính năng tìm kiếm của Refine hoặc đẩy lên URL qua router query.


* **TopBarActions** `[SMART]`
* **NotificationBell** `[DUMB]`
* **LanguageSelector** `[SMART]` *(Tích hợp hook `useGetLocale` và `useSetLocale` của Refine i18n)*
* **AdminProfile** `[SMART]`
* *Vai trò:* Sử dụng hook `useGetIdentity` để tự động lấy `avatar`, `name`, `role` từ `authProvider`. Chứa menu Logout gọi hook `useLogout`.






* **PageContent** `[DUMB]`
* *Vai trò:* Bọc nội dung trang con, tự động render `<Breadcrumb />` của Refine ở đầu trang nếu cần, padding `p-8`.







---

### 2. QUẢN LÝ TRẠNG THÁI (STATE MANAGEMENT)

* **Trạng thái Admin (`identity`, `isAuthenticated`)**:
* *Chiến lược:* **Không dùng Zustand/Context API thủ công**. Refine quản lý tập trung thông qua `authProvider`.
* Sử dụng hook `useGetIdentity<IUser>()` tại `AdminProfile` để lấy thông tin user real-time. Guard định tuyến được xử lý tự động bởi component `<Authenticated>` bọc ngoài hệ thống routes.


* **Trạng thái Active của Menu (`selectedKey`)**:
* *Chiến lược:* Sử dụng hook **`useMenu()`** của Refine. Hook này tự động phân tích URL hiện tại, đối chiếu với danh sách `resources` để trả về thuộc tính `selectedKey` (trang nào đang active) và `openKeys` (nhóm menu nào đang mở rộng).


* **Trạng thái đóng/mở Sidebar (`siderCollapsed`)**:
* *Chiến lược:* Sử dụng **`useThemedLayoutContext()`** từ `@refinedev/mui` (hoặc core). Trạng thái này được quản lý toàn cục trong layout provider giúp đồng bộ nút bấm ở Header và hành vi thu phóng của Sidebar mà không cần truyền props thủ công.


* **Trạng thái Đa ngôn ngữ (`locale`)**:
* *Chiến lược:* Sử dụng **`i18nProvider`** của Refine kết hợp với các hooks `useTranslation`, `useSetLocale`.



---

### 3. CẤU TRÚC DỮ LIỆU (DATA INTERFACES)

Khi dùng Refine, ta tận dụng các kiểu dữ liệu từ hệ sinh thái của nó thay vì tự định nghĩa từ đầu:

```typescript
import React from "react";
import { ITreeMenu } from "@refinedev/core";

// --- Custom Sidebar Menu Props ---
// Thay vì truyền thủ công từng item, ta map qua mảng `menuItems` trả về từ useMenu()
export interface NavigationMenuProps {
  menuItems: ITreeMenu[]; // Cấu trúc menu phân cấp (gồm cả con, cha, icon) của Refine
  selectedKey: string;
}

export interface NavItemProps {
  item: ITreeMenu; // Chứa sẵn key, route, name, icon, và children
  isActive: boolean;
}

// --- Auth & Identity Interfaces (Dữ liệu từ authProvider.getIdentity) ---
export interface IAdminIdentity {
  id: string;
  name: string;
  avatar: string;
  role: string; // Vd: "Admin", "Manager"
  email?: string;
}

export interface AdminProfileProps {
  identity: IAdminIdentity | undefined;
  onLogout: () => void; // Trigger hàm logout() từ useLogout
}

// --- Layout Component Props ---
export interface AdminMasterLayoutProps {
  children: React.ReactNode;
}

```