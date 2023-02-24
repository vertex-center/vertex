import Header from "./components/Header/Header";
import { HashRouter, Route, Routes } from "react-router-dom";
import Marketplace from "./pages/Marketplace/Marketplace";
import Installed from "./pages/Installed";

function App() {
    return (
        <div className="bg-zinc-50 min-h-screen">
            <HashRouter>
                <Header />
                <Routes>
                    <Route path="/marketplace" element={<Marketplace />} />
                    <Route path="/installed" element={<Installed />} />
                    <Route path="/" element={<></>} />
                </Routes>
            </HashRouter>
        </div>
    );
}

export default App;
