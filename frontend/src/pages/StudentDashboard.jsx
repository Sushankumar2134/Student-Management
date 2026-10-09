import { useEffect, useState } from "react";

import api from "../services/api";
import ChangePassword from "./ChangePassword";

import "./StudentDashboard.css";


function StudentDashboard({ onLogout }) {

    const [student, setStudent] = useState(null);
    const [loading, setLoading] = useState(true);

    const [showEdit, setShowEdit] = useState(false);
    const [showChangePassword, setShowChangePassword] = useState(false);

    const [email, setEmail] = useState("");
    const [city, setCity] = useState("");

    const token = localStorage.getItem("token");


    // ==============================
    // GET STUDENT PROFILE
    // ==============================

    const fetchProfile = async () => {

        try {

            const response = await api.get("/profile", {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setStudent(response.data);

            setEmail(response.data.email || "");
            setCity(response.data.city || "");

        } catch (error) {

            console.error(error);

            if (error.response?.status === 401) {
                handleLogout();
            }

        } finally {

            setLoading(false);

        }
    };


    useEffect(() => {
        fetchProfile();
    }, []);


    // ==============================
    // LOGOUT
    // ==============================

    const handleLogout = () => {

        localStorage.removeItem("token");
        localStorage.removeItem("role");

        if (onLogout) {
            onLogout();
        }
    };


    // ==============================
    // UPDATE PROFILE
    // ==============================

    const handleUpdateProfile = async (e) => {

        e.preventDefault();

        try {

            await api.put(
                "/profile",
                {
                    email: email,
                    city: city,
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            alert("Profile updated successfully");

            setShowEdit(false);

            fetchProfile();

        } catch (error) {

            console.error(error);

            alert(
                error.response?.data?.message ||
                "Failed to update profile"
            );

        }
    };


    // ==============================
    // LOADING
    // ==============================

    if (loading) {

        return (
            <div className="student-loading">
                Loading profile...
            </div>
        );

    }


    // ==============================
    // NO PROFILE
    // ==============================

    if (!student) {

        return (
            <div className="student-loading">
                Unable to load profile.
            </div>
        );

    }


    // ==============================
    // CHANGE PASSWORD PAGE
    // ==============================

    if (showChangePassword) {

        return (
            <ChangePassword
                onBack={() => setShowChangePassword(false)}
            />
        );

    }


    // ==============================
    // DASHBOARD
    // ==============================

    return (

        <div className="student-dashboard">


            {/* ================= SIDEBAR ================= */}

            <aside className="student-sidebar">


                {/* LOGO */}

                <div className="student-logo">

                    <div className="logo-icon">
                        🎓
                    </div>

                    <div>

                        <h2>
                            StudentHub
                        </h2>

                        <span>
                            Student Portal
                        </span>

                    </div>

                </div>


                {/* MAIN MENU */}

                <div className="menu-section">

                    <p className="menu-title">
                        MAIN MENU
                    </p>


                    {/* Dashboard */}

                    <button
                        className="student-nav active"
                    >

                        <span>
                            🏠
                        </span>

                        Dashboard

                    </button>


                    {/* My Profile */}

                    <button
                        className="student-nav"
                        onClick={() => setShowEdit(true)}
                    >

                        <span>
                            👤
                        </span>

                        My Profile

                    </button>

                </div>


                {/* ACCOUNT MENU */}

                <div className="menu-section account-section">

                    <p className="menu-title">
                        ACCOUNT
                    </p>


                    {/* Change Password */}

                    <button
                        className="student-nav"
                        onClick={() =>
                            setShowChangePassword(true)
                        }
                    >

                        <span>
                            🔐
                        </span>

                        Change Password

                    </button>


                    {/* Logout */}

                    <button
                        className="student-nav logout-btn"
                        onClick={handleLogout}
                    >

                        <span>
                            🚪
                        </span>

                        Logout

                    </button>

                </div>

            </aside>



            {/* ================= MAIN ================= */}

            <main className="student-main">


                {/* TOPBAR */}

                <header className="student-topbar">

                    <div></div>


                    <div className="student-user">


                        <div className="user-avatar">

                            {student.name
                                ?.charAt(0)
                                .toUpperCase()
                            }

                        </div>


                        <span>
                            {student.name}
                        </span>


                        <button
                            onClick={handleLogout}
                        >
                            Logout
                        </button>

                    </div>

                </header>



                {/* CONTENT */}

                <section className="student-content">


                    {/* WELCOME */}

                    <div className="welcome-section">

                        <h1>
                            Welcome back, {student.name} 👋
                        </h1>

                        <p>
                            Here's your profile overview
                        </p>

                    </div>



                    {/* PROFILE CARD */}

                    <div className="student-profile-card">


                        {/* PROFILE HEADER */}

                        <div className="profile-header">


                            <div className="big-avatar">

                                {student.name
                                    ?.charAt(0)
                                    .toUpperCase()
                                }

                            </div>


                            <div>

                                <h2>
                                    {student.name}
                                </h2>

                                <p>
                                    Student ID: #{student.id}
                                </p>

                            </div>

                        </div>



                        {/* PROFILE DETAILS */}

                        <div className="profile-details">


                            {/* EMAIL */}

                            <div className="profile-item">

                                <span>
                                    Email
                                </span>

                                <strong>
                                    {student.email ||
                                        "Not provided"
                                    }
                                </strong>

                            </div>



                            {/* AGE */}

                            <div className="profile-item">

                                <span>
                                    Age
                                </span>

                                <strong>
                                    {student.age}
                                </strong>

                            </div>



                            {/* GENDER */}

                            <div className="profile-item">

                                <span>
                                    Gender
                                </span>

                                <strong>
                                    {student.gender ||
                                        "Not provided"
                                    }
                                </strong>

                            </div>



                            {/* CITY */}

                            <div className="profile-item">

                                <span>
                                    City
                                </span>

                                <strong>
                                    {student.city ||
                                        "Not provided"
                                    }
                                </strong>

                            </div>

                        </div>



                        {/* EDIT BUTTON */}

                        <div className="profile-footer">

                            <button
                                className="edit-profile-btn"
                                onClick={() =>
                                    setShowEdit(true)
                                }
                            >

                                ✏️ Edit Profile

                            </button>

                        </div>

                    </div>

                </section>

            </main>



            {/* ================= EDIT PROFILE MODAL ================= */}

            {showEdit && (

                <div className="student-modal-overlay">


                    <div className="student-modal">


                        {/* MODAL HEADER */}

                        <div className="modal-header">


                            <div>

                                <h2>
                                    Edit Profile
                                </h2>

                                <p>
                                    Update your profile information
                                </p>

                            </div>


                            <button
                                className="close-btn"
                                onClick={() =>
                                    setShowEdit(false)
                                }
                            >
                                ×
                            </button>

                        </div>



                        {/* FORM */}

                        <form
                            onSubmit={handleUpdateProfile}
                        >


                            {/* EMAIL */}

                            <div className="form-group">

                                <label>
                                    Email
                                </label>

                                <input
                                    type="email"
                                    value={email}
                                    onChange={(e) =>
                                        setEmail(e.target.value)
                                    }
                                    placeholder="Enter your email"
                                    required
                                />

                            </div>



                            {/* CITY */}

                            <div className="form-group">

                                <label>
                                    City
                                </label>

                                <input
                                    type="text"
                                    value={city}
                                    onChange={(e) =>
                                        setCity(e.target.value)
                                    }
                                    placeholder="Enter your city"
                                    required
                                />

                            </div>



                            {/* BUTTONS */}

                            <div className="modal-actions">


                                <button
                                    type="button"
                                    className="cancel-btn"
                                    onClick={() =>
                                        setShowEdit(false)
                                    }
                                >
                                    Cancel
                                </button>


                                <button
                                    type="submit"
                                    className="save-btn"
                                >
                                    Save Changes
                                </button>

                            </div>

                        </form>

                    </div>

                </div>

            )}

        </div>

    );
}


export default StudentDashboard;