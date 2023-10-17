import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import ContainersApp from "./apps/Containers/pages/ContainersApp/ContainersApp";
import ContainerDetails from "./apps/Containers/pages/Container/Container";
import ContainerLogs from "./apps/Containers/pages/ContainerLogs/ContainerLogs";
import ContainerEnv from "./apps/Containers/pages/ContainerEnv/ContainerEnv";
import ContainerHome from "./apps/Containers/pages/ContainerHome/ContainerHome";
import SettingsApp from "./apps/Settings/SettingsApp/SettingsApp";
import SettingsTheme from "./apps/Settings/SettingsTheme/SettingsTheme";
import { useContext } from "react";
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

const queryClient = new QueryClient();

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={classNames("app", theme)}>
            <QueryClientProvider client={queryClient}>
                <ReactQueryDevtools initialIsOpen={false} />
                <HashRouter>
                    <Header />
                    <div className="appContent">
                        <Routes>
                            <Route
                                path="/"
                                element={<Navigate to="/app/vx-containers" />}
                                index
                            />
                            <Route
                                path="/app/vx-containers"
                                element={<ContainersApp />}
                            />
                            <Route
                                path="/app/vx-containers/add"
                                element={<ContainersStore />}
                            />
                            <Route
                                path="/app/vx-devtools-service-editor"
                                element={<ServiceEditor />}
                            />
                            <Route path="/app/vx-sql" element={<SqlApp />}>
                                <Route
                                    path="/app/vx-sql/install"
                                    element={<SqlInstaller />}
                                />
                                <Route
                                    path="/app/vx-sql/db/:uuid"
                                    element={<SqlDatabase />}
                                />
                            </Route>
                            <Route
                                path="/app/vx-monitoring"
                                element={<MonitoringApp />}
                            >
                                <Route
                                    path="/app/vx-monitoring/metrics"
                                    element={<MetricsList />}
                                />
                                <Route
                                    path="/app/vx-monitoring/prometheus"
                                    element={<Prometheus />}
                                />
                                <Route
                                    path="/app/vx-monitoring/grafana"
                                    element={<Grafana />}
                                />
                            </Route>
                            <Route
                                path="/app/vx-tunnels"
                                element={<TunnelsApp />}
                            >
                                <Route
                                    path="/app/vx-tunnels/cloudflare"
                                    element={<CloudflareTunnels />}
                                />
                            </Route>
                            <Route
                                path="/app/vx-reverse-proxy"
                                element={<ReverseProxyApp />}
                            >
                                <Route
                                    path="/app/vx-reverse-proxy/vertex"
                                    element={<VertexReverseProxy />}
                                />
                            </Route>
                            <Route
                                path="/app/vx-containers/:uuid/"
                                element={<ContainerDetails />}
                            >
                                <Route
                                    path="/app/vx-containers/:uuid/home"
                                    element={<ContainerHome />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/docker"
                                    element={<ContainerDocker />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/logs"
                                    element={<ContainerLogs />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/environment"
                                    element={<ContainerEnv />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/database"
                                    element={<ContainerDetailsDatabase />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/update"
                                    element={<ContainerUpdate />}
                                />
                                <Route
                                    path="/app/vx-containers/:uuid/settings"
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
                                    path="/settings/security"
                                    element={<SettingsSecurity />}
                                />
                                <Route
                                    path="/settings/updates"
                                    element={<SettingsUpdates />}
                                />
                                <Route
                                    path="/settings/about"
                                    element={<SettingsAbout />}
                                />
                            </Route>
                        </Routes>
                    </div>
                    <Dock />
                </HashRouter>
            </QueryClientProvider>
        </div>
    );
}

export default App;
