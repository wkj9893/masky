import { Dashboard } from "./components/Dashboard";
import { Link, Outlet, Route, Routes } from "react-router-dom";
import { AboutPage } from "./pages/about";
import { IndexPage } from "./pages";
import { SettingsPage } from "./pages/settings";
import { LogsPage } from "./pages/logs";
import { useTheme } from "./hooks/useTheme";
import "./app.scss";
import { ConnectionsPage } from "./pages/connections";

export default function App() {
  const [theme, setTheme] = useTheme();
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="/" element={<IndexPage />} />
          <Route path="/setting" element={<SettingsPage />} />
          <Route path="/connections" element={<ConnectionsPage />} />
          <Route path="/logs" element={<LogsPage />} />
          <Route path="/about" element={<AboutPage />} />
          <Route path="*" element={<NoMatch />}></Route>
        </Route>
      </Routes>
      <button
        style={{
          position: "absolute",
          right: "100px",
          width: "60px",
          height: "40px",
        }}
        onClick={() => {
          setTheme(theme === "light" ? "dark" : "light");
        }}
      >
        Switch Theme
      </button>
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
