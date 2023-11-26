import "./Login.sass";
import {
    Button,
    Horizontal,
    Logo,
    MaterialIcon,
    TextField,
    Title,
    Vertical,
} from "@vertex-center/components";
import Spacer from "../../../../components/Spacer/Spacer";

export default function Login() {
    return (
        <div className="login">
            <div className="login-container">
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Login</Title>
                </Horizontal>
                <Vertical gap={20}>
                    <TextField label="Username" />
                    <TextField label="Password" type="password" />
                    <Horizontal>
                        <Spacer />
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="login" />}
                        >
                            Login
                        </Button>
                    </Horizontal>
                </Vertical>
            </div>
        </div>
    );
}
