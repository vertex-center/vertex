import React, { useContext } from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Documentation from "./pages/Documentation/Documentation.tsx";
import Home from "./pages/Home/Home.tsx";
import "@vertex-center/components/dist/style.css";
import "./reset.css";
import "./styles/index.sass";
import { ThemeContext, ThemeProvider } from "./theme.tsx";
import cx from "classnames";

const router = createBrowserRouter(
    [
        {
            path: "/",
            element: <Home />,
        },
        {
            path: "/docs",
            element: <Documentation />,
        },
    ],
    {
        basename: "/",
    }
);

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={cx("app", theme)}>
            <RouterProvider router={router} />
        </div>
    );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <ThemeProvider>
            <App />
        </ThemeProvider>
    </React.StrictMode>
);
