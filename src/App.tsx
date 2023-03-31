import { HashRouter, Route, Routes } from "react-router-dom";
import Marketplace from "./pages/Marketplace/Marketplace";
import Infrastructure from "./pages/Infrastructure/Infrastructure";
import BayDetails from "./pages/BayDetails/BayDetails";
import BayDetailsLogs from "./pages/BayDetailsLogs/BayDetailsLogs";
import BayDetailsEnv from "./pages/BayDetailsEnv/BayDetailsEnv";
import BayDetailsDependencies from "./pages/BayDetailsDependencies";
import BayDetailsHome from "./pages/BayDetailsHome/BayDetailsHome";

function App() {
    return (
        <HashRouter>
            {/*<Header />*/}
            <Routes>
                <Route path="/" element={<Infrastructure />} index />
                <Route path="/marketplace" element={<Marketplace />} />
                <Route path="/bay/:uuid/" element={<BayDetails />}>
                    <Route path="/bay/:uuid/" element={<BayDetailsHome />} />
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
    );
}

export default App;
