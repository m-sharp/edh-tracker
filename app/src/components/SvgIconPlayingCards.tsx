import { ReactElement } from "react";
import { SvgIcon } from "@mui/material";

export default function SvgIconPlayingCards({ fontSize }: { fontSize?: string | number }): ReactElement {
    // Via https://fonts.google.com/icons?selected=Material%20Symbols%20Outlined%3Aplaying_cards%3AFILL%400%3Bwght%40400%3BGRAD%400%3Bopsz%4024
    return (
        <SvgIcon sx={{ display: "inline-flex", mr: 0, fontSize }}>
            <svg
                xmlns="http://www.w3.org/2000/svg"
                height="24"
                viewBox="0 -960 960 960"
                width="24"
                strokeWidth={1.5}
                stroke="currentColor"
                fill="white"
            >
                <path d="m608-368 46-166-142-98-46 166 142 98ZM160-207l-33-16q-31-13-42-44.5t3-62.5l72-156v279Zm160 87q-33 0-56.5-24T240-201v-239l107 294q3 7 5 13.5t7 12.5h-39Zm206-5q-31 11-62-3t-42-45L245-662q-11-31 3-61.5t45-41.5l301-110q31-11 61.5 3t41.5 45l178 489q11 31-3 61.5T827-235L526-125Zm-28-75 302-110-179-490-301 110 178 490Zm62-300Z" />
            </svg>
        </SvgIcon>
    );
}
