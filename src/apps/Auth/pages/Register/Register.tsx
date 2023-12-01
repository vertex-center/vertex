import "../Login/Login.sass";
import { useState } from "react";
import { useRegister } from "../../hooks/useRegister";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import {
    Button,
    Horizontal,
    Logo,
    MaterialIcon,
    TextField,
    Title,
    Vertical,
} from "@vertex-center/components";
import { APIError } from "../../../../components/Error/APIError";
import Spacer from "../../../../components/Spacer/Spacer";
import { Link, useNavigate } from "react-router-dom";

export default function Register() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const navigate = useNavigate();

    const { register, isRegistering, errorRegister } = useRegister({
        onSuccess: () => navigate("/"),
    });

    const onRegister = () => register({ username, password });
    const onUsernameChange = (e: any) => setUsername(e.target.value);
    const onPasswordChange = (e: any) => setPassword(e.target.value);

    return (
        <div className="login">
            <div className="login-container">
                <ProgressOverlay show={isRegistering} />
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Register</Title>
                </Horizontal>
                <Vertical gap={20}>
                    <TextField
                        id="username"
                        label="Username"
                        onChange={onUsernameChange}
                        required
                    />
                    <TextField
                        id="password"
                        label="Password"
                        onChange={onPasswordChange}
                        type="password"
                        required
                    />
                    <APIError error={errorRegister} />
                    <Link to="/login">I already have an account</Link>
                    <Horizontal>
                        <Spacer />
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="login" />}
                            onClick={onRegister}
                        >
                            Register
                        </Button>
                    </Horizontal>
                </Vertical>
            </div>
        </div>
    );
}
