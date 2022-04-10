import { Link } from "react-router-dom";
import "./Dashboard.scss";

export function Dashboard() {
  return (
    <nav>
      <header>
        Masky Client
      </header>
      <div className="container">
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
