import React, { useState } from "react";
import Register from "./components/Register";
import Login from "./components/Login";
import SubmitComplaint from "./components/SubmitComplaint";
import ComplaintsList from "./components/ComplaintsList";
import AdminComplaints from "./components/AdminComplaints";
import ViewComplaint from "./components/ViewComplaint";

export default function App() {
  const [user, setUser] = useState(() => {
    const raw = localStorage.getItem("user");
    return raw ? JSON.parse(raw) : null;
  });

  function onLogin(u) {
    setUser(u);
    localStorage.setItem("user", JSON.stringify(u));
  }

  function onLogout() {
    setUser(null);
    localStorage.removeItem("user");
  }

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">

      <header className="bg-blue-600 text-white px-6 py-4 shadow-lg flex justify-between items-center">
        <h1 className="text-2xl font-bold bg-black">Complaint Portal</h1>

        <div className="flex items-center gap-4">
          {user ? (
            <>
              <span className="text-white font-medium">
                Hi, {user.name} {user.is_admin ? "(Admin)" : ""}
              </span>
              <button
                onClick={onLogout}
                className="px-4 py-2 bg-red-500 hover:bg-red-600 rounded-lg text-white shadow"
              >
                Logout
              </button>
            </>
          ) : (
            <span className="text-white">Please login or register</span>
          )}
        </div>
      </header>


      <main className="flex-1 container mx-auto px-4 py-6">

 
        {!user && (
          <section className="grid md:grid-cols-2 gap-6">
            <Register />
            <Login onLogin={onLogin} />
          </section>
        )}

    
        {user && (
          <section className="space-y-8">

    
            <SubmitComplaint user={user} />

            <ComplaintsList user={user} />


            {user.is_admin && <AdminComplaints user={user} />}


            <ViewComplaint user={user} />
          </section>
        )}
      </main>

      <footer className="bg-gray-800 text-white text-center py-4 text-sm">
        &copy; 2025 Complaint Portal. All rights reserved.
      </footer>

    </div>
  );
}
