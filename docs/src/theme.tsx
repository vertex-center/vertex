import { createContext, PropsWithChildren, useEffect, useState } from "react";
import { useCookies } from "react-cookie";

export type Theme = "theme-vertex-dark" | "theme-vertex-light";

export const ThemeContext = createContext<{
    theme: string;
    setTheme: any;
}>({
    theme: "theme-vertex-dark",
    setTheme: undefined,
});

export function ThemeProvider({ children }: PropsWithChildren) {
    const [cookies, setCookie] = useCookies(["theme"]);
    const [theme, setTheme] = useState<Theme>(cookies.theme);

    useEffect(() => {
        if (cookies.theme !== theme) setCookie("theme", theme);
    }, [cookies.theme, setCookie, theme]);

    useEffect(() => {
        if (theme === undefined) {
            setTheme("theme-vertex-dark");
            return;
        }
        setTheme(theme);
    }, [theme]);

    return (
        <ThemeContext.Provider value={{ theme, setTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}
