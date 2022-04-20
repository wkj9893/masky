import { Link } from "react-router-dom";

export function Dashboard() {
  return (
    <nav className="flex flex-col gap-4 relative left-2 top-2 bg-gray-200 w-1/5">
      <header className="text-4xl">
        Masky Client
      </header>
      <div className="text-2xl flex flex-col">
        <Link to="/">
          Overview
        </Link>
        <Link to="/config">
          Config
        </Link>
        <Link to="/proxies">
          Proxies
        </Link>
        <Link to="/logs">
          Logs
        </Link>
        <Link to="/about">
          About
        </Link>
      </div>
    </nav>
  );
}
