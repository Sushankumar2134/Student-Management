import { useState } from "react";
import api from "../services/api";
import "./Login.css";

function ForgotPassword({ onBack }) {

    const [username, setUsername] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleResetPassword = async (e) => {

        e.preventDefault();

        setMessage("");
        setError("");

        if (newPassword !== confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        if (newPassword.length < 6) {
            setError("Password must contain at least 6 characters");
            return;
        }

        setLoading(true);

        try {

            const response = await api.post(
                "/forgot-password",
                {
                    username: username.trim(),
                    new_password: newPassword,
                    confirm_password: confirmPassword
                }
            );

            setMessage(
                response.data.message ||
                "Password reset successfully"
            );

            setUsername("");
            setNewPassword("");
            setConfirmPassword("");

        } catch (error) {

            console.error("Forgot password error:", error);

            if (error.response) {

                setError(
                    error.response.data.message ||
                    "Password reset failed"
                );

            } else {

                setError(
                    "Unable to connect to server"
                );
            }

        } finally {

            setLoading(false);
        }
    };

    return (
        <div
            className="forgot-page"
            onClick={onBack}
        >

            <div
                className="forgot-card"
                onClick={(e) => e.stopPropagation()}
            >

                {/* Close */}

                <button
                    type="button"
                    className="forgot-close"
                    onClick={onBack}
                >
                    ×
                </button>


                {/* Icon */}

                <div className="forgot-icon-large">
                    🔐
                </div>


                {/* Header */}

                <div className="forgot-header">

                    <h2>Reset your password</h2>

                    <p>
                        Enter your username and create a new
                        password for your account.
                    </p>

                </div>


                {/* Form */}

                <form onSubmit={handleResetPassword}>

                    {/* Username */}

                    <div className="forgot-form-group">

                        <label>
                            Username
                        </label>

                        <div className="forgot-input-wrapper">

                            <span>
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


                    {/* New password */}

                    <div className="forgot-form-group">

                        <label>
                            New Password
                        </label>

                        <div className="forgot-input-wrapper">

                            <span>
                                🔒
                            </span>

                            <input
                                type={
                                    showNewPassword
                                        ? "text"
                                        : "password"
                                }
                                placeholder="Enter new password"
                                value={newPassword}
                                onChange={(e) =>
                                    setNewPassword(e.target.value)
                                }
                                required
                            />

                            <button
                                type="button"
                                className="forgot-eye"
                                onClick={() =>
                                    setShowNewPassword(
                                        !showNewPassword
                                    )
                                }
                            >
                                {showNewPassword
                                    ? "🙈"
                                    : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* Confirm password */}

                    <div className="forgot-form-group">

                        <label>
                            Confirm Password
                        </label>

                        <div className="forgot-input-wrapper">

                            <span>
                                🔒
                            </span>

                            <input
                                type={
                                    showConfirmPassword
                                        ? "text"
                                        : "password"
                                }
                                placeholder="Confirm new password"
                                value={confirmPassword}
                                onChange={(e) =>
                                    setConfirmPassword(
                                        e.target.value
                                    )
                                }
                                required
                            />

                            <button
                                type="button"
                                className="forgot-eye"
                                onClick={() =>
                                    setShowConfirmPassword(
                                        !showConfirmPassword
                                    )
                                }
                            >
                                {showConfirmPassword
                                    ? "🙈"
                                    : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* Password match */}

                    {confirmPassword && (

                        <div
                            className={
                                newPassword === confirmPassword
                                    ? "password-match success"
                                    : "password-match danger"
                            }
                        >
                            {newPassword === confirmPassword
                                ? "✓ Passwords match"
                                : "✕ Passwords do not match"}
                        </div>

                    )}


                    {/* Reset */}

                    <button
                        type="submit"
                        className="forgot-reset-button"
                        disabled={loading}
                    >

                        {loading
                            ? "Resetting..."
                            : "Reset Password"}

                    </button>

                </form>


                {/* Messages */}

                {message && (

                    <div className="forgot-success">
                        ✓ {message}
                    </div>

                )}

                {error && (

                    <div className="forgot-error">
                        {error}
                    </div>

                )}


                {/* Back */}

                <button
                    type="button"
                    className="forgot-back"
                    onClick={onBack}
                >
                    ← Back to Login
                </button>

            </div>

        </div>
    );
}

export default ForgotPassword;