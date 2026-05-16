export async function showNotification(title: string, body: string) {
  try {
    const w = window as Window & { __TAURI_INTERNALS__?: unknown };
    if (w.__TAURI_INTERNALS__) {
      const { isPermissionGranted, requestPermission, sendNotification } = await import(
        "@tauri-apps/plugin-notification"
      );
      let granted = await isPermissionGranted();
      if (!granted) {
        const p = await requestPermission();
        granted = p === "granted";
      }
      if (granted) {
        sendNotification({ title, body });
        return;
      }
    }
  } catch {
    /* fallback */
  }
  if (typeof Notification !== "undefined" && Notification.permission === "granted") {
    new Notification(title, { body });
  } else if (typeof Notification !== "undefined" && Notification.permission !== "denied") {
    const p = await Notification.requestPermission();
    if (p === "granted") new Notification(title, { body });
  }
}
