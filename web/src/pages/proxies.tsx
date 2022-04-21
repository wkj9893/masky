import { useEffect, useState } from "react";
import type { Config, Proxy } from "./config";

export function ProxiesPage() {
  const [proxies, setProxies] = useState<Array<Proxy>>([]);

  useEffect(() => {
    const fn = async () => {
      const c: Config = await (await fetch("/api/config")).json();
      setProxies(c.proxies);
    };
    fn();
  }, []);

  return (
    <main>
      <p className="text-2xl">Proxies</p>
      {proxies.map((p) => {
        return (
          <ul key={p.id} className="text-xl my-4">
            <li>id: {p.id}</li>
            <li>name: {p.name}</li>
            <li>server: [ {p.server} ]</li>
          </ul>
        );
      })}
    </main>
  );
}
