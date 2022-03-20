import { useEffect, useState } from "react";

export function useTheme(): [
  string,
  React.Dispatch<React.SetStateAction<string>>,
] {
  const [theme, setTheme] = useState("");

  useEffect(() => {
    if (!theme) {
      const t = localStorage.getItem("theme");
      if (t) {
        setTheme(t);
        document.querySelector(":root")!.setAttribute("theme", t);
        return;
      }
      if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
        setTheme("dark");
        document.querySelector(":root")!.setAttribute("theme", "dark");
        localStorage.setItem("theme", "dark");
      } else {
        setTheme("light");
        document.querySelector(":root")!.setAttribute("theme", "light");
        localStorage.setItem("theme", "light");
      }
    } else {
      document.querySelector(":root")!.setAttribute("theme", theme);
      localStorage.setItem("theme", theme);
    }
  }, [theme]);
  return [theme, setTheme];
}
