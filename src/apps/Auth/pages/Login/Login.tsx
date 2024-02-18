import "./Login.sass";
import {
    Button,
    FormItem,
    Horizontal,
    Input,
    Logo,
    MaterialIcon,
    Title,
    Vertical,
} from "@vertex-center/components";
import Spacer from "../../../../components/Spacer/Spacer";
import { APIError } from "../../../../components/Error/APIError";
import { useState } from "react";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useLogin } from "../../hooks/useLogin";
import { Link, useNavigate } from "react-router-dom";

export default function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const navigate = useNavigate();

    const { login, isLoggingIn, errorLogin } = useLogin({
        onSuccess: () => navigate("/"),
    });

    const onRegister = () => login({ username, password });
    const onUsernameChange = (e: any) => setUsername(e.target.value);
    const onPasswordChange = (e: any) => setPassword(e.target.value);

    return (
        <div className="login">
            <div className="login-container">
                <ProgressOverlay show={isLoggingIn} />
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Login</Title>
                </Horizontal>
                <Vertical gap={20}>
                    <FormItem label="Username" required>
                        <Input onChange={onUsernameChange} />
                    </FormItem>
                    <FormItem label="Password" required>
                        <Input onChange={onPasswordChange} type="password" />
                    </FormItem>
                    <APIError error={errorLogin} />
                    <Link to="/register">I don't have an account</Link>
                    <Horizontal>
                        <Spacer />
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="login" />}
                            onClick={onRegister}
                        >
                            Login
                        </Button>
                    </Horizontal>
                </Vertical>
            </div>
        </div>
    );
}
