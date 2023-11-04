import { useContext } from "react";
import { ThemeContext } from "./theme.tsx";
import cx from "classnames";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Documentation from "./pages/Documentation/Documentation.tsx";
import { Header } from "@vertex-center/components";

const router = createBrowserRouter(
    [
        {
            path: "/",
            element: <Documentation />,
        },
    ],
    {
        basename: "/",
    }
);

export function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={cx("app", theme)}>
            <Header />
            <RouterProvider router={router} />
        </div>
    );
}
