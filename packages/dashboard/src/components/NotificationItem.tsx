import {
  MenuItem,
  Typography,
  styled,
  MenuItemProps,
  TypographyProps,
} from "@mui/material";

import { NotificationOnResponse } from "@/types/dtos/notification";
import { formatISO } from "@/lib/util/time";

interface NotificationItemProps {
  data: NotificationOnResponse;
  onClick: () => void;
}

interface NotificationItemElementProps extends MenuItemProps {
  seen: boolean;
}

function NotificationItem({ data, onClick }: NotificationItemProps) {
  return (
    <NotificationItemElement seen={data.seen} onClick={onClick}>
      <NotificationTitle>{data.title}</NotificationTitle>
      <Typography variant="body2">{formatISO(data.createdAt)}</Typography>
      <NotificationBody>{data.message}</NotificationBody>
    </NotificationItemElement>
  );
}

const NotificationItemElement = styled(
  ({ seen, ...props }: NotificationItemElementProps) => <MenuItem {...props} />,
)(({ theme, seen }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "start",
  whiteSpace: "normal",
  backgroundColor: seen ? "transparent" : theme.palette.grey[300],
  ":hover": {
    backgroundColor: seen ? theme.palette.grey[200] : theme.palette.grey[400],
  },
}));

const NotificationTitle = styled((props: TypographyProps) => (
  <Typography {...props} variant="subtitle1" />
))(() => ({
  fontWeight: "bold",
}));

const NotificationBody = styled((props: TypographyProps) => (
  <Typography {...props} variant="body1" />
))(({ theme }) => ({
  paddingTop: theme.typography.pxToRem(5),
  color: theme.palette.grey[700],
}));

export default NotificationItem;
