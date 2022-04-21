import { Link } from "react-router-dom";

export function Dashboard() {
  return (
    <nav className="grid grid-rows-[1fr_2fr] place-items-center">
      <header className="text-4xl flex flex-col">
        <p>Masky</p>
        <p>Client</p>
      </header>
      <ul className="bg-gray-200 dark:bg-gray-800 text-2xl flex flex-col gap-2 p-6 divide-y divide-slate-200">
        <li>
          <Link to="/">
            Overview
          </Link>
        </li>
        <li>
          <Link to="/config">
            Config
          </Link>
        </li>
        <li>
          <Link to="/proxies">
            Proxies
          </Link>
        </li>
        <li>
          <Link to="/logs">
            Logs
          </Link>
        </li>
        <li>
          <Link to="/about">
            About
          </Link>
        </li>
      </ul>
    </nav>
  );
}
