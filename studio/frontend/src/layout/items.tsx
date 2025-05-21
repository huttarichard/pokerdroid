import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import {
  DashboardOutlined,
  LocalLibraryOutlined,
  SportsEsportsOutlined,
} from "@mui/icons-material";
import { Fragment } from "react/jsx-runtime";

export const mainListItems = (
  <Fragment>
    <ListItemButton href="/">
      <ListItemIcon>
        <DashboardOutlined />
      </ListItemIcon>
    </ListItemButton>

    <ListItemButton href="/tree">
      <ListItemIcon>
        <LocalLibraryOutlined />
      </ListItemIcon>
    </ListItemButton>

    <ListItemButton href="/game">
      <ListItemIcon>
        <SportsEsportsOutlined />
      </ListItemIcon>
    </ListItemButton>
  </Fragment>
);
