import { createContext, PropsWithChildren } from "react";

export const ThemeContext = createContext<{
    theme?: string;
    setTheme?: unknown;
}>({
    theme: undefined,
    setTheme: undefined,
});

type ThemeProviderProps = PropsWithChildren<{
    theme?: string;
    setTheme?: unknown;
}>;

export default function ThemeProvider(props: ThemeProviderProps) {
    const { theme, setTheme, children } = props;

    return (
        <ThemeContext.Provider value={{ theme, setTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}
