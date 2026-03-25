import { ReactElement } from "react";
import { IconButton, SxProps, Theme, Tooltip } from "@mui/material";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

interface TooltipIconProps {
    title: string;
    icon?: ReactElement;
    placement?: "top" | "bottom" | "left" | "right";
}

export function TooltipIcon({ title, icon, placement = "top" }: TooltipIconProps): ReactElement {
    return (
        <Tooltip title={title} enterTouchDelay={0} placement={placement}>
            <span style={{ display: "inline-flex", alignItems: "center", cursor: "default" }}>
                {icon ?? <InfoOutlinedIcon fontSize="small" color="action" />}
            </span>
        </Tooltip>
    );
}

interface TooltipIconButtonProps {
    title: string;
    onClick: () => void;
    icon: ReactElement;
    placement?: "top" | "bottom" | "left" | "right";
    color?: "default" | "inherit" | "primary" | "secondary" | "error" | "info" | "success" | "warning";
    size?: "small" | "medium" | "large";
    disabled?: boolean;
    sx?: SxProps<Theme>;
}

export function TooltipIconButton({
    title,
    onClick,
    icon,
    placement = "top",
    color,
    size = "medium",
    disabled,
    sx,
}: TooltipIconButtonProps): ReactElement {
    return (
        <Tooltip title={title} placement={placement}>
            <span>
                <IconButton size={size} onClick={onClick} color={color} disabled={disabled} sx={sx}>
                    {icon}
                </IconButton>
            </span>
        </Tooltip>
    );
}
