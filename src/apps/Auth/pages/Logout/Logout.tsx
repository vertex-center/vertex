import "../Login/Login.sass";
import { APIError } from "../../../../components/Error/APIError";
import { useLogout } from "../../hooks/useLogout";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function Logout() {
    const navigate = useNavigate();

    const { logout, isLoggingOut, errorLogout } = useLogout({
        onSuccess: () => {
            navigate("/login");
        },
    });

    useEffect(() => {
        logout();
    }, []);

    return (
        <div className="login">
            <div className="login-container">
                {isLoggingOut && "Logging out..."}
                {errorLogout && "Failed to logout."}
                <ProgressOverlay show={isLoggingOut} />
                <APIError error={errorLogout} />
            </div>
        </div>
    );
}
