import React from "react";
import { Card, CardContent, Stack, Typography } from "@mui/material";
import ImageOutlinedIcon from "@mui/icons-material/ImageOutlined";
import SearchOutlinedIcon from "@mui/icons-material/SearchOutlined";

interface PagePostCreateNoticeProps {
  title: string;
  description: string;
  icon: "media" | "seo";
}

export const PagePostCreateNotice: React.FC<PagePostCreateNoticeProps> = ({ title, description, icon }) => (
  <Card sx={{ border: "1px solid", borderColor: "divider", borderRadius: 2, boxShadow: "none", opacity: 0.8 }}>
    <CardContent sx={{ p: { xs: 2, md: 3 } }}>
      <Stack direction="row" spacing={1.5} alignItems="center" mb={1.5}>
        {icon === "media" ? <ImageOutlinedIcon color="action" /> : <SearchOutlinedIcon color="action" />}
        <Typography variant="h6" fontWeight={700}>{title}</Typography>
      </Stack>
      <Typography variant="body2" color="text.secondary">{description}</Typography>
    </CardContent>
  </Card>
);
