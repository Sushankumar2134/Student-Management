import { useState } from "react";
import api from "../services/api";
import "./Login.css";

function Login({ onLogin }) {

    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const [showPassword, setShowPassword] = useState(false);

    const [message, setMessage] = useState("");
    const [loading, setLoading] = useState(false);

    // Forgot password states
    const [showForgot, setShowForgot] = useState(false);
    const [forgotUsername, setForgotUsername] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [forgotMessage, setForgotMessage] = useState("");
    const [forgotError, setForgotError] = useState("");
    const [resetLoading, setResetLoading] = useState(false);

    // =========================
    // LOGIN
    // =========================

    const handleLogin = async (e) => {

        e.preventDefault();

        setMessage("");
        setLoading(true);

        try {

            const response = await api.post("/login", {
                username: username.trim(),
                password: password
            });

            const { token, role } = response.data;

            // Save JWT
            localStorage.setItem("token", token);

            // Save role
            localStorage.setItem("role", role);

            // Send role to App.jsx
            onLogin(role);

        } catch (error) {

            console.error("Login error:", error);

            if (error.response) {

                setMessage(
                    error.response.data.message ||
                    "Invalid username or password"
                );

            } else {

                setMessage(
                    "Unable to connect to server. Make sure Go backend is running."
                );
            }

        } finally {

            setLoading(false);
        }
    };


    // =========================
    // FORGOT PASSWORD
    // =========================

    const handleForgotPassword = async (e) => {

        e.preventDefault();

        setForgotMessage("");
        setForgotError("");

        if (newPassword !== confirmPassword) {

            setForgotError(
                "New password and confirm password do not match"
            );

            return;
        }

        setResetLoading(true);

        try {

            const response = await api.post("/forgot-password", {
                username: forgotUsername.trim(),
                new_password: newPassword,
                confirm_password: confirmPassword
            });

            setForgotMessage(response.data.message);

            setForgotUsername("");
            setNewPassword("");
            setConfirmPassword("");

        } catch (error) {

            console.error("Forgot password error:", error);

            if (error.response) {

                setForgotError(
                    error.response.data.message ||
                    "Password reset failed"
                );

            } else {

                setForgotError(
                    "Unable to connect to server."
                );
            }

        } finally {

            setResetLoading(false);
        }
    };


    // =========================
    // CLOSE FORGOT PASSWORD
    // =========================

    const closeForgotPassword = () => {

        setShowForgot(false);

        setForgotUsername("");
        setNewPassword("");
        setConfirmPassword("");

        setForgotMessage("");
        setForgotError("");
    };


    // =========================
    // GO TO REGISTER
    // =========================

    const goToRegister = () => {
        onLogin("register");
    };

    return (
        <div className="login-page">

            {/* Background decoration */}

            <div className="background-circle circle-one"></div>
            <div className="background-circle circle-two"></div>


            {/* LOGIN CARD */}

            <div className="login-card">

                {/* Logo */}

                <div className="logo">

                    <div className="logo-icon">
                        S
                    </div>

                    <div>
                        <h1>StudentHub</h1>
                        <span>Student Management System</span>
                    </div>

                </div>


                {/* Heading */}

                <div className="login-heading">

                    <h2>Welcome back</h2>

                    <p>
                        Sign in to continue to your account
                    </p>

                </div>


                {/* Login form */}

                <form onSubmit={handleLogin}>

                    {/* Username */}

                    <div className="form-group">

                        <label>Username</label>

                        <div className="input-with-side-icon">

                            <span className="side-input-icon">
                                👤
                            </span>

                            <input
                                type="text"
                                placeholder="Enter your username"
                                value={username}
                                onChange={(e) =>
                                    setUsername(e.target.value)
                                }
                                required
                            />

                        </div>

                    </div>


                    {/* Password */}

                    <div className="form-group">

                        <div className="password-label">

                            <label>Password</label>

                            <button
                                type="button"
                                className="forgot-link"
                                onClick={() => setShowForgot(true)}
                            >
                                Forgot password?
                            </button>

                        </div>

                        <div className="input-wrapper">

                            <span className="input-with-side-icon">
                                🔒
                            </span>

                            <input
                                type={
                                    showPassword
                                        ? "text"
                                        : "password"
                                }
                                placeholder="Enter your password"
                                value={password}
                                onChange={(e) =>
                                    setPassword(e.target.value)
                                }
                                required
                            />

                            <button
                                type="button"
                                className="eye-button"
                                onClick={() =>
                                    setShowPassword(!showPassword)
                                }
                                title={
                                    showPassword
                                        ? "Hide password"
                                        : "Show password"
                                }
                            >
                                {showPassword ? "🙈" : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* Login button */}

                    <button
                        type="submit"
                        className="login-button"
                        disabled={loading}
                    >

                        {loading ? (
                            "Signing in..."
                        ) : (
                            <>
                                Sign In
                                <span className="arrow">
                                    →
                                </span>
                            </>
                        )}

                    </button>

                </form>


                {/* Login error */}

                {message && (

                    <div className="error-message">
                        {message}
                    </div>

                )}


                {/* Register link */}

                <div className="register-link">
                    Don't have an account?
                    <button
                        type="button"
                        className="link-button"
                        onClick={goToRegister}
                    >
                        Register as a Student
                    </button>
                </div>


                {/* Footer */}

                <div className="login-footer">

                    <span className="secure-icon">
                        🔐
                    </span>

                    Secure authentication

                </div>

            </div>


            {/* =====================================
                FORGOT PASSWORD MODAL
            ===================================== */}

            {showForgot && (

                <div
                    className="modal-overlay"
                    onClick={closeForgotPassword}
                >

                    <div
                        className="forgot-modal"
                        onClick={(e) =>
                            e.stopPropagation()
                        }
                    >

                        <button
                            className="close-button"
                            onClick={closeForgotPassword}
                        >
                            ×
                        </button>


                        <div className="forgot-icon">
                            🔐
                        </div>


                        <h2>Reset Password</h2>

                        <p className="forgot-description">
                            Enter your username and create a
                            new password.
                        </p>


                        <form onSubmit={handleForgotPassword}>

                            {/* Username */}

                            <div className="form-group">

                                <label>Username</label>

                                <input
                                    type="text"
                                    placeholder="Enter your username"
                                    value={forgotUsername}
                                    onChange={(e) =>
                                        setForgotUsername(
                                            e.target.value
                                        )
                                    }
                                    required
                                />

                            </div>


                            {/* New password */}

                            <div className="form-group">

                                <label>New Password</label>

                                <input
                                    type="password"
                                    placeholder="Enter new password"
                                    value={newPassword}
                                    onChange={(e) =>
                                        setNewPassword(
                                            e.target.value
                                        )
                                    }
                                    required
                                />

                            </div>


                            {/* Confirm password */}

                            <div className="form-group">

                                <label>
                                    Confirm Password
                                </label>

                                <input
                                    type="password"
                                    placeholder="Confirm new password"
                                    value={confirmPassword}
                                    onChange={(e) =>
                                        setConfirmPassword(
                                            e.target.value
                                        )
                                    }
                                    required
                                />

                            </div>


                            <button
                                type="submit"
                                className="reset-button"
                                disabled={resetLoading}
                            >

                                {resetLoading
                                    ? "Resetting..."
                                    : "Reset Password"}

                            </button>

                        </form>


                        {forgotMessage && (

                            <div className="success-message">
                                {forgotMessage}
                            </div>

                        )}


                        {forgotError && (

                            <div className="error-message">
                                {forgotError}
                            </div>

                        )}


                        <button
                            className="back-login"
                            onClick={closeForgotPassword}
                        >
                            ← Back to Login
                        </button>

                    </div>

                </div>

            )}

        </div>
    );
}

export default Login;