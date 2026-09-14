import { DeleteButton } from "@refinedev/mui";

export const ImageContentTypeDeleteButton = ({ id }: { id: number }) => <DeleteButton hideText recordItemId={id} confirmTitle="Xóa loại nội dung?" confirmOkText="Xóa" confirmCancelText="Hủy" />;

