import { ReactElement } from "react";
import { IconButton, Tooltip } from "@mui/material";
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
}

export function TooltipIconButton({ title, onClick, icon, placement = "top" }: TooltipIconButtonProps): ReactElement {
    return (
        <Tooltip title={title} placement={placement}>
            <IconButton size="medium" onClick={onClick}>
                {icon}
            </IconButton>
        </Tooltip>
    );
}
