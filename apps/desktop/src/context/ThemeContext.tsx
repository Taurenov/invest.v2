import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { fetchSettings, updateSettings, type UserSettings } from "../api/client";
import i18n from "../i18n";

type ThemeCtx = {
  settings: UserSettings | null;
  apply: (patch: Partial<UserSettings>) => Promise<void>;
};

const Ctx = createContext<ThemeCtx | null>(null);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [settings, setSettings] = useState<UserSettings | null>(null);

  useEffect(() => {
    const cached = localStorage.getItem("fin_settings_cache");
    if (cached) {
      try {
        const parsed = JSON.parse(cached) as UserSettings;
        setSettings(parsed);
        applyDom(parsed);
      } catch {
        /* ignore */
      }
    }
    fetchSettings()
      .then((r) => {
        setSettings(r.data);
        applyDom(r.data);
        localStorage.setItem("fin_settings_cache", JSON.stringify(r.data));
        i18n.changeLanguage(r.data.locale);
      })
      .catch(() => {
        const fallback: UserSettings = {
          locale: "ru",
          base_currency: "RUB",
          theme: "dark",
          timezone: "Europe/Moscow",
        };
        setSettings(fallback);
        applyDom(fallback);
      });
  }, []);

  const apply = async (patch: Partial<UserSettings>) => {
    const next = { ...settings!, ...patch } as UserSettings;
    try {
      const r = await updateSettings(next);
      setSettings(r.data);
      localStorage.setItem("fin_settings_cache", JSON.stringify(r.data));
      applyDom(r.data);
      if (r.data.locale) i18n.changeLanguage(r.data.locale);
    } catch {
      // offline fallback: keep settings locally
      setSettings(next);
      localStorage.setItem("fin_settings_cache", JSON.stringify(next));
      applyDom(next);
      if (next.locale) i18n.changeLanguage(next.locale);
    }
  };

  return <Ctx.Provider value={{ settings, apply }}>{children}</Ctx.Provider>;
}

function applyDom(s: UserSettings) {
  const theme = s.theme === "system"
    ? window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
    : s.theme;
  document.documentElement.setAttribute("data-theme", theme);
}

export function useTheme() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error("useTheme outside provider");
  return ctx;
}
