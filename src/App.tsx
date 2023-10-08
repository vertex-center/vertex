import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import InstancesApp from "./apps/Instances/InstancesApp/InstancesApp";
import Instance from "./apps/Instances/Details/Instance/Instance";
import InstanceLogs from "./apps/Instances/Details/InstanceLogs/InstanceLogs";
import InstanceEnv from "./apps/Instances/Details/InstanceEnv/InstanceEnv";
import InstanceHome from "./apps/Instances/Details/InstanceHome/InstanceHome";
import SettingsApp from "./apps/Settings/SettingsApp/SettingsApp";
import SettingsTheme from "./apps/Settings/SettingsTheme/SettingsTheme";
import { useContext } from "react";
import { ThemeContext } from "./main";
import classNames from "classnames";
import SettingsAbout from "./apps/Settings/SettingsAbout/SettingsAbout";
import SettingsUpdates from "./apps/Settings/SettingsUpdates/SettingsUpdates";
import InstanceDocker from "./apps/Instances/Details/InstanceDocker/InstanceDocker";
import InstanceSettings from "./apps/Instances/Details/InstanceSettings/InstanceSettings";
import InstancesStore from "./apps/Instances/InstancesStore/InstancesStore";
import Dock from "./components/Dock/Dock";
import InstanceUpdate from "./apps/Instances/Details/InstanceUpdate/InstanceUpdate";
import ReverseProxyApp from "./apps/ReverseProxy/ReverseProxyApp/ReverseProxyApp";
import SettingsNotifications from "./apps/Settings/SettingsNotifications/SettingsNotifications";
import Header from "./components/Header/Header";
import InstanceDetailsDatabase from "./apps/Instances/Details/InstanceDatabase/InstanceDetailsDatabase";
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

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={classNames("app", theme)}>
            <HashRouter>
                <Header />
                <div className="appContent">
                    <Routes>
                        <Route
                            path="/"
                            element={<Navigate to="/app/vx-instances" />}
                            index
                        />
                        <Route
                            path="/app/vx-instances"
                            element={<InstancesApp />}
                        />
                        <Route
                            path="/app/vx-instances/add"
                            element={<InstancesStore />}
                        />
                        <Route path="/app/vx-sql" element={<SqlApp />}>
                            <Route
                                path="/app/vx-sql/install"
                                element={<SqlInstaller />}
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
                        <Route path="/app/vx-tunnels" element={<TunnelsApp />}>
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
                            path="/app/vx-instances/:uuid/"
                            element={<Instance />}
                        >
                            <Route
                                path="/app/vx-instances/:uuid/home"
                                element={<InstanceHome />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/docker"
                                element={<InstanceDocker />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/logs"
                                element={<InstanceLogs />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/environment"
                                element={<InstanceEnv />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/database"
                                element={<InstanceDetailsDatabase />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/update"
                                element={<InstanceUpdate />}
                            />
                            <Route
                                path="/app/vx-instances/:uuid/settings"
                                element={<InstanceSettings />}
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
        </div>
    );
}

export default App;
