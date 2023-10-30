import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Documentation from "./pages/Documentation/Documentation.tsx";
import Home from "./pages/Home/Home.tsx";
import "@vertex-center/components/dist/style.css";

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

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);
