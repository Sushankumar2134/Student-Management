import { useState } from "react";
import api from "../services/api";
import "./Register.css";

function Register({ onLogin }) {

    const [formData, setFormData] = useState({
        name: "",
        email: "",
        username: "",
        password: "",
        confirmPassword: "",
        age: "",
        gender: "",
        city: ""
    });

    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const [submitted, setSubmitted] = useState(false);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
        if (error) setError("");
        if (message) setMessage("");
    };

    const validateForm = () => {
        if (!formData.name.trim()) {
            setError("Full name is required");
            return false;
        }
        if (!formData.email.trim()) {
            setError("Email is required");
            return false;
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(formData.email)) {
            setError("Please enter a valid email address");
            return false;
        }
        if (!formData.username.trim()) {
            setError("Username is required");
            return false;
        }
        if (formData.username.length < 3) {
            setError("Username must be at least 3 characters");
            return false;
        }
        if (!formData.password) {
            setError("Password is required");
            return false;
        }
        if (formData.password.length < 8) {
            setError("Password must be at least 8 characters");
            return false;
        }
        if (formData.password !== formData.confirmPassword) {
            setError("Passwords do not match");
            return false;
        }
        const age = parseInt(formData.age);
        if (!formData.age || isNaN(age) || age < 16 || age > 100) {
            setError("Age must be between 16 and 100");
            return false;
        }
        if (!formData.gender) {
            setError("Gender is required");
            return false;
        }
        if (!formData.city.trim()) {
            setError("City is required");
            return false;
        }
        return true;
    };

    const handleRegister = async (e) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        setLoading(true);
        setError("");
        setMessage("");

        try {
            const response = await api.post("/register", {
                name: formData.name.trim(),
                email: formData.email.trim(),
                username: formData.username.trim(),
                password: formData.password,
                age: parseInt(formData.age),
                gender: formData.gender,
                city: formData.city.trim()
            });

            setMessage(response.data.message || "Registration successful!");
            setSubmitted(true);

        } catch (err) {
            console.error("Registration error:", err);
            if (err.response) {
                setError(err.response.data.message || "Registration failed");
            } else {
                setError("Unable to connect to server. Make sure Go backend is running.");
            }
        } finally {
            setLoading(false);
        }
    };

    const handleBackToLogin = () => {
        onLogin(null);
    };

    return (
        <div className="register-page">
            <div className="background-circle circle-one"></div>
            <div className="background-circle circle-two"></div>

            <div className="register-card">
                <div className="logo">
                    <div className="logo-icon">S</div>
                    <div>
                        <h1>StudentHub</h1>
                        <span>Student Management System</span>
                    </div>
                </div>

                <div className="register-heading">
                    <h2>Create Student Account</h2>
                    <p>Register to access your student portal</p>
                </div>

                {!submitted ? (
                    <form onSubmit={handleRegister} noValidate>
                        <div className="form-row">
                            <div className="form-group">
                                <label>Full Name <span className="required">*</span></label>
                                <input
                                    type="text"
                                    name="name"
                                    placeholder="Enter your full name"
                                    value={formData.name}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Email <span className="required">*</span></label>
                                <input
                                    type="email"
                                    name="email"
                                    placeholder="Enter your email"
                                    value={formData.email}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Username <span className="required">*</span></label>
                                <input
                                    type="text"
                                    name="username"
                                    placeholder="Choose a username"
                                    value={formData.username}
                                    onChange={handleChange}
                                    required
                                    minLength={3}
                                    maxLength={50}
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Password <span className="required">*</span></label>
                                <div className="input-wrapper">
                                    <input
                                        type={showPassword ? "text" : "password"}
                                        name="password"
                                        placeholder="Create a password (min 8 characters)"
                                        value={formData.password}
                                        onChange={handleChange}
                                        required
                                        minLength={8}
                                    />
                                    <button
                                        type="button"
                                        className="eye-button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        title={showPassword ? "Hide password" : "Show password"}
                                    >
                                        {showPassword ? "🙈" : "👁️"}
                                    </button>
                                </div>
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Confirm Password <span className="required">*</span></label>
                                <div className="input-wrapper">
                                    <input
                                        type={showConfirmPassword ? "text" : "password"}
                                        name="confirmPassword"
                                        placeholder="Confirm your password"
                                        value={formData.confirmPassword}
                                        onChange={handleChange}
                                        required
                                    />
                                    <button
                                        type="button"
                                        className="eye-button"
                                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                        title={showConfirmPassword ? "Hide password" : "Show password"}
                                    >
                                        {showConfirmPassword ? "🙈" : "👁️"}
                                    </button>
                                </div>
                            </div>
                        </div>

                        <div className="form-row two-cols">
                            <div className="form-group">
                                <label>Age <span className="required">*</span></label>
                                <input
                                    type="number"
                                    name="age"
                                    placeholder="Enter your age"
                                    value={formData.age}
                                    onChange={handleChange}
                                    required
                                    min={16}
                                    max={100}
                                />
                            </div>
                            <div className="form-group">
                                <label>Gender <span className="required">*</span></label>
                                <select
                                    name="gender"
                                    value={formData.gender}
                                    onChange={handleChange}
                                    required
                                >
                                    <option value="">Select gender</option>
                                    <option value="Male">Male</option>
                                    <option value="Female">Female</option>
                                </select>
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label>City <span className="required">*</span></label>
                                <input
                                    type="text"
                                    name="city"
                                    placeholder="Enter your city"
                                    value={formData.city}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        {error && (
                            <div className="error-message">
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            className="register-button"
                            disabled={loading}
                        >
                            {loading ? "Registering..." : "Create Account"}
                        </button>
                    </form>
                ) : (
                    <div className="success-state">
                        <div className="success-icon">✓</div>
                        <h3>Registration Successful!</h3>
                        <p>{message}</p>
                        <p className="success-detail">You can now log in with your username and password.</p>
                        <button
                            className="login-link-button"
                            onClick={handleBackToLogin}
                        >
                            ← Back to Login
                        </button>
                    </div>
                )}

                <div className="register-footer">
                    Already have an account? 
                    <button className="link-button" onClick={handleBackToLogin}>
                        Sign In
                    </button>
                </div>
            </div>
        </div>
    );
}

export default Register;