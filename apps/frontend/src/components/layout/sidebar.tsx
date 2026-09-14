import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Drawer from "@mui/material/Drawer";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Typography from "@mui/material/Typography";
import { useGetIdentity, useLogout, useMenu } from "@refinedev/core";
import { useThemedLayoutContext } from "@refinedev/mui";
import React from "react";
import { NavLink } from "react-router";
import LogoutIcon from "@mui/icons-material/Logout";
import SettingsIcon from "@mui/icons-material/Settings";

const drawerWidth = 240;

export const CustomSidebar: React.FC = () => {
  const { menuItems, selectedKey } = useMenu();
  const { data: identity } = useGetIdentity<{ role?: string }>();
  const { siderCollapsed, mobileSiderOpen, setMobileSiderOpen } =
    useThemedLayoutContext();
  const { mutate: logout } = useLogout();

  const handleDrawerToggle = () => {
    setMobileSiderOpen(!mobileSiderOpen);
  };

  const drawerContent = (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%", bgcolor: 'background.default' }}>
      {/* Logo */}
      <Box
        sx={{
          height: 70,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          borderBottom: "1px solid",
          borderColor: "divider",
          flexShrink: 0,
        }}
      >
        <Typography
          variant="h6"
          sx={{
            fontWeight: 800,
            fontSize: "20px",
            letterSpacing: "-0.025em",
          }}
        >
          <Box component="span" sx={{ color: "primary.main" }}>
            Dash
          </Box>
          <Box component="span" sx={{ color: "text.primary" }}>
            Stack
          </Box>
        </Typography>
      </Box>

      {/* Navigation Menu */}
      <Box
        sx={{
          flexGrow: 1,
          overflowY: "auto",
          px: 2,
          py: 2,
          "&::-webkit-scrollbar": {
            width: "4px",
          },
          "&::-webkit-scrollbar-thumb": {
            backgroundColor: "#E0E0E0",
            borderRadius: "4px",
          },
        }}
      >
        <List sx={{ pt: 0 }}>
          {menuItems.filter((item) => {
            const meta = item.meta as { requiredRoles?: string[] } | undefined;
            return !meta?.requiredRoles || (identity?.role ? meta.requiredRoles.includes(identity.role) : false);
          }).map((item) => {
            const isSelected = item.key === selectedKey;

            return (
              <ListItemButton
                key={item.key}
                component={NavLink}
                to={item.route || "/"}
                sx={{
                  borderRadius: "12px",
                  mb: 1,
                  backgroundColor: isSelected ? "primary.main" : "transparent",
                  color: isSelected ? "white" : "text.primary",
                  "&:hover": {
                    backgroundColor: isSelected
                      ? "primary.dark"
                      : "action.hover",
                  },
                }}
              >
                <ListItemIcon
                  sx={{
                    minWidth: "40px",
                    color: isSelected ? "white" : "action.active",
                    opacity: isSelected ? 1 : 0.7,
                  }}
                >
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.label}
                  primaryTypographyProps={{
                    fontSize: "14px",
                    fontWeight: 600,
                    sx: { opacity: isSelected ? 1 : 0.7 },
                  }}
                />
              </ListItemButton>
            );
          })}
        </List>

        <Box sx={{ pt: 3, pb: 1, px: 2 }}>
           <Typography variant="caption" sx={{ fontWeight: 'bold', color: 'text.primary', opacity: 0.6, letterSpacing: '0.05em' }}>
              PAGES
           </Typography>
        </Box>

        <Divider sx={{ my: 2 }} />

        {/* Bottom Actions */}
        <List>
          <ListItemButton
            sx={{
              borderRadius: "12px",
              mb: 1,
              "&:hover": { backgroundColor: "action.hover" },
            }}
          >
            <ListItemIcon sx={{ minWidth: "40px", opacity: 0.7 }}>
              <SettingsIcon />
            </ListItemIcon>
            <ListItemText
              primary="Settings"
              primaryTypographyProps={{
                fontSize: "14px",
                fontWeight: 600,
                sx: { opacity: 0.7 },
              }}
            />
          </ListItemButton>
          <ListItemButton
            onClick={() => logout()}
            sx={{
              borderRadius: "12px",
              "&:hover": { backgroundColor: "action.hover" },
            }}
          >
            <ListItemIcon sx={{ minWidth: "40px", opacity: 0.7 }}>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText
              primary="Logout"
              primaryTypographyProps={{
                fontSize: "14px",
                fontWeight: 600,
                sx: { opacity: 0.7 },
              }}
            />
          </ListItemButton>
        </List>
      </Box>
    </Box>
  );

  return (
    <Box
      component="nav"
      sx={{
        width: { md: siderCollapsed ? 0 : drawerWidth },
        flexShrink: { md: 0 },
        transition: "width 0.2s",
      }}
    >
      <Drawer
        variant="temporary"
        open={mobileSiderOpen}
        onClose={handleDrawerToggle}
        ModalProps={{
          keepMounted: true,
        }}
        sx={{
          display: { xs: "block", md: "none" },
          "& .MuiDrawer-paper": {
            boxSizing: "border-box",
            width: drawerWidth,
            borderRight: "1px solid",
            borderColor: "divider",
          },
        }}
      >
        {drawerContent}
      </Drawer>
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: "none", md: "block" },
          "& .MuiDrawer-paper": {
            boxSizing: "border-box",
            width: siderCollapsed ? 0 : drawerWidth,
            overflowX: "hidden",
            transition: "width 0.2s",
            borderRight: "1px solid",
            borderColor: "divider",
          },
        }}
        open
      >
        {drawerContent}
      </Drawer>
    </Box>
  );
};
