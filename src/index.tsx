import React, { createContext, PropsWithChildren, useState } from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./reset.css";
import "./index.sass";

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
    const [theme, setTheme] = useState<Theme>();

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
