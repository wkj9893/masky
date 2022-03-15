import { FormEvent, FormEventHandler, useEffect, useState } from "react";
import "./setting.css";

interface Setting {
  port: number;
  mode: Mode;
  addr: string;
  dns: string;
  password: string;
  allowLan: boolean;
  logLevel: LogLevel;
}
type Mode = "direct" | "rule" | "global";
type LogLevel = "info" | "warn" | "error";

export function SettingsPage() {
  const [port, setPort] = useState(0);
  const [mode, setMode] = useState("");
  const [addr, setAddr] = useState("");
  const [dns, setDns] = useState("");
  const [password, setPassword] = useState("");
  const [allowLan, setAllowLan] = useState("true");
  const [logLevel, setLogLevel] = useState("info");

  useEffect(() => {
    const fn = async () => {
      const response = await fetch("/api/setting");
      const setting: Setting = await response.json();
      setPort(setting.port);
      setMode(setting.mode);
      setAddr(setting.addr);
      setDns(setting.dns);
      setPassword(setting.password);
      setAllowLan(setting.allowLan ? "true" : "false");
      setLogLevel(setting.logLevel);
    };
    fn();
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const setting: Setting = {
      port,
      mode: mode as Mode,
      addr,
      dns,
      password,
      allowLan: allowLan === "true" ? true : false,
      logLevel: logLevel as LogLevel,
    };
    const response = await fetch("/api/setting", {
      method: "PATCH",
      body: JSON.stringify(setting),
    });
    if (response.ok) {
      alert("successfully update setting");
      return;
    }
    console.error("fail to update setting", response);
  }

  return (
    <main>
      <h3>Setting</h3>
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
          Addr:
          <input
            value={addr}
            onChange={(e) => {
              setAddr(e.target.value);
            }}
          />
        </label>

        <label>
          Dns:
          <input
            value={dns}
            onChange={(e) => {
              setDns(e.target.value);
            }}
          />
        </label>

        <label>
          Password:
          <input
            type="password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
          />
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

        <button>Update Setting</button>
      </form>
    </main>
  );
}
