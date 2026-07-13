import React from "react";
import { CustomSidebar } from "./sidebar";
import { CustomTopBar } from "./header";
import Box from "@mui/material/Box";
import { ThemedLayout } from "@refinedev/mui";

export interface AdminMasterLayoutProps {
  children: React.ReactNode;
}

export const AdminMasterLayout: React.FC<AdminMasterLayoutProps> = ({ children }) => {
  return (
    <ThemedLayout
      Header={CustomTopBar}
      Sider={CustomSidebar}
      Title={() => null} // Header is handled in Sidebar
    >
      <Box sx={{ 
        p: { xs: 2, md: 4 }, 
        bgcolor: '#F5F6FA', 
        minHeight: '100vh',
        width: '100%',
      }}>
        {children}
      </Box>
    </ThemedLayout>
  );
};
