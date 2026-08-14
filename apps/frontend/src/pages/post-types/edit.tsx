import { Edit } from "@refinedev/mui";
import { Box, Typography, Button, Stack } from "@mui/material";
import { useForm } from "@refinedev/react-hook-form";
import { PostTypeForm, IPostType, IPostTypeForm } from "./components/PostTypeForm";
import { HttpError, useNavigation } from "@refinedev/core";

export const PostTypeEdit = () => {
  const { list } = useNavigation();
  
  const form = useForm<IPostType, HttpError, IPostTypeForm>({
    refineCoreProps: {
      action: "edit",
      resource: "post-types",
      redirect: false,
      onMutationSuccess: () => {
        list("post-types");
      },
      meta: {
        method: "put",
      },
    },
  });
  const { saveButtonProps, refineCore: { formLoading } } = form;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Chỉnh sửa loại bài viết
        </Typography>
      </Box>
      <Edit
        title=""
        wrapperProps={{
          sx: {
            backgroundColor: "transparent",
            boxShadow: "none",
            p: 0,
          },
        }}
        headerProps={{
          sx: { display: "none" }
        }}
        footerButtons={
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: '100%' }}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={() => window.history.back()}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", color: 'text.secondary', borderColor: '#D5D5D5' }}
            >
              Hủy bỏ
            </Button>
            <Button
              {...saveButtonProps}
              variant="contained"
              color="success"
              disabled={formLoading}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'success.main', boxShadow: 'none', '&:hover': { bgcolor: 'success.dark', boxShadow: 'none' } }}
            >
              Lưu thay đổi
            </Button>
          </Stack>
        }
      >
        <PostTypeForm form={form} isEdit={true} />
      </Edit>
    </Box>
  );
};
