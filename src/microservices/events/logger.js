export function logEvent(direction, topic, payload) {
  const arrow = direction === "out" ? "" : "";
  console.log(`${arrow}  [${topic}] ${JSON.stringify(payload)}`);
}
