 import React, { useEffect, useState } from "react";

export default function ComplaintsList({ user }) {
  const [list, setList] = useState(() => {
    const raw = localStorage.getItem(`complaints_${user.secret_code}`);
    return raw ? JSON.parse(raw) : [];
  });
  const [err, setErr] = useState("");

  async function load() {
    setErr("");
    try {
      const res = await fetch("https://complaint-portal-tuyt.onrender.com/getAllComplaintsForUser", {
        headers: { "X-Secret-Code": user.secret_code }
      });
      const data = await res.json();
      if (!res.ok) {
        setErr(data.error || "Something went wrong");
      } else {
        setList(data);
      }
    } catch (e) {
      setErr("Network error: " + e.message);
    }
  }

  useEffect(() => { load(); }, [user.secret_code]);

  useEffect(() => {
    localStorage.setItem(`complaints_${user.secret_code}`, JSON.stringify(list));
  }, [list, user.secret_code]);

  return (
    <div className="max-w-2xl mx-auto mt-6 bg-white shadow-lg rounded-xl p-6">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-xl font-semibold text-gray-800">
          Your Complaints
        </h3>

        <button
          onClick={load}
          className="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition"
        >
          Refresh
        </button>
      </div>
      {err && (
        <div className="mb-3 p-3 rounded-lg bg-red-100 text-red-600 border border-red-300">
          {err}
        </div>
      )}
      <ul className="space-y-3">
        {list.map((c) => (
          <li
            key={c.id}
            className="p-4 bg-gray-50 border rounded-lg hover:shadow-md transition"
          >
            <div className="text-lg font-medium text-gray-800">
              {c.title}
            </div>
            <div className="text-sm text-gray-500">Complaint ID: {c.id}</div>
          </li>
        ))}
      </ul>

      {list.length === 0 && !err && (
        <p className="mt-4 text-center text-gray-500">No complaints found.</p>
      )}
    </div>
  );
}
