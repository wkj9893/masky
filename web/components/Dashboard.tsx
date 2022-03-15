import { Link } from "react-router-dom";
import "./Dashboard.css";

export function Dashboard() {
  return (
    <nav>
      <header>
        Masky Client
      </header>
      <div className="container ">
        <Link to="/">
          Overview
        </Link>
        <Link to="/setting">
          Setting
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
