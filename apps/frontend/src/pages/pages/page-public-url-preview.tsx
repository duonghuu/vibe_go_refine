import React from "react";
import { Card, CardContent, Link, Stack, Typography } from "@mui/material";
import OpenInNewOutlinedIcon from "@mui/icons-material/OpenInNewOutlined";

interface PagePublicUrlPreviewProps {
  slug: string;
}

export const PagePublicUrlPreview: React.FC<PagePublicUrlPreviewProps> = ({ slug }) => (
  <Card sx={{ border: "1px solid", borderColor: "divider", borderRadius: 2, boxShadow: "none", mb: 3 }}>
    <CardContent sx={{ p: { xs: 2, md: 3 } }}>
      <Typography variant="h6" fontWeight={700} mb={1.5}>URL public</Typography>
      <Stack direction="row" spacing={1} alignItems="center">
        <Link href={slug ? `/${slug}` : undefined} underline="hover" color="primary.main" sx={{ overflowWrap: "anywhere", fontSize: 14 }}>
          /{slug || "slug-trang"}
        </Link>
        <OpenInNewOutlinedIcon fontSize="small" color="action" />
      </Stack>
      <Typography variant="caption" color="text.secondary" display="block" mt={1}>
        Preview theo slug hiện tại, không tạo redirect tự động.
      </Typography>
    </CardContent>
  </Card>
);
