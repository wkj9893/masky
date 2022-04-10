import { Link, Outlet, Route, Routes } from "react-router-dom";
import { Dashboard } from "./components/Dashboard";
import { AboutPage } from "./pages/about";
import { IndexPage } from "./pages";
import { ConfigPage } from "./pages/config";
import { LogsPage } from "./pages/logs";
import { useTheme } from "./hooks/useTheme";
import { ThemeButton } from "./components/ThemeButton";
import { ProxiesPage } from "./pages/proxies";
import "./app.scss";

export default function App() {
  const [theme, setTheme] = useTheme();
  return (
    <div className="app">
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
      <ThemeButton
        theme={theme}
        setTheme={setTheme}
        style={{
          position: "absolute",
          right: "100px",
          display: "flex",
        }}
      />
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
    <>
      <h3>Nothing to see here</h3>
      <Link to="/">
        <h3>Go to the home page</h3>
      </Link>
    </>
  );
}
