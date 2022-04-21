import { FormEvent, useEffect, useState } from "react";

export interface Config {
  port: number;
  mode: string;
  allowLan: boolean;
  logLevel: string;
  proxies: Array<Proxy>;
}

export interface Proxy {
  id: string;
  name: string;
  server: Array<string>;
}

export function ConfigPage() {
  const [port, setPort] = useState(0);
  const [mode, setMode] = useState("");
  const [allowLan, setAllowLan] = useState("true");
  const [logLevel, setLogLevel] = useState("info");

  useEffect(() => {
    const fn = async () => {
      const r = await fetch("/api/config");
      const c: Config = await r.json();
      setPort(c.port);
      setMode(c.mode);
      setAllowLan(c.allowLan ? "true" : "false");
      setLogLevel(c.logLevel);
    };
    fn();
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const r = await fetch("/api/config", {
      method: "PATCH",
      body: JSON.stringify({
        port,
        mode: mode,
        allowLan: allowLan === "true" ? true : false,
        logLevel: logLevel,
      }),
    });
    if (r.ok) {
      alert("successfully update config");
      return;
    }
    console.error("fail to update config", r);
  }

  return (
    <main className="flex flex-col items-center gap-4 mt-8">
      <p className="text-2xl">Config</p>
      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 text-lg"
      >
        <label>
          Port:{"  "}
          <input
            type="number"
            value={port}
            disabled
            className="border p-1 rounded-md dark:bg-gray-600"
          />
        </label>
        <label>
          Mode:{"  "}
          <select
            value={mode}
            onChange={(e) => {
              setMode(e.target.value);
            }}
            className="bg-gray-100 border p-1 rounded-md dark:bg-gray-600"
          >
            <option value={"direct"}>direct</option>
            <option value={"rule"}>rule</option>
            <option value={"global"}>global</option>
          </select>
        </label>

        <label>
          AllowLan:{"  "}
          <select
            value={allowLan}
            onChange={(e) => {
              setAllowLan(e.target.value);
            }}
            className="bg-gray-100 border p-1 rounded-md dark:bg-gray-600"
          >
            <option value="true">true</option>
            <option value="false">false</option>
          </select>
        </label>

        <label>
          LogLevel:{"  "}
          <select
            value={logLevel}
            onChange={(e) => {
              setLogLevel(e.target.value);
            }}
            className="bg-gray-100 border p-1 rounded-md dark:bg-gray-600"
          >
            <option value={"info"}>info</option>
            <option value={"warn"}>warn</option>
            <option value={"error"}>error</option>
          </select>
        </label>

        <button className="bg-gray-200 border p-1 rounded-md dark:bg-gray-600">
          Update Config
        </button>
      </form>
    </main>
  );
}
