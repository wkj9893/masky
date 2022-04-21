import { Link, Outlet, Route, Routes } from "react-router-dom";
import { Dashboard } from "./components/Dashboard";
import { AboutPage } from "./pages/about";
import { IndexPage } from "./pages";
import { ConfigPage } from "./pages/config";
import { LogsPage } from "./pages/logs";
import { ProxiesPage } from "./pages/proxies";
import { ThemeButton } from "./components/ThemeButton";
import { useTheme } from "./hooks/useTheme";

export default function App() {
  return (
    <div className="dark:bg-[#121212] dark:text-white grid grid-cols-[1fr_4fr_1fr] place-items-center">
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
  const [theme, setTheme] = useTheme();
  return (
    <>
      <Dashboard />
      <Outlet />
      <ThemeButton
        theme={theme}
        setTheme={setTheme}
        className="p-2"
      />
    </>
  );
}

function NoMatch() {
  return (
    <div className="grid grid-rows-[1fr_1fr_4fr] place-items-center">
      <p className="row-start-2 row-end-3 text-xl">
        Nothing to see here
      </p>
      <Link
        to="/"
        className="row-start-3 row-end-4 text-3xl underline underline-offset-8"
      >
        Back to the home page
      </Link>
    </div>
  );
}
