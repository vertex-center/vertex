import React, {
    createContext,
    PropsWithChildren,
    useEffect,
    useState,
} from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./reset.css";
import "./index.sass";
import "@vertex-center/components/dist/style.css";
import { useCookies } from "react-cookie";
import { themes } from "./models/theme";
import { HeaderProvider } from "./components/Header/Header";

export type Theme =
    | "theme-vertex-dark"
    | "theme-vertex-light"
    | "catppuccin-mocha"
    | "catppuccin-macchiato"
    | "catppuccin-frappe"
    | "catppuccin-latte";

export const ThemeContext = createContext<{
    theme: string;
    setTheme: any;
}>({
    theme: undefined,
    setTheme: undefined,
});

function ThemeProvider({ children }: PropsWithChildren) {
    const [cookies, setCookie] = useCookies(["theme"]);
    const [theme, setTheme] = useState<Theme>(cookies.theme);

    useEffect(() => {
        if (cookies.theme !== theme) setCookie("theme", theme);
    }, [cookies.theme, setCookie, theme]);

    useEffect(() => {
        const t = themes.find((t) => t.key === theme);
        if (t === undefined) {
            setTheme("theme-vertex-dark");
        }
    }, [theme]);

    return (
        <ThemeContext.Provider value={{ theme, setTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
    <React.StrictMode>
        <ThemeProvider>
            <HeaderProvider>
                <App />
            </HeaderProvider>
        </ThemeProvider>
    </React.StrictMode>
);
