import { HashRouter, Route, Routes } from "react-router-dom";
import Marketplace from "./pages/Marketplace/Marketplace";
import Infrastructure from "./pages/Infrastructure/Infrastructure";

function App() {
    return (
        <div className="bg-zinc-50 min-h-screen">
            <HashRouter>
                {/*<Header />*/}
                <Routes>
                    <Route path="/" element={<Infrastructure />} />
                    <Route path="/marketplace" element={<Marketplace />} />
                </Routes>
            </HashRouter>
        </div>
    );
}

export default App;
