import React, { useState } from "react";

export default function Login({ onLogin }) {
  const [secret, setSecret] = useState("");
  const [msg, setMsg] = useState("");

  async function submit(e) {
    e.preventDefault();
    setMsg("");

    try {
      const res = await fetch("https://complaint-portal-tuyt.onrender.com/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ secret_code: secret }),
      });

      const data = await res.json();

      if (!res.ok) {
        setMsg(" " + (data.error || "Login failed"));
      } else {
        setMsg(" Logged in as " + data.name);
        onLogin(data);
        setSecret("");
      }
    } catch (err) {
      setMsg(" Network error: " + err.message);
    }
  }

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100 px-4">
      <div className="w-full max-w-md bg-white shadow-xl rounded-2xl p-8">
        <h3 className="text-2xl font-semibold text-gray-800 mb-6 text-center">
          Login
        </h3>

        <form onSubmit={submit} className="space-y-5">
          <div>
            <label className="block text-gray-700 mb-1 font-medium">
              Secret Code
            </label>
            <input
              type="text"
              value={secret}
              onChange={(e) => setSecret(e.target.value)}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg 
                         focus:outline-none focus:ring-2 focus:ring-blue-500 
                         text-gray-700"
              placeholder="Enter your secret code"
            />
          </div>

          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded-lg 
                       font-medium text-lg
                       hover:bg-blue-700 transition-all 
                       active:scale-[0.98]"
          >
            Login
          </button>
        </form>

        {msg && (
          <p
            className={`mt-4 text-center text-sm font-medium ${
              msg.includes("NO") ? "text-red-600" : "text-green-600"
            }`}
          >
            {msg}
          </p>
        )}
      </div>
    </div>
  );
}
