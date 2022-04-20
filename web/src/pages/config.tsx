import { FormEvent, useEffect, useState } from "react";

interface Config {
  port: number;
  mode: Mode;
  allowLan: boolean;
  logLevel: LogLevel;
}
type Mode = "direct" | "rule" | "global";
type LogLevel = "info" | "warn" | "error";

export function ConfigPage() {
  const [port, setPort] = useState(0);
  const [mode, setMode] = useState("");
  const [allowLan, setAllowLan] = useState("true");
  const [logLevel, setLogLevel] = useState("info");

  useEffect(() => {
    const fn = async () => {
      const response = await fetch("/api/config");
      const config: Config = await response.json();
      setPort(config.port);
      setMode(config.mode);
      setAllowLan(config.allowLan ? "true" : "false");
      setLogLevel(config.logLevel);
    };
    fn();
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const config: Config = {
      port,
      mode: mode as Mode,
      allowLan: allowLan === "true" ? true : false,
      logLevel: logLevel as LogLevel,
    };
    const response = await fetch("/api/config", {
      method: "PUT",
      body: JSON.stringify(config),
    });
    if (response.ok) {
      alert("successfully update config");
      return;
    }
    console.error("fail to update config", response);
  }

  return (
    <main>
      <h3>Config</h3>
      <form
        onSubmit={handleSubmit}
      >
        <label>
          Port:
          <input
            type="number"
            value={port}
            onChange={(e) => {
              setPort(e.target.valueAsNumber);
            }}
          />
        </label>
        <label>
          Mode:
          <select
            value={mode}
            onChange={(e) => {
              setMode(e.target.value);
            }}
          >
            <option value={"direct"}>direct</option>
            <option value={"rule"}>rule</option>
            <option value={"global"}>global</option>
          </select>
        </label>

        <label>
          AllowLan:
          <select
            value={allowLan}
            onChange={(e) => {
              setAllowLan(e.target.value);
            }}
          >
            <option value="true">true</option>
            <option value="false">false</option>
          </select>
        </label>

        <label>
          LogLevel:
          <select
            value={logLevel}
            onChange={(e) => {
              setLogLevel(e.target.value);
            }}
          >
            <option value={"info"}>info</option>
            <option value={"warn"}>warn</option>
            <option value={"error"}>error</option>
          </select>
        </label>

        <button>Update Config</button>
      </form>
    </main>
  );
}
