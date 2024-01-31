import {
    HashRouter,
    Navigate,
    Route,
    Routes,
    useLocation,
    useNavigate,
} from "react-router-dom";
import ContainersApp from "./apps/Containers/pages/ContainersApp/ContainersApp";
import ContainerDetails from "./apps/Containers/pages/Container/Container";
import ContainerLogs from "./apps/Containers/pages/ContainerLogs/ContainerLogs";
import ContainerEnv from "./apps/Containers/pages/ContainerEnv/ContainerEnv";
import ContainerHome from "./apps/Containers/pages/ContainerHome/ContainerHome";
import SettingsApp from "./apps/AdminSettings/SettingsApp/SettingsApp";
import SettingsTheme from "./apps/AdminSettings/SettingsTheme/SettingsTheme";
import { Fragment, useContext } from "react";
import { ThemeContext } from "./main";
import classNames from "classnames";
import SettingsAbout from "./apps/AdminSettings/SettingsAbout/SettingsAbout";
import SettingsUpdates from "./apps/AdminSettings/SettingsUpdates/SettingsUpdates";
import ContainerDocker from "./apps/Containers/pages/ContainerDocker/ContainerDocker";
import ContainerSettings from "./apps/Containers/pages/ContainerSettings/ContainerSettings";
import ContainersStore from "./apps/Containers/pages/ContainersStore/ContainersStore";
import Dock from "./components/Dock/Dock";
import ContainerUpdate from "./apps/Containers/pages/ContainerUpdate/ContainerUpdate";
import ReverseProxyApp from "./apps/ReverseProxy/ReverseProxyApp/ReverseProxyApp";
import SettingsNotifications from "./apps/AdminSettings/SettingsNotifications/SettingsNotifications";
import Header from "./components/Header/Header";
import ContainerDetailsDatabase from "./apps/Containers/pages/ContainerDatabase/ContainerDetailsDatabase";
import MonitoringApp from "./apps/Monitoring/MonitoringApp/MonitoringApp";
import Prometheus from "./apps/Monitoring/Prometheus/Prometheus";
import Grafana from "./apps/Monitoring/Grafana/Grafana";
import TunnelsApp from "./apps/Tunnels/TunnelsApp/TunnelsApp";
import CloudflareTunnels from "./apps/Tunnels/CloudflareTunnels/CloudflareTunnels";
import VertexReverseProxy from "./apps/ReverseProxy/VertexReverseProxy/VertexReverseProxy";
import SqlApp from "./apps/Sql/SqlApp/SqlApp";
import SqlInstaller from "./apps/Sql/Installer/SqlInstaller";
import SqlDatabase from "./apps/Sql/SqlDatabase/SqlDatabase";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import ServiceEditor from "./apps/DevToolsServiceEditor/ServiceEditor/ServiceEditor";
import Login from "./apps/Auth/pages/Login/Login";
import SettingsDb from "./apps/AdminSettings/SettingsData/SettingsDb";
import SettingsChecks from "./apps/AdminSettings/SettingsChecks/SettingsChecks";
import Register from "./apps/Auth/pages/Register/Register";
import Logout from "./apps/Auth/pages/Logout/Logout";
import Account from "./apps/Auth/pages/Account/Account";
import AccountInfo from "./apps/Auth/pages/Account/AccountInfo";
import AccountSecurity from "./apps/Auth/pages/Account/AccountSecurity";
import useUser from "./apps/Auth/hooks/useUser";
import AccountEmails from "./apps/Auth/pages/Account/AccountEmails";
import { getAuthToken } from "./backend/server";

const queryClient = new QueryClient();

