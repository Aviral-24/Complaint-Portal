# Complaint Management System

This project is a simple  management system with a React frontend and a Go backend. It allows users to submit complaints, view their own complaints, and provides an admin interface to view and resolve all complaints.

## Features

*   **User Registration**: Users can register with a name and email.
*   **User Login**: Users can log in using a secret code.
*   **Complaint Submission**: Logged-in users can submit new complaints with a title, summary, and severity.
*   **View User Complaints**: Users can view a list of their submitted complaints.
*   **View Single Complaint**: Users can view the details of a specific complaint by its ID.
*   **Admin Dashboard**: Admins can view all complaints, including user details.
*   **Resolve Complaints**: Admins can mark complaints as resolved.

## Technologies Used

### Frontend
*   **React**: A JavaScript library for building user interfaces.
*   **Vite**: A fast build tool for modern web projects.
*   **Tailwind CSS**: A utility-first CSS framework for styling.

### Backend
*   **Go**: A statically typed, compiled programming language.
*   **Standard Library**: Utilizes Go's built-in `net/http` for API handling, `encoding/json` for JSON operations, and `crypto/rand` for secure random key generation.
*   **In-memory Store**: For simplicity, data is stored in memory (not persistent across server restarts).


### Prerequisites

*   Go (version 1.16 or higher)
*   npm or yarn

### Setup

2.  **Backend Setup:**
    cd backend
    go run main.go


3.  **Frontend Setup:**
    cd frontend
    npm install
    npm run dev


## Workflow Overview

This system uses a React frontend and a Go backend.
Users register/login with a secret code.
They can submit and view their complaints.
Admins can view and resolve all complaints.
Communication happens via RESTful APIs.
Data is stored in-memory for simplicity.

