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
import { useCookies } from "react-cookie";

const root = ReactDOM.createRoot(document.getElementById("root"));

export type Theme = "vertex-dark" | "vertex-light";

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

    return (
        <ThemeContext.Provider value={{ theme, setTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}

root.render(
    <React.StrictMode>
        <ThemeProvider>
            <App />
        </ThemeProvider>
    </React.StrictMode>
);
