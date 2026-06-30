 import React, { useState } from "react";

export default function ViewComplaint({ user }) {
  const [id, setId] = useState("");
  const [data, setData] = useState(null);
  const [msg, setMsg] = useState("");

  async function view(e) {
    e.preventDefault();
    setMsg("");
    setData(null);

    if (!id.trim()) {
      setMsg("Please enter a complaint ID");
      return;
    }

    try {
      const url =
        "https://complaint-portal-tuyt.onrender.com/viewComplaint?id=" +
        encodeURIComponent(id.trim());

      const res = await fetch(url, {
        headers: { "X-Secret-Code": user.secret_code },
      });

      const body = await res.json();

      if (!res.ok) {
        setMsg("Error: " + (body.error || JSON.stringify(body)));
      } else {
        setData(body);
      }
    } catch (e) {
      setMsg("Network error: " + e.message);
    }
  }

  return (
    <div className="max-w-xl mx-auto mt-10 bg-white shadow-lg rounded-xl p-6">
      <h3 className="text-2xl font-semibold text-gray-800 mb-4">
        View Complaint by ID
      </h3>

      <form onSubmit={view} className="space-y-4">
        <input
          type="text"
          className="w-full px-4 py-2 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 outline-none"
          placeholder="Enter Complaint ID"
          value={id}
          onChange={(e) => setId(e.target.value)}
        />

        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition"
        >
          View Complaint
        </button>
      </form>
      {msg && (
        <div className="mt-4 p-3 bg-red-100 text-red-700 rounded-lg border border-red-300">
          {msg}
        </div>
      )}
      {data && (
        <div className="mt-6 bg-gray-50 p-4 rounded-lg border">
          <h4 className="text-lg font-bold text-gray-800">{data.title}</h4>

          <p className="text-gray-600 mt-1">
            <strong>Severity:</strong> {data.severity}
          </p>
          <p className="text-gray-700 mt-2">{data.summary}</p>

          <p className="mt-3 text-gray-600">
            <strong>Resolved:</strong>{" "}
            {data.resolved ? (
              <span className="text-green-600 font-semibold">Yes</span>
            ) : (
              <span className="text-red-600 font-semibold">No</span>
            )}
          </p>

          <p className="text-gray-600 mt-1">
            <strong>Owner ID:</strong> {data.user_id}
          </p>
        </div>
      )}
    </div>
  );
}
