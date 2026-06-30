 const BASE = "http://localhost:8080";

async function apiPost(path, body, secret) {
  const headers = { "Content-Type": "application/json" };
  if (secret !== undefined && secret !== null) {
    headers["X-Secret-Code"] = secret;
  }

  const res = await fetch(BASE + path, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });

  return res;
}

async function apiGet(path, secret) {
  const headers = {};


  if (secret !== undefined && secret !== null) {
    headers["X-Secret-Code"] = secret;
  }

  const res = await fetch(BASE + path, { headers });
  return res;
}

export { apiPost, apiGet };
