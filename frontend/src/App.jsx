import { useState } from "react";

import Login from "./pages/Login";
import Register from "./pages/Register";
import AdminDashboard from "./pages/AdminDashboard";
import StudentDashboard from "./pages/StudentDashboard";

function App() {

    const [role, setRole] = useState(
        localStorage.getItem("role")
    );

    const handleLogin = (userRole) => {
        setRole(userRole);
    };

    const handleLogout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("role");
        setRole(null);
    };

    // Show registration page
    if (role === "register") {
        return <Register onLogin={handleLogin} />;
    }

    // No login
    if (!role) {
        return <Login onLogin={handleLogin} />;
    }

    // Admin
    if (role === "admin") {
        return <AdminDashboard onLogout={handleLogout} />;
    }

    // Student
    if (role === "student") {
        return <StudentDashboard onLogout={handleLogout} />;
    }

    return <Login onLogin={handleLogin} />;
}

export default App;