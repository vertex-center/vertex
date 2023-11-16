import { Fragment, useContext } from "react";
import { ThemeContext } from "./theme.tsx";
import cx from "classnames";
import {
    createBrowserRouter,
    NavLink,
    Link as RouterLink,
    LinkProps as RouterLinkProps,
    Outlet,
    RouterProvider,
    useLocation,
    useNavigate,
} from "react-router-dom";
import {
    Header,
    LinkProps,
    SelectField,
    SelectOption,
    Sidebar,
    useHasSidebar,
} from "@vertex-center/components";
import Documentation from "./pages/Documentation/Documentation.tsx";
import Generator from "./documentation.ts";
import AllDocs from "./pages/All/AllDocs.tsx";

const docs = new Generator("/docs");

const router = createBrowserRouter([
    {
        element: <Root />,
        children: [
            {
                path: "/",
                element: <AllDocs />,
            },
            {
                element: <Docs />,
                children: [
                    ...docs.getRoutes().map((route) => ({
                        path: route.path,
                        element: (
                            <Documentation content={route.page?.default} />
                        ),
                    })),
                    {
                        path: "*",
                        element: <Documentation content={"div"} />,
                    },
                ],
            },
        ],
    },
]);

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
        Object.entries(hierarchy ?? {}).map(([label, hierarchy]) => {
            if (label === "_path") return null;
            return <SidebarItems key={label} hierarchy={hierarchy} />;
        });

    if (root) return children;

    return (
        <Sidebar.Item label={title ?? "---"} link={link}>
            {children}
        </Sidebar.Item>
    );
};

function Docs() {
    const navigate = useNavigate();
    const location = useLocation();

    const app = location.pathname.split("/")?.[1];
    const version = location.pathname.split("/")?.[2];
    console.log(app, version);
    if (version === undefined || version === "") navigate(`/${app}/next/`);

    const onVersionChange = (v: unknown) => {
        navigate(`/${app}/${v}/`);
    };

    useHasSidebar(true);

    return (
        <Fragment>
            <Sidebar>
                <SelectField
                    label="Version"
                    onChange={onVersionChange}
                    value={version}
                >
                    {Object.keys(docs.hierarchy).map((version) => (
                        <SelectOption key={version} value={version}>
                            {version}
                        </SelectOption>
                    ))}
                </SelectField>
                <Sidebar.Group title={version}>
                    <SidebarItems root hierarchy={docs.hierarchy[version]} />
                </Sidebar.Group>
            </Sidebar>
            <Outlet />
        </Fragment>
    );
}

export function Root() {
    const { theme } = useContext(ThemeContext);

    const linkLogo: LinkProps<RouterLinkProps> = { as: RouterLink, to: "/" };

    return (
        <div className={cx("app", theme)}>
            <Header appName="Vertex Docs" linkLogo={linkLogo} />
            <div className="app-content">
                <Outlet />
            </div>
        </div>
    );
}

export function App() {
    return <RouterProvider router={router} />;
}
