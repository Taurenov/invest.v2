import { useTranslation } from "react-i18next";
import { useTheme } from "../context/ThemeContext";
import { exportLocalBackup, restoreLocalBackup } from "../api/client";
import { useState } from "react";

export default function SettingsPage() {
  const { t } = useTranslation();
  const { settings, apply } = useTheme();
  const [restoreMsg, setRestoreMsg] = useState("");
  const [updateMsg, setUpdateMsg] = useState("");

  if (!settings) return <p>{t("common.loading")}</p>;

  return (
    <>
      <h1>{t("settings.title")}</h1>
      <article className="card">
        <label>
          {t("settings.theme")}
          <select
            value={settings.theme}
            onChange={(e) => apply({ theme: e.target.value as "dark" | "light" | "system" })}
          >
            <option value="dark">{t("settings.theme_dark")}</option>
            <option value="light">{t("settings.theme_light")}</option>
            <option value="system">{t("settings.theme_system")}</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.locale")}
          <select value={settings.locale} onChange={(e) => apply({ locale: e.target.value })}>
            <option value="ru">Русский</option>
            <option value="en">English</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.currency")}
          <select
            value={settings.base_currency}
            onChange={(e) => apply({ base_currency: e.target.value })}
          >
            <option value="RUB">RUB</option>
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
          </select>
        </label>
        <label style={{ display: "block", marginTop: "1rem" }}>
          {t("settings.timezone")}
          <input
            value={settings.timezone}
            onChange={(e) => apply({ timezone: e.target.value })}
          />
        </label>
        <hr style={{ margin: "1rem 0", opacity: 0.2 }} />
        <section className="toolbar">
          <button type="button" onClick={exportLocalBackup}>
            {t("settings.backup_export")}
          </button>
          <label className="secondary btn" style={{ cursor: "pointer" }}>
            {t("settings.backup_import")}
            <input
              type="file"
              accept="application/json"
              style={{ display: "none" }}
              onChange={async (e) => {
                const f = e.target.files?.[0];
                if (!f) return;
                try {
                  await restoreLocalBackup(f);
                  setRestoreMsg(t("settings.backup_ok"));
                } catch {
                  setRestoreMsg(t("settings.backup_fail"));
                }
              }}
            />
          </label>
          <button
            type="button"
            className="secondary"
            onClick={async () => {
              try {
                const mod = await import("@tauri-apps/plugin-updater");
                const update = await mod.check();
                if (update) {
                  setUpdateMsg(t("settings.update_available"));
                } else {
                  setUpdateMsg(t("settings.update_none"));
                }
              } catch {
                setUpdateMsg(t("settings.update_unavailable"));
              }
            }}
          >
            {t("settings.update_check")}
          </button>
        </section>
        {restoreMsg && <p className="label" style={{ marginTop: "0.5rem" }}>{restoreMsg}</p>}
        {updateMsg && <p className="label" style={{ marginTop: "0.25rem" }}>{updateMsg}</p>}
      </article>
    </>
  );
}
