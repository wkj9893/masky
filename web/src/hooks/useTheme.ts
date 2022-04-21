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
        if (t == "dark") {
          document.documentElement.setAttribute("class", "dark");
        } else {
          document.documentElement.removeAttribute("class");
        }
      } else if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
        setTheme("dark");
        document.documentElement.setAttribute("class", "dark");
        localStorage.setItem("theme", "dark");
      } else {
        setTheme("light");
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
      }
    } else {
      if (theme == "dark") {
        document.documentElement.setAttribute("class", "dark");
      } else {
        document.documentElement.removeAttribute("class");
      }
      localStorage.setItem("theme", theme);
    }
  }, [theme]);
  return [theme, setTheme];
}
