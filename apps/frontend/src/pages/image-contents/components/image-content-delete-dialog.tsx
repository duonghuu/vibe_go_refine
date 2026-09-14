import { DeleteButton } from "@refinedev/mui";

export const ImageContentDeleteButton = ({ id }: { id: number }) => <DeleteButton hideText recordItemId={id} confirmTitle="Xóa nội dung hình ảnh?" confirmOkText="Xóa" confirmCancelText="Hủy" />;

