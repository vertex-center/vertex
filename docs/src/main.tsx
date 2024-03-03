import React from "react";
import ReactDOM from "react-dom/client";
import "./reset.css";
import "./styles/index.sass";
import { ThemeProvider } from "./theme.tsx";
import { App } from "./App.tsx";
import { PageProvider } from "@vertex-center/components";

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <ThemeProvider>
            <PageProvider>
                <App />
            </PageProvider>
        </ThemeProvider>
    </React.StrictMode>
);
