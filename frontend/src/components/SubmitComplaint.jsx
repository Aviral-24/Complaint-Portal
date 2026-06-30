 import React, { useState } from "react";

export default function SubmitComplaint({ user }) {
  const [title, setTitle] = useState("");
  const [summary, setSummary] = useState("");
  const [severity, setSeverity] = useState(3);
  const [msg, setMsg] = useState("");

  async function submit(e) {
    e.preventDefault();
    setMsg("");

    try {
      const res = await fetch("http://localhost:8080/submitComplaint", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Secret-Code": user.secret_code,
        },
        body: JSON.stringify({ title, summary, severity }),
      });

      const data = await res.json();

      if (!res.ok) {
        setMsg("NO " + (data.error || "Something went wrong"));
      } else {
        setMsg(data.id);
        setTitle("");
        setSummary("");
        setSeverity(3);
      }
    } catch (err) {
      setMsg(" Network error: " + err.message);
    }
  }

  return (
    <div className="flex justify-center items-center w-full px-4 py-6">
      <div className="w-full max-w-xl bg-white shadow-xl rounded-2xl p-8">
        <h3 className="text-2xl font-semibold text-gray-800 mb-6 text-center">
          Submit Complaint
        </h3>

        <form onSubmit={submit} className="space-y-5">
   <div>
            <label className="block text-gray-700 font-medium mb-1">
              Title
            </label>
            <input
              className="w-full px-4 py-2 border border-gray-300 rounded-lg 
                         text-gray-700 focus:outline-none focus:ring-2 
                         focus:ring-blue-500"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Complaint title"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 font-medium mb-1">
              Summary
            </label>
            <textarea
              className="w-full px-4 py-2 border border-gray-300 rounded-lg 
                         text-gray-700 focus:outline-none focus:ring-2 
                         focus:ring-blue-500"
              rows="4"
              value={summary}
              placeholder="Describe the issue..."
              onChange={(e) => setSummary(e.target.value)}
              required
            ></textarea>
          </div>
          <div>
            <label className="block text-gray-700 font-medium mb-1">
              Severity
            </label>
            <select
              className="w-full px-4 py-2 border border-gray-300 rounded-lg 
                         text-gray-700 focus:outline-none focus:ring-2 
                         focus:ring-blue-500"
              value={severity}
              onChange={(e) => setSeverity(parseInt(e.target.value))}
            >
              <option value={1}>1 - Low</option>
              <option value={2}>2</option>
              <option value={3}>3 - Medium</option>
              <option value={4}>4</option>
              <option value={5}>5 - Critical</option>
            </select>
          </div>
          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded-lg 
                       text-lg font-medium hover:bg-blue-700 
                       transition-all active:scale-[0.98]"
          >
            Submit Complaint
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
