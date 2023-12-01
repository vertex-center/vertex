import {
    BrowserRouter,
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
import SettingsApp from "./apps/Settings/SettingsApp/SettingsApp";
import SettingsTheme from "./apps/Settings/SettingsTheme/SettingsTheme";
import { Fragment, useContext } from "react";
import { ThemeContext } from "./main";
import classNames from "classnames";
import SettingsAbout from "./apps/Settings/SettingsAbout/SettingsAbout";
import SettingsUpdates from "./apps/Settings/SettingsUpdates/SettingsUpdates";
import ContainerDocker from "./apps/Containers/pages/ContainerDocker/ContainerDocker";
import ContainerSettings from "./apps/Containers/pages/ContainerSettings/ContainerSettings";
import ContainersStore from "./apps/Containers/pages/ContainersStore/ContainersStore";
import Dock from "./components/Dock/Dock";
import ContainerUpdate from "./apps/Containers/pages/ContainerUpdate/ContainerUpdate";
import ReverseProxyApp from "./apps/ReverseProxy/ReverseProxyApp/ReverseProxyApp";
import SettingsNotifications from "./apps/Settings/SettingsNotifications/SettingsNotifications";
import Header from "./components/Header/Header";
import ContainerDetailsDatabase from "./apps/Containers/pages/ContainerDatabase/ContainerDetailsDatabase";
import SettingsHardware from "./apps/Settings/SettingsHardware/SettingsHardware";
import SettingsSecurity from "./apps/Settings/SettingsSecurity/SettingsSecurity";
import MonitoringApp from "./apps/Monitoring/MonitoringApp/MonitoringApp";
import Prometheus from "./apps/Monitoring/Prometheus/Prometheus";
import MetricsList from "./apps/Monitoring/MetricsList/MetricsList";
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
import SettingsDb from "./apps/Settings/SettingsData/SettingsDb";
import SettingsChecks from "./apps/Settings/SettingsChecks/SettingsChecks";
import Register from "./apps/Auth/pages/Register/Register";
import Logout from "./apps/Auth/pages/Logout/Logout";
import { getAuthToken } from "./backend/api/backend";
import Account from "./apps/Auth/pages/Account/Account";
import AccountInfo from "./apps/Auth/pages/Account/AccountInfo";
import AccountSecurity from "./apps/Auth/pages/Account/AccountSecurity";

const queryClient = new QueryClient();

function AllRoutes() {
    const { pathname } = useLocation();
    const navigate = useNavigate();

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
            navigate("/login");
        }
    }

    return (
        <Fragment>
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
                        </Route>
                        <Route
                            path="/"
                            element={<Navigate to="/app/containers" />}
                            index
                        />
                        <Route
                            path="/app/containers"
                            element={<ContainersApp />}
                        />
                        <Route
                            path="/app/containers/add"
                            element={<ContainersStore />}
                        />
                        <Route
                            path="/app/devtools-service-editor"
                            element={<ServiceEditor />}
                        />
                        <Route path="/app/sql" element={<SqlApp />}>
                            <Route
                                path="/app/sql/install"
                                element={<SqlInstaller />}
                            />
                            <Route
                                path="/app/sql/db/:uuid"
                                element={<SqlDatabase />}
                            />
                        </Route>
                        <Route
                            path="/app/monitoring"
                            element={<MonitoringApp />}
                        >
                            <Route
                                path="/app/monitoring/metrics"
                                element={<MetricsList />}
                            />
                            <Route
                                path="/app/monitoring/prometheus"
                                element={<Prometheus />}
                            />
                            <Route
                                path="/app/monitoring/grafana"
                                element={<Grafana />}
                            />
                        </Route>
                        <Route path="/app/tunnels" element={<TunnelsApp />}>
                            <Route
                                path="/app/tunnels/cloudflare"
                                element={<CloudflareTunnels />}
                            />
                        </Route>
                        <Route
                            path="/app/reverse-proxy"
                            element={<ReverseProxyApp />}
                        >
                            <Route
                                path="/app/reverse-proxy/vertex"
                                element={<VertexReverseProxy />}
                            />
                        </Route>
                        <Route
                            path="/app/containers/:uuid/"
                            element={<ContainerDetails />}
                        >
                            <Route
                                path="/app/containers/:uuid/home"
                                element={<ContainerHome />}
                            />
                            <Route
                                path="/app/containers/:uuid/docker"
                                element={<ContainerDocker />}
                            />
                            <Route
                                path="/app/containers/:uuid/logs"
                                element={<ContainerLogs />}
                            />
                            <Route
                                path="/app/containers/:uuid/environment"
                                element={<ContainerEnv />}
                            />
                            <Route
                                path="/app/containers/:uuid/database"
                                element={<ContainerDetailsDatabase />}
                            />
                            <Route
                                path="/app/containers/:uuid/update"
                                element={<ContainerUpdate />}
                            />
                            <Route
                                path="/app/containers/:uuid/settings"
                                element={<ContainerSettings />}
                            />
                        </Route>
                        <Route path="/settings" element={<SettingsApp />}>
                            <Route
                                path="/settings/theme"
                                element={<SettingsTheme />}
                            />
                            <Route
                                path="/settings/notifications"
                                element={<SettingsNotifications />}
                            />
                            <Route
                                path="/settings/hardware"
                                element={<SettingsHardware />}
                            />
                            <Route
                                path="/settings/database"
                                element={<SettingsDb />}
                            />
                            <Route
                                path="/settings/security"
                                element={<SettingsSecurity />}
                            />
                            <Route
                                path="/settings/updates"
                                element={<SettingsUpdates />}
                            />
                            <Route
                                path="/settings/checks"
                                element={<SettingsChecks />}
                            />
                            <Route
                                path="/settings/about"
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
                <ReactQueryDevtools initialIsOpen={false} />
                <BrowserRouter>
                    <AllRoutes />
                </BrowserRouter>
            </QueryClientProvider>
        </div>
    );
}

export default App;
