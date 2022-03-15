import { useState } from "react";

export function LogsPage() {
  const [logs, setLogs] = useState(["hello", "world"]);
  return (
    <main>
      <h3>Logs</h3>
      {logs.map((log) => {
        <p>{log}</p>;
      })}
    </main>
  );
}
