import { ReactElement } from "react";
import { IconButton, Tooltip } from "@mui/material";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

interface TooltipIconProps {
    title: string;
    icon?: ReactElement;
}

export function TooltipIcon({ title, icon }: TooltipIconProps): ReactElement {
    return (
        <Tooltip title={title} enterTouchDelay={0}>
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
}

export function TooltipIconButton({ title, onClick, icon }: TooltipIconButtonProps): ReactElement {
    return (
        <Tooltip title={title}>
            <IconButton size="medium" onClick={onClick}>
                {icon}
            </IconButton>
        </Tooltip>
    );
}
