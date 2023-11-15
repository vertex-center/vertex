import { useContext, useState } from "react";
import { ThemeContext } from "./theme.tsx";
import cx from "classnames";
import {
    createBrowserRouter,
    NavLink,
    Outlet,
    RouterProvider,
} from "react-router-dom";
import {
    Header,
    SelectField,
    SelectOption,
    Sidebar,
    useHasSidebar,
} from "@vertex-center/components";
import Documentation from "./pages/Documentation/Documentation.tsx";
import Docs from "./documentation.ts";

const docs = new Docs("/docs");

const router = createBrowserRouter(
    [
        {
            element: <Root />,
            children: [
                {
                    path: "/",
                    element: <Documentation content={undefined} />,
                },
                ...docs.getRoutes().map((route) => ({
                    path: route.path,
                    element: <Documentation content={route.page?.default} />,
                })),
            ],
        },
    ],
    { basename: "/" }
);

type SidebarItemsProps = {
    root?: boolean;
    hierarchy: any;
};

const SidebarItems = (props: SidebarItemsProps) => {
    const { root, hierarchy } = props;

    console.log(docs.pages);
    const { title, path, isPage } = docs?.getPage(hierarchy?._path) ?? {};

    let link = {};
    if (isPage) {
        link = { as: NavLink, to: path, end: true };
    }

    const hasChildren = Object.keys(hierarchy ?? {}).length > 1;

    const children =
        hasChildren &&
        Object.entries(hierarchy ?? {}).map(([label, hierarchy], i) => {
            if (label === "_path") return null;
            return <SidebarItems key={i} label={label} hierarchy={hierarchy} />;
        });

    if (root) return children;

    return (
        <Sidebar.Item label={title ?? "---"} link={link}>
            {children}
        </Sidebar.Item>
    );
};

function Root() {
    const { theme } = useContext(ThemeContext);

    const [version, setVersion] = useState("next");

    useHasSidebar(true);

    return (
        <div className={cx("app", theme)}>
            <Header />
            <div className="app-content">
                <Sidebar>
                    <SelectField
                        label="Version"
                        onChange={(v) => setVersion(v)}
                        value={version}
                    >
                        {Object.keys(docs.hierarchy).map((version) => (
                            <SelectOption key={version} value={version}>
                                {version}
                            </SelectOption>
                        ))}
                    </SelectField>
                    <Sidebar.Group title={version}>
                        <SidebarItems
                            root
                            hierarchy={docs.hierarchy[version]}
                        />
                    </Sidebar.Group>
                </Sidebar>
                <Outlet />
            </div>
        </div>
    );
}

export function App() {
    return <RouterProvider router={router} />;
}
