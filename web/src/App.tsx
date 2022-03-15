import { Dashboard } from "../components/Dashboard";
import { Link, Outlet, Route, Routes } from "react-router-dom";
import { AboutPage } from "./pages/about";
import { IndexPage } from "./pages";
import { SettingsPage } from "./pages/settings";
import { LogsPage } from "./pages/logs";
import "./app.css";

export default function App() {
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="/" element={<IndexPage />} />
          <Route path="/setting" element={<SettingsPage />} />
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
    <>
      <h3>Nothing to see here</h3>
      <Link to="/">
        <h3>Go to the home page</h3>
      </Link>
    </>
  );
}
