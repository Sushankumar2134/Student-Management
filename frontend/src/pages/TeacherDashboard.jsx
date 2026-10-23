import { useEffect, useState } from "react";
import api from "../services/api";
import ChangePassword from "./ChangePassword";
import "./TeacherDashboard.css";

function TeacherDashboard({ onLogout }) {
    const [teacher, setTeacher] = useState(null);
    const [stats, setStats] = useState({
        assigned_classes: 0,
        assigned_subjects: 0,
        assigned_students: 0,
        recent_attendance: 0
    });
    const [classes, setClasses] = useState([]);
    const [subjects, setSubjects] = useState([]);
    const [students, setStudents] = useState([]);
    const [attendance, setAttendance] = useState([]);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("dashboard");
    const [selectedClassId, setSelectedClassId] = useState("");
    const [selectedSubjectId, setSelectedSubjectId] = useState("");
    const [attendanceDate, setAttendanceDate] = useState("");
    const [showChangePassword, setShowChangePassword] = useState(false);
    const [showEditProfile, setShowEditProfile] = useState(false);
    const [editForm, setEditForm] = useState({
        email: "",
        phone: "",
        address: "",
        subject: ""
    });
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [attendanceLoading, setAttendanceLoading] = useState(false);

    const token = localStorage.getItem("token");

    useEffect(() => {
        fetchDashboard();
        fetchClasses();
    }, []);

    useEffect(() => {
        if (selectedClassId) {
            fetchSubjects(selectedClassId);
        }
    }, [selectedClassId]);

    const fetchDashboard = async () => {
        try {
            const dashRes = await api.get("/teacher/dashboard", {
    headers: {
        Authorization: `Bearer ${token}`
    }
});
            setTeacher(dashRes.data.teacher);
            setStats(dashRes.data.statistics);
            setEditForm({
                email: dashRes.data.teacher.email || "",
                phone: dashRes.data.teacher.phone || "",
                address: dashRes.data.teacher.address || "",
                subject: dashRes.data.teacher.subject || ""
            });
        } catch (err) {
            console.error("Dashboard error:", err);
            if (err.response?.status === 401) handleLogout();
        } finally {
            setLoading(false);
        }
    };

    const fetchClasses = async () => {
        try {
            const res = await api.get("/teacher/classes", {
                headers: { Authorization: `Bearer ${token}` }
            });
            setClasses(res.data);
            if (res.data.length > 0 && !selectedClassId) {
                setSelectedClassId(res.data[0].id);
            }
        } catch (err) {
            console.error("Classes error:", err);
        }
    };

    const fetchSubjects = async (classId) => {
        try {
            const res = await api.get("/teacher/subjects", {
                params: { class_id: classId },
                headers: { Authorization: `Bearer ${token}` }
            });
            setSubjects(res.data);
            if (res.data.length > 0) {
                setSelectedSubjectId(res.data[0].id);
            }
        } catch (err) {
            console.error("Subjects error:", err);
        }
    };

    const fetchStudents = async () => {
        if (!selectedClassId) return;
        try {
            const res = await api.get("/teacher/students", {
                params: { class_id: selectedClassId },
                headers: { Authorization: `Bearer ${token}` }
            });
            setStudents(res.data);
        } catch (err) {
            console.error("Students error:", err);
        }
    };

    const fetchAttendance = async () => {
        if (!selectedClassId || !selectedSubjectId || !attendanceDate) return;
        setAttendanceLoading(true);
        try {
            const res = await api.get("/teacher/attendance", {
                params: {
                    class_id: selectedClassId,
                    subject_id: selectedSubjectId,
                    date: attendanceDate
                },
                headers: { Authorization: `Bearer ${token}` }
            });
            setAttendance(res.data);
        } catch (err) {
            console.error("Attendance error:", err);
        } finally {
            setAttendanceLoading(false);
        }
    };

    const handleMarkAttendance = async (e) => {
        e.preventDefault();
        if (!selectedClassId || !selectedSubjectId || !attendanceDate) {
            setError("Please select class, subject, and date");
            return;
        }
        const records = students.map(s => ({
            student_id: s.id,
            status: s.attendanceStatus || "Present"
        }));
        try {
            await api.post("/teacher/attendance", {
                class_id: selectedClassId,
                subject_id: selectedSubjectId,
                date: attendanceDate,
                records
            }, { headers: { Authorization: `Bearer ${token}` } });
            setMessage("Attendance saved successfully");
            fetchAttendance();
        } catch (err) {
            console.error("Mark attendance error:", err);
            setError(err.response?.data?.message || "Failed to save attendance");
        }
    };

    const handleUpdateProfile = async (e) => {
        e.preventDefault();
        try {
            await api.put("/teacher/profile", editForm, {
                headers: { Authorization: `Bearer ${token}` }
            });
            setMessage("Profile updated successfully");
            setShowEditProfile(false);
            fetchDashboard();
        } catch (err) {
            console.error("Update profile error:", err);
            setError(err.response?.data?.message || "Failed to update profile");
        }
    };

    const handleLogout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("role");
        onLogout();
    };

    const clearMessages = () => { setMessage(""); setError(""); };

    if (loading) {
        return <div className="teacher-loading">Loading dashboard...</div>;
    }

    if (!teacher) {
        return <div className="teacher-loading">Unable to load profile.</div>;
    }

    if (showChangePassword) {
        return (
            <ChangePassword
                onBack={() => setShowChangePassword(false)}
            />
        );
    }

    return (
        <div className="teacher-layout">
            <aside className="teacher-sidebar">
                <div className="teacher-logo">
                    <div className="logo-icon">👨‍🏫</div>
                    <div>
                        <h2>StudentHub</h2>
                        <span>Teacher Portal</span>
                    </div>
                </div>

                <nav className="teacher-nav">
                    <div className="nav-title">MAIN MENU</div>
                    <button className={`nav-item ${activeTab === "dashboard" ? "active" : ""}`} onClick={() => setActiveTab("dashboard")}>
                        <span>📊</span> Dashboard
                    </button>
                    <button className={`nav-item ${activeTab === "profile" ? "active" : ""}`} onClick={() => setActiveTab("profile")}>
                        <span>👤</span> My Profile
                    </button>
                    <button className={`nav-item ${activeTab === "classes" ? "active" : ""}`} onClick={() => { setActiveTab("classes"); fetchClasses(); }}>
                        <span>🏫</span> My Classes & Subjects
                    </button>
                    <button className={`nav-item ${activeTab === "students" ? "active" : ""}`} onClick={() => { setActiveTab("students"); fetchStudents(); }}>
                        <span>👥</span> My Students
                    </button>
                    <button className={`nav-item ${activeTab === "attendance" ? "active" : ""}`} onClick={() => { setActiveTab("attendance"); }}>
                        <span>📋</span> Attendance
                    </button>

                    <div className="nav-title second">ACCOUNT</div>
                    <button className="nav-item" onClick={() => setShowChangePassword(true)}>
                        <span>🔐</span> Change Password
                    </button>

                    <div className="sidebar-bottom">
                        <button className="logout-sidebar" onClick={handleLogout}>
                            <span>↪</span> Logout
                        </button>
                    </div>
                </nav>
            </aside>

            <main className="teacher-main">
                <header className="teacher-topbar">
                    <div>
                        <h1>{activeTab === "dashboard" ? "Dashboard" : activeTab === "profile" ? "My Profile" : activeTab === "classes" ? "My Classes & Subjects" : activeTab === "students" ? "My Students" : "Attendance"}</h1>
                        <p>Welcome back, {teacher.name} 👋</p>
                    </div>
                    <div className="teacher-user">
                        <div className="user-avatar">{teacher.name?.charAt(0).toUpperCase()}</div>
                        <span>{teacher.name}</span>
                    </div>
                </header>

                <section className="teacher-content">
                    {message && <div className="teacher-success">✓ {message}</div>}
                    {error && <div className="teacher-error">{error}</div>}

                    {activeTab === "dashboard" && (
                        <div className="dashboard-tabs">
                            <div className="stats-grid">
                                <div className="stat-card blue">
                                    <div className="stat-icon">🏫</div>
                                    <div className="stat-info">
                                        <span>Assigned Classes</span>
                                        <strong>{stats.assigned_classes}</strong>
                                        <small>Classes you teach</small>
                                    </div>
                                </div>
                                <div className="stat-card green">
                                    <div className="stat-icon">📚</div>
                                    <div className="stat-info">
                                        <span>Assigned Subjects</span>
                                        <strong>{stats.assigned_subjects}</strong>
                                        <small>Subjects you teach</small>
                                    </div>
                                </div>
                                <div className="stat-card purple">
                                    <div className="stat-icon">👥</div>
                                    <div className="stat-info">
                                        <span>Assigned Students</span>
                                        <strong>{stats.assigned_students}</strong>
                                        <small>Students in your classes</small>
                                    </div>
                                </div>
                                <div className="stat-card orange">
                                    <div className="stat-icon">📋</div>
                                    <div className="stat-info">
                                        <span>Recent Attendance</span>
                                        <strong>{stats.recent_attendance}</strong>
                                        <small>Last 7 days</small>
                                    </div>
                                </div>
                            </div>

                            <div className="teacher-info-card">
                                <h2>Teacher Information</h2>
                                <div className="info-grid">
                                    <div className="info-item"><span>Teacher ID</span><strong>#{teacher.id}</strong></div>
                                    <div className="info-item"><span>Name</span><strong>{teacher.name}</strong></div>
                                    <div className="info-item"><span>Email</span><strong>{teacher.email}</strong></div>
                                    <div className="info-item"><span>Subject</span><strong>{teacher.subject || "Not specified"}</strong></div>
                                    <div className="info-item"><span>Phone</span><strong>{teacher.phone || "Not provided"}</strong></div>
                                    <div className="info-item"><span>Address</span><strong>{teacher.address || "Not provided"}</strong></div>
                                    <div className="info-item"><span>Hire Date</span><strong>{teacher.hireDate || "Not specified"}</strong></div>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === "profile" && (
                        <div className="profile-section">
                            <div className="profile-card">
                                <div className="profile-header">
                                    <div className="big-avatar">{teacher.name?.charAt(0).toUpperCase()}</div>
                                    <div>
                                        <h2>{teacher.name}</h2>
                                        <p>Teacher ID: #{teacher.id}</p>
                                    </div>
                                </div>
                                <div className="profile-details">
                                    <div className="profile-item"><span>Email</span><strong>{teacher.email}</strong></div>
                                    <div className="profile-item"><span>Subject</span><strong>{teacher.subject || "Not specified"}</strong></div>
                                    <div className="profile-item"><span>Phone</span><strong>{teacher.phone || "Not provided"}</strong></div>
                                    <div className="profile-item"><span>Address</span><strong>{teacher.address || "Not provided"}</strong></div>
                                    <div className="profile-item"><span>Hire Date</span><strong>{teacher.hireDate || "Not specified"}</strong></div>
                                </div>
                                <div className="profile-footer">
                                    <button className="edit-profile-btn" onClick={() => setShowEditProfile(true)}>✏️ Edit Profile</button>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === "classes" && (
                        <div className="classes-section">
                            <div className="classes-card">
                                <div className="classes-header">
                                    <div><h2>My Classes & Subjects</h2><p>Classes and subjects assigned to you</p></div>
                                </div>
                                <div className="table-container">
                                    {classes.length === 0 ? (
                                        <div className="empty-state"><div>🏫</div><h3>No classes assigned</h3><p>Contact admin to assign classes</p></div>
                                    ) : (
                                        <table>
                                            <thead>
                                                <tr><th>Class</th><th>Section</th><th>Academic Year</th><th>Students</th><th>Subjects</th></tr>
                                            </thead>
                                            <tbody>
                                                {classes.map(c => (
                                                    <tr key={c.id} className={selectedClassId === c.id ? "selected" : ""} onClick={() => setSelectedClassId(c.id)}>
                                                        <td><strong>{c.class_name}</strong></td>
                                                        <td>{c.section}</td>
                                                        <td>{c.academic_year}</td>
                                                        <td>{c.student_count}</td>
                                                        <td>
                                                            {subjects.filter(s => true).map(s => <span key={s.id} className="subject-badge">{s.subject_name}</span>)}
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    )}
                                </div>
                                {selectedClassId && subjects.length > 0 && (
                                    <div className="subjects-list">
                                        <h3>Subjects in {classes.find(c => c.id === selectedClassId)?.class_name} {classes.find(c => c.id === selectedClassId)?.section}</h3>
                                        <div className="subject-chips">
                                            {subjects.map(s => (
                                                <span key={s.id} className={`subject-chip ${selectedSubjectId === s.id ? "selected" : ""}`} onClick={() => setSelectedSubjectId(s.id)}>
                                                    {s.subject_name} ({s.subject_code})
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {activeTab === "students" && (
                        <div className="students-section">
                            <div className="students-card">
                                <div className="students-header">
                                    <div><h2>My Students</h2><p>Students in selected class</p></div>
                                    <select className="class-select" value={selectedClassId} onChange={e => { setSelectedClassId(e.target.value); fetchStudents(); }}>
                                        {classes.map(c => <option key={c.id} value={c.id}>{c.class_name} {c.section}</option>)}
                                    </select>
                                </div>
                                <div className="table-container">
                                    {students.length === 0 ? (
                                        <div className="empty-state"><div>👥</div><h3>No students in this class</h3><p>Students will appear when assigned to this class</p></div>
                                    ) : (
                                        <table>
                                            <thead>
                                                <tr><th>ID</th><th>Student</th><th>Email</th><th>Age</th><th>Gender</th><th>City</th></tr>
                                            </thead>
                                            <tbody>
                                                {students.map(s => (
                                                    <tr key={s.id}>
                                                        <td><span className="student-id">#{s.id}</span></td>
                                                        <td><div className="student-name"><div className="student-avatar">{s.name?.charAt(0).toUpperCase()}</div><strong>{s.name}</strong></div></td>
                                                        <td><span className="email-text">{s.email}</span></td>
                                                        <td>{s.age}</td>
                                                        <td><span className={`gender-badge ${s.gender?.toLowerCase() === "male" ? "male" : "female"}`}>{s.gender}</span></td>
                                                        <td>{s.city}</td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === "attendance" && (
                        <div className="attendance-section">
                            <div className="attendance-card">
                                <div className="attendance-header">
                                    <h2>Mark Attendance</h2>
                                </div>
                                <form onSubmit={handleMarkAttendance} className="attendance-form">
                                    <div className="form-row">
                                        <div className="form-group">
                                            <label>Class</label>
                                            <select value={selectedClassId} onChange={e => { setSelectedClassId(e.target.value); setSelectedSubjectId(""); }} required>
                                                <option value="">Select class</option>
                                                {classes.map(c => <option key={c.id} value={c.id}>{c.class_name} {c.section}</option>)}
                                            </select>
                                        </div>
                                        <div className="form-group">
                                            <label>Subject</label>
                                            <select value={selectedSubjectId} onChange={e => setSelectedSubjectId(e.target.value)} required disabled={!selectedClassId}>
                                                <option value="">Select subject</option>
                                                {subjects.map(s => <option key={s.id} value={s.id}>{s.subject_name}</option>)}
                                            </select>
                                        </div>
                                        <div className="form-group">
                                            <label>Date</label>
                                            <input type="date" value={attendanceDate} onChange={e => setAttendanceDate(e.target.value)} required max={new Date().toISOString().split("T")[0]} />
                                        </div>
                                    </div>
                                    <div className="form-actions">
                                        <button type="button" className="load-btn" onClick={fetchAttendance} disabled={!selectedClassId || !selectedSubjectId || !attendanceDate || attendanceLoading}>
                                            {attendanceLoading ? "Loading..." : "Load Attendance"}
                                        </button>
                                        <button type="submit" className="save-btn" disabled={students.length === 0 || !selectedClassId || !selectedSubjectId || !attendanceDate}>
                                            Save Attendance
                                        </button>
                                    </div>
                                </form>

                                {attendance.length > 0 && (
                                    <div className="attendance-records">
                                        <h3>Attendance for {attendanceDate}</h3>
                                        <div className="table-container">
                                            <table>
                                                <thead>
                                                    <tr><th>Student</th><th>Status</th></tr>
                                                </thead>
                                                <tbody>
                                                    {attendance.map(a => (
                                                        <tr key={a.id}>
                                                            <td><div className="student-name"><div className="student-avatar">{a.student_name?.charAt(0).toUpperCase()}</div><strong>{a.student_name}</strong></div></td>
                                                            <td><span className={`status-badge ${a.status.toLowerCase()}`}>{a.status}</span></td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                )}

                                {students.length > 0 && attendance.length === 0 && selectedClassId && selectedSubjectId && attendanceDate && (
                                    <div className="attendance-marking">
                                        <h3>Mark Attendance for {attendanceDate}</h3>
                                        <p className="mark-hint">Click Load Attendance first to see existing records, then modify and save</p>
                                        <div className="students-grid">
                                            {students.map(s => (
                                                <div key={s.id} className="attendance-student">
                                                    <div className="student-info">
                                                        <div className="student-avatar-small">{s.name?.charAt(0).toUpperCase()}</div>
                                                        <span>{s.name}</span>
                                                    </div>
                                                    <div className="attendance-toggle">
                                                        <label className="toggle-label">
                                                            <input type="radio" name={`attendance-${s.id}`} value="Present" defaultChecked onChange={e => { setStudents(students.map(st => st.id === s.id ? {...st, attendanceStatus: "Present"} : st)); }} /> Present
                                                        </label>
                                                        <label className="toggle-label">
                                                            <input type="radio" name={`attendance-${s.id}`} value="Absent" onChange={e => { setStudents(students.map(st => st.id === s.id ? {...st, attendanceStatus: "Absent"} : st)); }} /> Absent
                                                        </label>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {showEditProfile && (
                        <div className="teacher-modal-overlay">
                            <div className="teacher-modal">
                                <button className="modal-close" onClick={() => setShowEditProfile(false)}>×</button>
                                <div className="modal-icon">✏️</div>
                                <h2>Edit Profile</h2>
                                <p>Update your profile information</p>
                                <form onSubmit={handleUpdateProfile}>
                                    <div className="form-grid">
                                        <div className="teacher-form-group">
                                            <label>Email</label>
                                            <input type="email" name="email" value={editForm.email} onChange={e => setEditForm({...editForm, email: e.target.value})} required />
                                        </div>
                                        <div className="teacher-form-group">
                                            <label>Phone</label>
                                            <input type="text" name="phone" value={editForm.phone} onChange={e => setEditForm({...editForm, phone: e.target.value})} />
                                        </div>
                                        <div className="teacher-form-group">
                                            <label>Subject</label>
                                            <input type="text" name="subject" value={editForm.subject} onChange={e => setEditForm({...editForm, subject: e.target.value})} />
                                        </div>
                                        <div className="teacher-form-group">
                                            <label>Address</label>
                                            <input type="text" name="address" value={editForm.address} onChange={e => setEditForm({...editForm, address: e.target.value})} />
                                        </div>
                                    </div>
                                    <div className="modal-actions">
                                        <button type="button" className="cancel-button" onClick={() => setShowEditProfile(false)}>Cancel</button>
                                        <button type="submit" className="save-button">Save Changes</button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    )}
                </section>
            </main>
        </div>
    );
}

export default TeacherDashboard;