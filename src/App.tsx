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
                            element={<Navigate to="/instances" />}
                            index
                        />
                        <Route path="/instances" element={<InstancesApp />} />
                        <Route
                            path="/instances/add"
                            element={<InstancesStore />}
                        />
                        <Route path="/proxy" element={<ReverseProxyApp />} />
                        <Route path="/instances/:uuid/" element={<Instance />}>
                            <Route
                                path="/instances/:uuid/home"
                                element={<InstanceHome />}
                            />
                            <Route
                                path="/instances/:uuid/docker"
                                element={<InstanceDocker />}
                            />
                            <Route
                                path="/instances/:uuid/logs"
                                element={<InstanceLogs />}
                            />
                            <Route
                                path="/instances/:uuid/environment"
                                element={<InstanceEnv />}
                            />
                            <Route
                                path="/instances/:uuid/database"
                                element={<InstanceDetailsDatabase />}
                            />
                            <Route
                                path="/instances/:uuid/update"
                                element={<InstanceUpdate />}
                            />
                            <Route
                                path="/instances/:uuid/settings"
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
