import React, { useState } from "react";

export default function Register() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [masterKey, setMasterKey] = useState("");
  const [msg, setMsg] = useState("");

  async function submit(e) {
    e.preventDefault();
    setMsg("");

    const body = { name, email, is_admin: masterKey ? true : false };

    try {
      const res = await fetch("https://complaint-portal-tuyt.onrender.com/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(masterKey && { "X-Master-Key": masterKey }),
        },
        body: JSON.stringify(body),
      });

      const data = await res.json();

      if (!res.ok) {
        setMsg("Error: " + (data.error || JSON.stringify(data)));
      } else {
        setMsg("Registered successfully! Secret Code: " + data.secret_code);
        setName("");
        setEmail("");
        setMasterKey("");
      }
    } catch (err) {
      setMsg("Network error: " + err.message);
    }
  }

  return (
    <div className="bg-white shadow-lg rounded-xl p-6">
      <h3 className="text-xl font-semibold text-gray-800 mb-4">Register</h3>

      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium mb-1">
            Name
          </label>
          <input
            type="text"
            className="w-full px-4 py-2 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 outline-none"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block text-gray-700 font-medium mb-1">
            Email
          </label>
          <input
            type="email"
            className="w-full px-4 py-2 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 outline-none"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
         
        </div>
        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition"
        >
          Register
        </button>
      </form>
      {msg && (
        <div className="mt-4 p-3 bg-gray-100 text-gray-700 rounded-lg border">
          {msg}
        </div>
      )}
    </div>
  );
}
