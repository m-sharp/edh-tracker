import { ReactElement } from "react";
import { Outlet, Link, useLocation } from "react-router-dom";
import {
    AppBar,
    Box,
    Button,
    Container,
    SvgIcon,
    Toolbar,
    Typography
} from "@mui/material";

export default function Root(): ReactElement {
    const location = useLocation();

    return (
        <Box sx={{ display: "flex", width: "auto" }}>
            <DrawerAppBar />
            <Container id="detail" component="main" sx={{ p: 3, width: "90%", bgcolor: "#f0f5fa", mt: 12, mb: 5 }} maxWidth="xl">
                <RootContent path={location.pathname} />
                <Outlet />
            </Container>
        </Box>
    );
}

// ToDo: Add mobile menu icon and link menu
// ToDo: Link button click should propagate down to <Link>
function DrawerAppBar(): ReactElement {
    return (
        <AppBar position="fixed">
            <Container maxWidth="xl">
                <Toolbar disableGutters>
                    <SvgIconPlayingCards />
                    <Typography
                        variant="h6"
                        noWrap
                        sx={{
                            mr: 2,
                            display: "flex",
                            fontFamily: "monospace",
                            fontWeight: 700,
                            letterSpacing: ".3rem",
                        }}
                    >
                        <Link to={`/`} style={{textDecoration: "none", color: "white"}}>EDH Tracker</Link>
                    </Typography>
                    <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }}>
                        <Button sx={{ my: 2, color: 'white', display: 'block' }}>
                            <Link to={`/players`} style={{color: "white"}}>Players</Link>
                        </Button>
                        <Button sx={{ my: 2, color: 'white', display: 'block' }}>
                            <Link to={`/decks`} style={{color: "white"}}>Decks</Link>
                        </Button>
                        <Button sx={{ my: 2, color: 'white', display: 'block' }}>
                            <Link to={`/games`} style={{color: "white"}}>Games</Link>
                        </Button>
                    </Box>
                </Toolbar>
            </Container>
        </AppBar>
    );
}

function SvgIconPlayingCards(): ReactElement {
    // Via https://fonts.google.com/icons?selected=Material%20Symbols%20Outlined%3Aplaying_cards%3AFILL%400%3Bwght%40400%3BGRAD%400%3Bopsz%4024
    return (
        <SvgIcon sx={{ display: "flex", mr: 2 }}>
            <svg
                xmlns="http://www.w3.org/2000/svg"
                height="24"
                viewBox="0 -960 960 960"
                width="24"
                strokeWidth={1.5}
                stroke="currentColor"
                fill="white"
            >
                <path
                    d="m608-368 46-166-142-98-46 166 142 98ZM160-207l-33-16q-31-13-42-44.5t3-62.5l72-156v279Zm160 87q-33 0-56.5-24T240-201v-239l107 294q3 7 5 13.5t7 12.5h-39Zm206-5q-31 11-62-3t-42-45L245-662q-11-31 3-61.5t45-41.5l301-110q31-11 61.5 3t41.5 45l178 489q11 31-3 61.5T827-235L526-125Zm-28-75 302-110-179-490-301 110 178 490Zm62-300Z"
                />
            </svg>
        </SvgIcon>
    );
}

interface RootContentProps {
    path: string;
}

function RootContent({path}: RootContentProps): ReactElement {
    if (path !== "/") {
        return (<></>);
    }

    return (
        <Box id="rootContent" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <h2>Welcome to the EDH Tracker!</h2>
            <p>We'll have some login stuff here eventually...</p>
        </Box>
    )
}
