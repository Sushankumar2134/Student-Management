import { useEffect, useState } from "react";
import api from "../services/api";
import "./AdminDashboard.css";
import ChangePassword from "./ChangePassword";
function AdminDashboard({ onLogout }) {

    const [students, setStudents] = useState([]);
    const [search, setSearch] = useState("");

    const [loading, setLoading] = useState(true);

    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [showChangePassword, setShowChangePassword] = useState(false);
    const [selectedStudent, setSelectedStudent] = useState(null);

    const [statistics, setStatistics] = useState({
        total_students: 0,
        male_students: 0,
        female_students: 0
    });

    const [formData, setFormData] = useState({
        name: "",
        email: "",
        age: "",
        gender: "",
        city: ""
    });

    const [message, setMessage] = useState("");
    const [error, setError] = useState("");

    const token = localStorage.getItem("token");

    // ==========================================
    // GET STUDENTS
    // ==========================================

    const getStudents = async (searchText = "") => {

        try {

            const response = await api.get("/students", {
                params: {
                    search: searchText
                },
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            setStudents(response.data || []);

        } catch (error) {

            console.error("Get students error:", error);

            setError("Unable to load students.");

        } finally {

            setLoading(false);

        }
    };


    // ==========================================
    // GET STATISTICS
    // ==========================================

    const getStatistics = async () => {

        try {

            const response = await api.get("/statistics", {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            setStatistics(response.data);

        } catch (error) {

            console.error("Statistics error:", error);

        }
    };


    // ==========================================
    // INITIAL LOAD
    // ==========================================

    useEffect(() => {

        getStudents();
        getStatistics();

    }, []);


    // ==========================================
    // SEARCH
    // ==========================================

    const handleSearch = (e) => {

        const value = e.target.value;

        setSearch(value);

        getStudents(value);
    };


    // ==========================================
    // FORM INPUT
    // ==========================================

    const handleInputChange = (e) => {

        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });

    };


    // ==========================================
    // ADD STUDENT
    // ==========================================

    const handleAddStudent = async (e) => {

        e.preventDefault();

        setMessage("");
        setError("");

        try {

            await api.post(
                "/students",
                {
                    name: formData.name,
                    email: formData.email,
                    age: Number(formData.age),
                    gender: formData.gender,
                    city: formData.city
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            setMessage("Student added successfully.");

            setShowAddModal(false);

            resetForm();

            getStudents(search);
            getStatistics();

        } catch (error) {

            console.error("Add student error:", error);

            setError(
                error.response?.data?.message ||
                "Failed to add student."
            );

        }

    };


    // ==========================================
    // EDIT STUDENT
    // ==========================================

    const openEditModal = (student) => {

        setSelectedStudent(student);

        setFormData({
            name: student.name || "",
            email: student.email || "",
            age: student.age || "",
            gender: student.gender || "",
            city: student.city || ""
        });

        setShowEditModal(true);
    };


    const handleEditStudent = async (e) => {

        e.preventDefault();

        setMessage("");
        setError("");

        try {

            await api.put(
                `/students/${selectedStudent.id}`,
                {
                    name: formData.name,
                    email: formData.email,
                    age: Number(formData.age),
                    gender: formData.gender,
                    city: formData.city
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            setMessage("Student updated successfully.");

            setShowEditModal(false);

            resetForm();

            getStudents(search);
            getStatistics();

        } catch (error) {

            console.error("Update student error:", error);

            setError(
                error.response?.data?.message ||
                "Failed to update student."
            );

        }

    };


    // ==========================================
    // VIEW STUDENT
    // ==========================================

    const openViewModal = async (student) => {

        try {

            const response = await api.get(
                `/students/${student.id}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            setSelectedStudent(response.data);

            setShowViewModal(true);

        } catch (error) {

            console.error("View student error:", error);

            setError("Unable to get student details.");

        }

    };


    // ==========================================
    // DELETE STUDENT
    // ==========================================

    const openDeleteModal = (student) => {

        setSelectedStudent(student);

        setShowDeleteModal(true);
    };


    const handleDeleteStudent = async () => {

        setMessage("");
        setError("");

        try {

            await api.delete(
                `/students/${selectedStudent.id}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            setMessage("Student deleted successfully.");

            setShowDeleteModal(false);

            setSelectedStudent(null);

            getStudents(search);
            getStatistics();

        } catch (error) {

            console.error("Delete student error:", error);

            setError(
                error.response?.data?.message ||
                "Failed to delete student."
            );

        }

    };


    // ==========================================
    // RESET FORM
    // ==========================================

    const resetForm = () => {

        setFormData({
            name: "",
            email: "",
            age: "",
            gender: "",
            city: ""
        });

        setSelectedStudent(null);
    };


    // ==========================================
    // LOGOUT
    // ==========================================

    const logout = () => {

        localStorage.removeItem("token");
        localStorage.removeItem("role");

        onLogout();
    };

    if (showChangePassword) {
    return (
        <ChangePassword
            onBack={() => setShowChangePassword(false)}
        />
    );
}
    return (

        <div className="admin-layout">

            {/* =====================================
                SIDEBAR
            ===================================== */}

            <aside className="sidebar">

                <div className="sidebar-logo">

                    <div className="sidebar-logo-icon">
                        S
                    </div>

                    <div>
                        <h2>StudentHub</h2>
                        <span>Management System</span>
                    </div>

                </div>


                <nav className="sidebar-nav">

                    <div className="nav-title">
                        MAIN MENU
                    </div>

                    <button className="nav-item active">
                        <span>▣</span>
                        Dashboard
                    </button>
                    <div className="nav-title second">
                        ACCOUNT
                    </div>

                    <button className="nav-item">
                        <span>⚙</span>
                        Settings
                    </button>

                   <button
                    className="nav-item"
                    onClick={() => setShowChangePassword(true)}
                >
                    <span>🔐</span>
                         Change Password
                    </button>

                </nav>


                <div className="sidebar-bottom">

                    <button
                        className="logout-sidebar"
                        onClick={logout}
                    >
                        <span>↪</span>
                        Logout
                    </button>

                </div>

            </aside>


            {/* =====================================
                MAIN AREA
            ===================================== */}

            <main className="admin-main">

                {/* TOP BAR */}

                <header className="topbar">

                    <div>

                        <h1>
                            Dashboard
                        </h1>

                        <p>
                            Welcome back, Admin 👋
                        </p>

                    </div>


                    <div className="topbar-right">

                        <button className="notification-button">
                            🔔
                            <span className="notification-dot"></span>
                        </button>


                        <div className="admin-profile">

                            <div className="admin-avatar">
                                A
                            </div>

                            <div>
                                <strong>
                                    Admin
                                </strong>

                                <span>
                                    Administrator
                                </span>
                            </div>

                        </div>

                    </div>

                </header>


                {/* CONTENT */}

                <div className="dashboard-content">


                    {/* SUCCESS / ERROR */}

                    {message && (

                        <div className="dashboard-success">
                            ✓ {message}
                        </div>

                    )}

                    {error && (

                        <div className="dashboard-error">
                            {error}
                        </div>

                    )}


                    {/* =================================
                        STATISTICS
                    ================================= */}

                    <div className="stats-grid">


                        <div className="stat-card blue">

                            <div className="stat-icon">
                                ♙
                            </div>

                            <div className="stat-info">

                                <span>
                                    Total Students
                                </span>

                                <strong>
                                    {statistics.total_students}
                                </strong>

                                <small>
                                    Registered students
                                </small>

                            </div>

                        </div>


                        <div className="stat-card green">

                            <div className="stat-icon">
                                ♂
                            </div>

                            <div className="stat-info">

                                <span>
                                    Male Students
                                </span>

                                <strong>
                                    {statistics.male_students}
                                </strong>

                                <small>
                                    Male students
                                </small>

                            </div>

                        </div>


                        <div className="stat-card purple">

                            <div className="stat-icon">
                                ♀
                            </div>

                            <div className="stat-info">

                                <span>
                                    Female Students
                                </span>

                                <strong>
                                    {statistics.female_students}
                                </strong>

                                <small>
                                    Female students
                                </small>

                            </div>

                        </div>

                    </div>


                    {/* =================================
                        STUDENTS SECTION
                    ================================= */}

                    <div className="students-card">


                        <div className="students-header">

                            <div>

                                <h2>
                                    Students
                                </h2>

                                <p>
                                    Manage all registered students
                                </p>

                            </div>


                            <button
                                className="add-student-button"
                                onClick={() => {

                                    resetForm();

                                    setShowAddModal(true);

                                }}
                            >
                                <span>+</span>
                                Add Student
                            </button>

                        </div>


                        {/* SEARCH */}

                        <div className="search-container">

                            <div className="search-box">

                                <span>
                                    🔍
                                </span>

                                <input
                                    type="text"
                                    placeholder="Search by name, email or city..."
                                    value={search}
                                    onChange={handleSearch}
                                />

                                {search && (

                                    <button
                                        onClick={() => {

                                            setSearch("");

                                            getStudents("");

                                        }}
                                    >
                                        ×
                                    </button>

                                )}

                            </div>

                        </div>


                        {/* TABLE */}

                        <div className="table-container">

                            {loading ? (

                                <div className="table-loading">
                                    Loading students...
                                </div>

                            ) : students.length === 0 ? (

                                <div className="empty-state">

                                    <div>
                                        📚
                                    </div>

                                    <h3>
                                        No students found
                                    </h3>

                                    <p>
                                        Try another search or add a new student.
                                    </p>

                                </div>

                            ) : (

                                <table>

                                    <thead>

                                        <tr>

                                            <th>
                                                ID
                                            </th>

                                            <th>
                                                STUDENT
                                            </th>

                                            <th>
                                                EMAIL
                                            </th>

                                            <th>
                                                AGE
                                            </th>

                                            <th>
                                                GENDER
                                            </th>

                                            <th>
                                                CITY
                                            </th>

                                            <th>
                                                ACTIONS
                                            </th>

                                        </tr>

                                    </thead>


                                    <tbody>

                                        {students.map((student) => (

                                            <tr key={student.id}>

                                                <td>
                                                    <span className="student-id">
                                                        #{student.id}
                                                    </span>
                                                </td>


                                                <td>

                                                    <div className="student-name">

                                                        <div className="student-avatar">
                                                            {student.name
                                                                ?.charAt(0)
                                                                .toUpperCase()}
                                                        </div>

                                                        <strong>
                                                            {student.name}
                                                        </strong>

                                                    </div>

                                                </td>


                                                <td>
                                                    <span className="email-text">
                                                        {student.email}
                                                    </span>
                                                </td>


                                                <td>
                                                    {student.age}
                                                </td>


                                                <td>

                                                    <span
                                                        className={
                                                            student.gender?.toLowerCase() === "male"
                                                                ? "gender-badge male"
                                                                : "gender-badge female"
                                                        }
                                                    >
                                                        {student.gender}
                                                    </span>

                                                </td>


                                                <td>
                                                    {student.city}
                                                </td>


                                                <td>

                                                    <div className="action-buttons">

                                                        <button
                                                            className="view-action"
                                                            title="View"
                                                            onClick={() =>
                                                                openViewModal(student)
                                                            }
                                                        >
                                                            👁
                                                        </button>

                                                        <button
                                                            className="edit-action"
                                                            title="Edit"
                                                            onClick={() =>
                                                                openEditModal(student)
                                                            }
                                                        >
                                                            ✏
                                                        </button>

                                                        <button
                                                            className="delete-action"
                                                            title="Delete"
                                                            onClick={() =>
                                                                openDeleteModal(student)
                                                            }
                                                        >
                                                            🗑
                                                        </button>

                                                    </div>

                                                </td>

                                            </tr>

                                        ))}

                                    </tbody>

                                </table>

                            )}

                        </div>

                    </div>

                </div>

            </main>


            {/* =====================================
                ADD / EDIT MODAL
            ===================================== */}

            {(showAddModal || showEditModal) && (

                <div className="modal-overlay">

                    <div className="student-modal">

                        <button
                            className="modal-close"
                            onClick={() => {

                                setShowAddModal(false);
                                setShowEditModal(false);
                                resetForm();

                            }}
                        >
                            ×
                        </button>


                        <div className="modal-icon">
                            {showEditModal ? "✏" : "+"}
                        </div>


                        <h2>
                            {showEditModal
                                ? "Edit Student"
                                : "Add New Student"}
                        </h2>

                        <p>
                            {showEditModal
                                ? "Update student information below."
                                : "Enter the student's information below."}
                        </p>


                        <form
                            onSubmit={
                                showEditModal
                                    ? handleEditStudent
                                    : handleAddStudent
                            }
                        >

                            <div className="form-grid">

                                <div className="admin-form-group">

                                    <label>
                                        Full Name
                                    </label>

                                    <input
                                        type="text"
                                        name="name"
                                        placeholder="Enter full name"
                                        value={formData.name}
                                        onChange={handleInputChange}
                                        required
                                    />

                                </div>


                                <div className="admin-form-group">

                                    <label>
                                        Email
                                    </label>

                                    <input
                                        type="email"
                                        name="email"
                                        placeholder="Enter email"
                                        value={formData.email}
                                        onChange={handleInputChange}
                                    />

                                </div>


                                <div className="admin-form-group">

                                    <label>
                                        Age
                                    </label>

                                    <input
                                        type="number"
                                        name="age"
                                        placeholder="Enter age"
                                        value={formData.age}
                                        onChange={handleInputChange}
                                        required
                                    />

                                </div>


                                <div className="admin-form-group">

                                    <label>
                                        Gender
                                    </label>

                                    <select
                                        name="gender"
                                        value={formData.gender}
                                        onChange={handleInputChange}
                                        required
                                    >

                                        <option value="">
                                            Select gender
                                        </option>

                                        <option value="Male">
                                            Male
                                        </option>

                                        <option value="Female">
                                            Female
                                        </option>

                                    </select>

                                </div>

                            </div>


                            <div className="admin-form-group">

                                <label>
                                    City
                                </label>

                                <input
                                    type="text"
                                    name="city"
                                    placeholder="Enter city"
                                    value={formData.city}
                                    onChange={handleInputChange}
                                />

                            </div>


                            <div className="modal-actions">

                                <button
                                    type="button"
                                    className="cancel-button"
                                    onClick={() => {

                                        setShowAddModal(false);
                                        setShowEditModal(false);
                                        resetForm();

                                    }}
                                >
                                    Cancel
                                </button>


                                <button
                                    type="submit"
                                    className="save-button"
                                >
                                    {showEditModal
                                        ? "Update Student"
                                        : "Add Student"}
                                </button>

                            </div>

                        </form>

                    </div>

                </div>

            )}


            {/* =====================================
                VIEW MODAL
            ===================================== */}

            {showViewModal && selectedStudent && (

                <div className="modal-overlay">

                    <div className="view-modal">

                        <button
                            className="modal-close"
                            onClick={() => {

                                setShowViewModal(false);
                                setSelectedStudent(null);

                            }}
                        >
                            ×
                        </button>


                        <div className="view-profile-header">

                            <div className="large-avatar">
                                {selectedStudent.name
                                    ?.charAt(0)
                                    .toUpperCase()}
                            </div>

                            <h2>
                                {selectedStudent.name}
                            </h2>

                            <span>
                                Student ID #{selectedStudent.id}
                            </span>

                        </div>


                        <div className="details-grid">

                            <div className="detail-item">
                                <span>Email</span>
                                <strong>
                                    {selectedStudent.email}
                                </strong>
                            </div>

                            <div className="detail-item">
                                <span>Age</span>
                                <strong>
                                    {selectedStudent.age}
                                </strong>
                            </div>

                            <div className="detail-item">
                                <span>Gender</span>
                                <strong>
                                    {selectedStudent.gender}
                                </strong>
                            </div>

                            <div className="detail-item">
                                <span>City</span>
                                <strong>
                                    {selectedStudent.city}
                                </strong>
                            </div>

                        </div>


                        <button
                            className="close-details-button"
                            onClick={() => {

                                setShowViewModal(false);
                                setSelectedStudent(null);

                            }}
                        >
                            Close
                        </button>

                    </div>

                </div>

            )}


            {/* =====================================
                DELETE MODAL
            ===================================== */}

            {showDeleteModal && selectedStudent && (

                <div className="modal-overlay">

                    <div className="delete-modal">

                        <div className="delete-icon">
                            🗑
                        </div>

                        <h2>
                            Delete Student?
                        </h2>

                        <p>
                            Are you sure you want to delete
                            <strong>
                                {" "}{selectedStudent.name}
                            </strong>?
                            This action cannot be undone.
                        </p>


                        <div className="modal-actions">

                            <button
                                className="cancel-button"
                                onClick={() =>
                                    setShowDeleteModal(false)
                                }
                            >
                                Cancel
                            </button>

                            <button
                                className="delete-confirm-button"
                                onClick={handleDeleteStudent}
                            >
                                Delete Student
                            </button>

                        </div>

                    </div>

                </div>

            )}

        </div>
    );
}

export default AdminDashboard;