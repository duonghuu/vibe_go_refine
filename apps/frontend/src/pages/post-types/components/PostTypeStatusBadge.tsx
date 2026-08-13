import { Chip } from "@mui/material";

interface PostTypeStatusBadgeProps {
  status: "ACTIVE" | "INACTIVE";
}

export const PostTypeStatusBadge = ({ status }: PostTypeStatusBadgeProps) => {
  const isActive = status === "ACTIVE";
  return (
    <Chip
      label={isActive ? "ACTIVE" : "INACTIVE"}
      size="small"
      sx={{
        backgroundColor: isActive ? "#e8f5e9" : "#ffebee",
        color: isActive ? "#2e7d32" : "#c62828",
        fontWeight: "bold",
      }}
    />
  );
};
