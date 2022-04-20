import { Link, Outlet, Route, Routes } from "react-router-dom";
import { Dashboard } from "./components/Dashboard";
import { AboutPage } from "./pages/about";
import { IndexPage } from "./pages";
import { ConfigPage } from "./pages/config";
import { LogsPage } from "./pages/logs";
import { ProxiesPage } from "./pages/proxies";

export default function App() {
  return (
    <div className="flex">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="/" element={<IndexPage />} />
          <Route path="/config" element={<ConfigPage />} />
          <Route path="/proxies" element={<ProxiesPage />} />
          <Route path="/logs" element={<LogsPage />} />
          <Route path="/about" element={<AboutPage />} />
          <Route path="*" element={<NoMatch />}></Route>
        </Route>
      </Routes>
    </div>
  );
}

function Layout() {
  return (
    <>
      <Dashboard />
      <Outlet />
    </>
  );
}

function NoMatch() {
  return (
    <div className="flex flex-col items-center gap-4">
      <p className="text-xl">Nothing to see here</p>
      <Link to="/">
        <p className="text-2xl underline underline-offset-8">
          Go to the home page
        </p>
      </Link>
    </div>
  );
}