function AllRoutes() {
    const { pathname } = useLocation();
    const navigate = useNavigate();
    const { errorUser } = useUser();

    let show = {
        header: true,
        dock: true,
    };

    if (
        pathname === "/login" ||
        pathname === "/register" ||
        pathname === "/logout"
    ) {
        show = {
            header: false,
            dock: false,
        };
    } else {
        const token = getAuthToken();
        if (!token) {
            console.log("No token found. Redirecting to login.");
            navigate("/login");
        }

        // @ts-ignore
        if (errorUser && errorUser?.response?.status === 401) {
            console.log("The token is invalid. Redirecting to login.");
            navigate("/login");
        }
    }

    return (
        <Fragment>
            <ReactQueryDevtools initialIsOpen={false} />
            {show.header && <Header />}
            <div className="app-main">
                <div className="app-sidebar" />
                <div className="app-content">
                    <Routes>
                        <Route path="/login" element={<Login />} />
                        <Route path="/register" element={<Register />} />
                        <Route path="/logout" element={<Logout />} />
                        <Route path="/account" element={<Account />}>
                            <Route
                                path="/account/info"
                                element={<AccountInfo />}
                            />
                            <Route
                                path="/account/security"
                                element={<AccountSecurity />}
                            />
                            <Route
                                path="/account/emails"
                                element={<AccountEmails />}
                            />
                        </Route>
                        <Route
                            path="/"
                            element={<Navigate to="/containers" />}
                            index
                        />
                        <Route path="/containers" element={<ContainersApp />} />
                        <Route
                            path="/containers/add"
                            element={<ContainersStore />}
                        />
                        <Route
                            path="/devtools-service-editor"
                            element={<ServiceEditor />}
                        />
                        <Route path="/sql" element={<SqlApp />}>
                            <Route
                                path="/sql/install"
                                element={<SqlInstaller />}
                            />
                            <Route
                                path="/sql/db/:uuid"
                                element={<SqlDatabase />}
                            />
                        </Route>
                        <Route path="/monitoring" element={<MonitoringApp />}>
                            <Route
                                path="/monitoring/prometheus"
                                element={<Prometheus />}
                            />
                            <Route
                                path="/monitoring/grafana"
                                element={<Grafana />}
                            />
                        </Route>
                        <Route path="/tunnels" element={<TunnelsApp />}>
                            <Route
                                path="/tunnels/cloudflare"
                                element={<CloudflareTunnels />}
                            />
                        </Route>
                        <Route
                            path="/reverse-proxy"
                            element={<ReverseProxyApp />}
                        >
                            <Route
                                path="/reverse-proxy/vertex"
                                element={<VertexReverseProxy />}
                            />
                        </Route>
                        <Route
                            path="/containers/:uuid/"
                            element={<ContainerDetails />}
                        >
                            <Route
                                path="/containers/:uuid/home"
                                element={<ContainerHome />}
                            />
                            <Route
                                path="/containers/:uuid/docker"
                                element={<ContainerDocker />}
                            />
                            <Route
                                path="/containers/:uuid/logs"
                                element={<ContainerLogs />}
                            />
                            <Route
                                path="/containers/:uuid/environment"
                                element={<ContainerEnv />}
                            />
                            <Route
                                path="/containers/:uuid/database"
                                element={<ContainerDetailsDatabase />}
                            />
                            <Route
                                path="/containers/:uuid/update"
                                element={<ContainerUpdate />}
                            />
                            <Route
                                path="/containers/:uuid/settings"
                                element={<ContainerSettings />}
                            />
                        </Route>
                        <Route path="/admin" element={<SettingsApp />}>
                            <Route
                                path="/admin/theme"
                                element={<SettingsTheme />}
                            />
                            <Route
                                path="/admin/notifications"
                                element={<SettingsNotifications />}
                            />
                            <Route
                                path="/admin/database"
                                element={<SettingsDb />}
                            />
                            <Route
                                path="/admin/updates"
                                element={<SettingsUpdates />}
                            />
                            <Route
                                path="/admin/checks"
                                element={<SettingsChecks />}
                            />
                            <Route
                                path="/admin/about"
                                element={<SettingsAbout />}
                            />
                        </Route>
                    </Routes>
                </div>
            </div>
            {show.dock && <Dock />}
        </Fragment>
    );
}

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div id="app" className={classNames("app", theme)}>
            <QueryClientProvider client={queryClient}>
                <HashRouter>
                    <AllRoutes />
                </HashRouter>
            </QueryClientProvider>
        </div>
    );
}

export default App;
