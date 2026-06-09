import { useTranslation } from "react-i18next";
import { useTheme } from "../context/ThemeContext";
import { exportLocalBackup, restoreLocalBackup } from "../api/client";
import { useState } from "react";
import { Flash } from "../components/Flash";

export default function SettingsPage() {
  const { t } = useTranslation();
  const { settings, apply } = useTheme();
  const [restoreMsg, setRestoreMsg] = useState("");
  const [updateMsg, setUpdateMsg] = useState("");
  const [saveMsg, setSaveMsg] = useState("");

  return (
    <>
      <div className="page-header">
        <h1>{t("settings.title")}</h1>
      </div>
      <article className="card settings-grid">
        <label>
          {t("settings.theme")}
          <select
            value={settings.theme}
            onChange={async (e) => {
              await apply({ theme: e.target.value as "dark" | "light" | "system" });
              setSaveMsg("Сохранено");
            }}
          >
            <option value="dark">{t("settings.theme_dark")}</option>
            <option value="light">{t("settings.theme_light")}</option>
            <option value="system">{t("settings.theme_system")}</option>
          </select>
        </label>
        <label>
          {t("settings.locale")}
          <select
            value={settings.locale}
            onChange={async (e) => {
              await apply({ locale: e.target.value });
              setSaveMsg("Сохранено");
            }}
          >
            <option value="ru">Русский</option>
            <option value="en">English</option>
          </select>
        </label>
        <label>
          {t("settings.currency")}
          <select
            value={settings.base_currency}
            onChange={async (e) => {
              await apply({ base_currency: e.target.value });
              setSaveMsg("Сохранено");
            }}
          >
            <option value="RUB">RUB</option>
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
          </select>
        </label>
        <label>
          {t("settings.timezone")}
          <input
            value={settings.timezone}
            onChange={async (e) => {
              await apply({ timezone: e.target.value });
              setSaveMsg("Сохранено");
            }}
          />
        </label>
        {saveMsg && <Flash kind="success">{saveMsg}</Flash>}
        <hr style={{ margin: "1.25rem 0", border: "none", borderTop: "1px solid var(--border)" }} />
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
                setUpdateMsg(update ? t("settings.update_available") : t("settings.update_none"));
              } catch {
                setUpdateMsg(t("settings.update_unavailable"));
              }
            }}
          >
            {t("settings.update_check")}
          </button>
        </section>
        {restoreMsg && <Flash kind="info">{restoreMsg}</Flash>}
        {updateMsg && <Flash kind="info">{updateMsg}</Flash>}
      </article>
    </>
  );
}
