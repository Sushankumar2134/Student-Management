import { useState } from "react";
import api from "../services/api";
import "./ChangePassword.css";

function ChangePassword({ onBack }) {

    const [currentPassword, setCurrentPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [showCurrent, setShowCurrent] = useState(false);
    const [showNew, setShowNew] = useState(false);
    const [showConfirm, setShowConfirm] = useState(false);

    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {

        e.preventDefault();

        setMessage("");
        setError("");

        if (newPassword !== confirmPassword) {
            setError("New password and confirm password do not match");
            return;
        }

        if (newPassword.length < 6) {
            setError("New password must contain at least 6 characters");
            return;
        }

        setLoading(true);

        try {

            const token = localStorage.getItem("token");

            const response = await api.put(
                "/change-password",
                {
                    current_password: currentPassword,
                    new_password: newPassword,
                    confirm_password: confirmPassword
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                }
            );

            setMessage(
                response.data.message ||
                "Password changed successfully"
            );

            setCurrentPassword("");
            setNewPassword("");
            setConfirmPassword("");

        } catch (error) {

            console.error("Change password error:", error);

            if (error.response) {

                setError(
                    error.response.data.message ||
                    "Unable to change password"
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

        <div className="change-password-page">

            <div className="change-password-card">

                <button
                    className="change-close"
                    onClick={onBack}
                >
                    ×
                </button>


                <div className="change-icon">
                    🔐
                </div>


                <h2>
                    Change Password
                </h2>

                <p className="change-subtitle">
                    Update your account password securely.
                </p>


                <form onSubmit={handleSubmit}>

                    {/* Current Password */}

                    <div className="change-group">

                        <label>
                            Current Password
                        </label>

                        <div className="password-input">

                            <span>🔒</span>

                            <input
                                type={
                                    showCurrent
                                        ? "text"
                                        : "password"
                                }
                                placeholder="Enter current password"
                                value={currentPassword}
                                onChange={(e) =>
                                    setCurrentPassword(e.target.value)
                                }
                                required
                            />

                            <button
                                type="button"
                                onClick={() =>
                                    setShowCurrent(!showCurrent)
                                }
                            >
                                {showCurrent ? "🙈" : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* New Password */}

                    <div className="change-group">

                        <label>
                            New Password
                        </label>

                        <div className="password-input">

                            <span>🔑</span>

                            <input
                                type={
                                    showNew
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
                                onClick={() =>
                                    setShowNew(!showNew)
                                }
                            >
                                {showNew ? "🙈" : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* Confirm Password */}

                    <div className="change-group">

                        <label>
                            Confirm New Password
                        </label>

                        <div className="password-input">

                            <span>🔑</span>

                            <input
                                type={
                                    showConfirm
                                        ? "text"
                                        : "password"
                                }
                                placeholder="Confirm new password"
                                value={confirmPassword}
                                onChange={(e) =>
                                    setConfirmPassword(e.target.value)
                                }
                                required
                            />

                            <button
                                type="button"
                                onClick={() =>
                                    setShowConfirm(!showConfirm)
                                }
                            >
                                {showConfirm ? "🙈" : "👁️"}
                            </button>

                        </div>

                    </div>


                    {/* Password match */}

                    {confirmPassword && (

                        <div
                            className={
                                newPassword === confirmPassword
                                    ? "password-match-ok"
                                    : "password-match-error"
                            }
                        >
                            {newPassword === confirmPassword
                                ? "✓ Passwords match"
                                : "✕ Passwords do not match"}
                        </div>

                    )}


                    <button
                        type="submit"
                        className="change-password-button"
                        disabled={loading}
                    >
                        {loading
                            ? "Changing Password..."
                            : "Change Password"}
                    </button>

                </form>


                {message && (

                    <div className="change-success">
                        ✓ {message}
                    </div>

                )}


                {error && (

                    <div className="change-error">
                        {error}
                    </div>

                )}


                <button
                    className="back-dashboard"
                    onClick={onBack}
                >
                    ← Back to Dashboard
                </button>

            </div>

        </div>
    );
}

export default ChangePassword;