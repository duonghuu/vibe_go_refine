import React from "react";
import AppBar from "@mui/material/AppBar";
import Avatar from "@mui/material/Avatar";
import IconButton from "@mui/material/IconButton";
import Stack from "@mui/material/Stack";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import InputBase from "@mui/material/InputBase";
import Box from "@mui/material/Box";
import Badge from "@mui/material/Badge";
import MenuIcon from "@mui/icons-material/Menu";
import SearchIcon from "@mui/icons-material/Search";
import NotificationsNoneIcon from "@mui/icons-material/NotificationsNone";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { useGetIdentity } from "@refinedev/core";
import { useThemedLayoutContext } from "@refinedev/mui";

type IUser = {
  id: string;
  name: string;
  avatar: string;
  role: string;
};

export const CustomTopBar: React.FC = () => {
  const { data: user } = useGetIdentity<IUser>();
  const { setMobileSiderOpen, mobileSiderOpen } = useThemedLayoutContext();

  return (
    <AppBar position="sticky" elevation={0} sx={{ backgroundColor: 'background.paper', borderBottom: '1px solid', borderColor: 'divider', height: 70, justifyContent: 'center' }}>
      <Toolbar sx={{ px: { xs: 2, md: 4 } }}>
        <Stack direction="row" width="100%" justifyContent="space-between" alignItems="center">
          {/* Left Side: Hamburger & Search */}
          <Stack direction="row" alignItems="center" gap={3} flex={1}>
            <IconButton
              edge="start"
              color="inherit"
              aria-label="menu"
              onClick={() => setMobileSiderOpen(!mobileSiderOpen)}
              sx={{ display: { md: 'none' }, color: 'text.primary' }}
            >
              <MenuIcon />
            </IconButton>

            {/* Search */}
            <Box
              sx={{
                display: { xs: 'none', sm: 'flex' },
                alignItems: 'center',
                backgroundColor: '#F5F6FA',
                borderRadius: '50px',
                px: 2,
                py: 0.5,
                width: '388px',
                border: '1px solid #D5D5D5',
                '&:focus-within': {
                  borderColor: 'primary.main',
                }
              }}
            >
              <SearchIcon sx={{ color: 'text.secondary', mr: 1 }} />
              <InputBase
                placeholder="Search"
                sx={{ ml: 1, flex: 1, fontSize: '14px' }}
              />
            </Box>
          </Stack>

          {/* Right Side: Actions & Profile */}
          <Stack direction="row" alignItems="center" gap={3}>
            {/* Notification */}
            <IconButton sx={{ color: 'text.primary', opacity: 0.7, '&:hover': { opacity: 1 } }}>
              <Badge badgeContent={6} color="error" sx={{ '& .MuiBadge-badge': { fontSize: '10px', height: '16px', minWidth: '16px' } }}>
                <NotificationsNoneIcon />
              </Badge>
            </IconButton>

            {/* Language */}
            <Stack direction="row" alignItems="center" gap={1} sx={{ cursor: 'pointer', p: 1, borderRadius: 1, '&:hover': { bgcolor: 'action.hover' }, display: { xs: 'none', sm: 'flex' } }}>
               <img src="https://flagcdn.com/w40/gb.png" alt="English" style={{ width: 28, height: 20, borderRadius: 4, objectFit: 'cover' }} />
               <Typography variant="body2" sx={{ fontWeight: 600, color: '#646464' }}>English</Typography>
               <ExpandMoreIcon sx={{ fontSize: 16, color: 'text.secondary' }} />
            </Stack>

            {/* Profile */}
            <Stack direction="row" alignItems="center" gap={1.5} sx={{ cursor: 'pointer', p: 0.5, pr: 1.5, borderRadius: '50px', border: '1px solid transparent', '&:hover': { borderColor: 'divider' } }}>
              <Avatar src={user?.avatar || "https://ui-avatars.com/api/?name=Admin&background=4880FF&color=fff"} alt={user?.name} sx={{ width: 44, height: 44 }} />
              <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
                <Typography variant="subtitle2" sx={{ fontWeight: 'bold', color: '#404040', lineHeight: 1.2 }}>
                  {user?.name || "Admin User"}
                </Typography>
                <Typography variant="caption" sx={{ fontWeight: 600, color: '#565656' }}>
                  {user?.role || "Admin"}
                </Typography>
              </Box>
              <ExpandMoreIcon sx={{ fontSize: 16, color: 'text.secondary', ml: 0.5 }} />
            </Stack>
          </Stack>
        </Stack>
      </Toolbar>
    </AppBar>
  );
};
