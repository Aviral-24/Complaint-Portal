import React, { useEffect, useState } from "react";

export default function AdminComplaints({ user }) {
  const [list, setList] = useState([]);
  const [msg, setMsg] = useState("");

  async function load() {
    setMsg("");
    try {
      const res = await fetch("https://complaint-portal-tuyt.onrender.com/getAllComplaintsForAdmin", {
        headers: { "X-Secret-Code": user.secret_code }
      });
      const data = await res.json();
      if (!res.ok) setMsg("Error: " + (data.error || JSON.stringify(data)));
      else setList(data);
    } catch (e) {
      setMsg("Network error: " + e.message);
    }
  }

  useEffect(() => { load(); }, []);

  async function resolve(id) {
    setMsg("");
    try {
      const res = await fetch("http://localhost:8080/resolveComplaint", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Secret-Code": user.secret_code
        },
        body: JSON.stringify({ id })
      });
      const data = await res.json();
      if (!res.ok) setMsg("Error: " + (data.error || JSON.stringify(data)));
      else {
        setMsg(`Complaint #${id} marked as resolved`);
        load();
      }
    } catch (e) {
      setMsg("Network error: " + e.message);
    }
  }

  return (
    <div className="max-w-5xl mx-auto mt-8 bg-white shadow-lg rounded-xl p-6">
      <div className="flex justify-between items-center mb-5">
        <h3 className="text-2xl font-semibold text-gray-800">
          Admin: All Complaints
        </h3>

        <button
          onClick={load}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
        >
          Refresh
        </button>
      </div>


      {msg && (
        <div className="mb-4 p-3 bg-blue-100 text-blue-700 rounded-lg border border-blue-300">
          {msg}
        </div>
      )}


      <div className="overflow-x-auto mt-4">
        <table className="w-full table-auto border-collapse">
          <thead>
            <tr className="bg-gray-100 text-left text-gray-700">
              <th className="p-3 font-medium">Title</th>
              <th className="p-3 font-medium">User</th>
              <th className="p-3 font-medium">Resolved</th>
              <th className="p-3 font-medium text-center">Action</th>
            </tr>
          </thead>

          <tbody>
            {list.map((c) => (
              <tr
                key={c.id}
                className="border-b hover:bg-gray-50 transition"
              >
                <td className="p-3">{c.title}</td>
                <td className="p-3">{c.user_name}</td>
                <td className="p-3">
                  {c.resolved ? (
                    <span className="px-3 py-1 bg-green-100 text-green-700 rounded-lg text-sm">
                      Yes
                    </span>
                  ) : (
                    <span className="px-3 py-1 bg-red-100 text-red-700 rounded-lg text-sm">
                      No
                    </span>
                  )}
                </td>

                <td className="p-3 text-center">
                  {!c.resolved && (
                    <button
                      onClick={() => resolve(c.id)}
                      className="px-3 py-1.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition text-sm"
                    >
                      Resolve
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {list.length === 0 && (
          <p className="text-center text-gray-500 mt-4">
            No complaints found.
          </p>
        )}
      </div>
    </div>
  );
}
