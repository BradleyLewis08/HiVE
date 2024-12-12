export function getFullURL(url: string) {
  if (!url) return "";
  if (url.startsWith("http")) return url;
  return `http://${url}`; // Use HTTP since we don't know if HTTPS is configured
}

export function formatDate(dateString: string) {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}
