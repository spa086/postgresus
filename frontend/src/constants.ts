export function getApplicationServer() {
  const origin = window.location.origin;
  const url = new URL(origin);
  return `${url.protocol}//${url.hostname}:4005`;
}
