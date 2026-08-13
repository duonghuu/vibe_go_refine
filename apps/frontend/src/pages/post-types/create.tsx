import { Create } from "@refinedev/mui";
import { Box, Typography, Button, Stack } from "@mui/material";
import { useForm } from "@refinedev/react-hook-form";
import { PostTypeForm } from "./components/PostTypeForm";

export const PostTypeCreate = () => {
  const form = useForm({
    defaultValues: {
      status: "ACTIVE",
      sort_order: 1
    }
  });

  const { saveButtonProps, refineCore: { formLoading } } = form;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Thêm loại bài viết mới
        </Typography>
      </Box>
      <Create
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
              Hủy
            </Button>
            <Button
              {...saveButtonProps}
              variant="contained"
              disabled={formLoading}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'primary.main', boxShadow: 'none' }}
            >
              Lưu & Thêm mới
            </Button>
          </Stack>
        }
      >
        <PostTypeForm form={form} />
      </Create>
    </Box>
  );
};
