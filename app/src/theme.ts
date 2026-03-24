import { createTheme } from "@mui/material/styles";

const theme = createTheme({
    palette: {
        mode: "dark",
        background: {
            default: "#0f1117",
            paper: "#1a1a2e",
        },
        primary: {
            main: "#c9a227",
            light: "#e0b83a",
        },
        text: {
            primary: "#e8edf5",
            // secondary: MUI dark-mode default opacity applied to text.primary automatically
        },
        // error: MUI default (no override — D-06)
    },
    typography: {
        fontFamily: '"Roboto", sans-serif',
        h1: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
        h2: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
        h3: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
        h4: {
            fontFamily: '"Josefin Sans", "Roboto", sans-serif',
            fontSize: "28px",
            fontWeight: 700,
            lineHeight: 1.2,
        },
        h6: {
            fontFamily: '"Josefin Sans", "Roboto", sans-serif',
            fontSize: "20px",
            fontWeight: 700,
            lineHeight: 1.2,
        },
        body1: { fontSize: "16px", fontWeight: 400, lineHeight: 1.5 },
        body2: { fontSize: "14px", fontWeight: 400, lineHeight: 1.4 },
    },
    components: {
        MuiButton: {
            styleOverrides: {
                root: { borderRadius: "6px" },
            },
            defaultProps: { disableElevation: true },
        },
        MuiAppBar: {
            styleOverrides: {
                root: { backgroundColor: "#1a1a2e" },
            },
        },
        MuiChip: {
            styleOverrides: {
                root: { fontFamily: '"Roboto", sans-serif' },
            },
        },
        MuiTab: {
            styleOverrides: {
                root: { textTransform: "none" },
            },
        },
        // MuiDataGrid: no override needed — adapts automatically when palette.mode is 'dark' (D-15)
    },
});

export default theme;
