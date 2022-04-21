import { useState } from "react";

export function LogsPage() {
  const [logs, setLogs] = useState(["hello", "world"]);
  return (
    <main>
      <p className="text-2xl">Logs</p>
    </main>
  );
}
