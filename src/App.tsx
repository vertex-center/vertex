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
import APIGenerator from "./api.ts";
import AllDocs from "./pages/All/AllDocs.tsx";
import { Api } from "./pages/Api/Api.tsx";

const docs = new Generator("/docs");
const apis = new APIGenerator();

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
                    ...apis.getRoutes().map((route) => ({
                        path: route.path,
                        element: <Api tag={route.tag} api={route.api} />,
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

    let title: string,
        path: string,
        isPage = true;

    const page = docs?.getPage(hierarchy?._path);
    if (page !== undefined) {
        title = page.title;
        path = page.path;
        isPage = page.isPage;
    } else {
        title = hierarchy?._title;
        path = hierarchy?._path;
    }

    let link = {};
    if (isPage) {
        link = { as: NavLink, to: path, end: true };
    }

    console.log(hierarchy);
    const hasChildren = Object.keys(hierarchy ?? {}).length > 2;

    const children =
        typeof hierarchy === "object" &&
        Object.entries(hierarchy ?? {}).map(([label, hierarchy]) => {
            if (label === "_path") return null;
            if (label === "_title") return null;
            return <SidebarItems key={label} hierarchy={hierarchy} />;
        });

    if (root) return children;

    return (
        <Sidebar.Item label={title ?? "---"} link={link}>
            {hasChildren && children}
        </Sidebar.Item>
    );
};

function Docs() {
    const navigate = useNavigate();
    const location = useLocation();

    const app = location.pathname.split("/")?.[1];
    const version = location.pathname.split("/")?.[2];
    if (version === undefined || version === "") navigate(`/${app}/next/`);

    const onVersionChange = (v: unknown) => {
        navigate(`/${app}/${v}/`);
    };

    useHasSidebar(true);

    const hierarchy = docs.hierarchy?.[app] ?? apis.hierarchy?.[app] ?? {};

    return (
        <Fragment>
            <Sidebar>
                <SelectField
                    label="Version"
                    onChange={onVersionChange}
                    value={version}
                >
                    {Object.keys(hierarchy).map((version) => (
                        <SelectOption key={version} value={version}>
                            {version}
                        </SelectOption>
                    ))}
                </SelectField>
                <Sidebar.Group title={version}>
                    <SidebarItems root hierarchy={hierarchy?.[version]} />
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
