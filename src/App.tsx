import { HashRouter, Route, Routes } from "react-router-dom";
import Marketplace from "./pages/Marketplace/Marketplace";
import Infrastructure from "./pages/Infrastructure/Infrastructure";
import BayDetails from "./pages/BayDetails/BayDetails";
import BayDetailsLogs from "./pages/BayDetailsLogs/BayDetailsLogs";
import BayDetailsEnv from "./pages/BayDetailsEnv/BayDetailsEnv";
import BayDetailsDependencies from "./pages/BayDetailsDependencies";
import BayDetailsHome from "./pages/BayDetailsHome/BayDetailsHome";
import Settings from "./pages/Settings/Settings";
import SettingsTheme from "./pages/SettingsTheme/SettingsTheme";
import { useContext } from "react";
import { ThemeContext } from "./main";
import classNames from "classnames";
import SettingsAbout from "./pages/SettingsAbout/SettingsAbout";
import SettingsUpdates from "./pages/SettingsUpdates/SettingsUpdates";

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={classNames("app", theme)}>
            <HashRouter>
                <Routes>
                    <Route path="/" element={<Infrastructure />} index />
                    <Route path="/settings" element={<Settings />}>
                        <Route
                            path="/settings/theme"
                            element={<SettingsTheme />}
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
                    <Route path="/marketplace" element={<Marketplace />} />
                    <Route path="/bay/:uuid/" element={<BayDetails />}>
                        <Route
                            path="/bay/:uuid/"
                            element={<BayDetailsHome />}
                        />
                        <Route
                            path="/bay/:uuid/logs"
                            element={<BayDetailsLogs />}
                        />
                        <Route
                            path="/bay/:uuid/environment"
                            element={<BayDetailsEnv />}
                        />
                        <Route
                            path="/bay/:uuid/dependencies"
                            element={<BayDetailsDependencies />}
                        />
                    </Route>
                </Routes>
            </HashRouter>
        </div>
    );
}

export default App;
