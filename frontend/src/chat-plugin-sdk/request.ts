export const LOBE_PLUGIN_SETTINGS = 'X-Lobe-Plugin-Settings';

export const createHeadersWithPluginSettings = (
  settings: any,
  header?: HeadersInit,
): HeadersInit => ({
  ...header,
  [LOBE_PLUGIN_SETTINGS]: typeof settings === 'string' ? settings : JSON.stringify(settings),
});
