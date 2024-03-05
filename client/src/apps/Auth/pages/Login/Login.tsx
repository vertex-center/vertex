import "./Login.sass";
import {
    Button,
    FormItem,
    Horizontal,
    Input,
    Logo,
    Title,
    Vertical,
} from "@vertex-center/components";
import Spacer from "../../../../components/Spacer/Spacer";
import { APIError } from "../../../../components/Error/APIError";
import { useState } from "react";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useLogin } from "../../hooks/useLogin";
import { Link, useNavigate } from "react-router-dom";
import { SignIn } from "@phosphor-icons/react";
import * as yup from "yup";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

const schema = yup
    .object({
        username: yup.string(),
        password: yup.string(),
    })
    .required();

export default function Login() {
    const navigate = useNavigate();

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(schema),
    });

    const { login, isLoggingIn, errorLogin } = useLogin({
        onSuccess: () => navigate("/"),
    });

    const onSubmit = handleSubmit((data) =>
        login({
            username: data.username,
            password: data.password,
        })
    );

    return (
        <div className="login">
            <form className="login-container" onSubmit={onSubmit}>
                <ProgressOverlay show={isLoggingIn} />
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Login</Title>
                </Horizontal>
                <Vertical gap={20}>
                    <FormItem
                        label="Username"
                        error={errors.username?.message?.toString()}
                        required
                    >
                        <Input {...register("username")} />
                    </FormItem>
                    <FormItem
                        label="Password"
                        error={errors.password?.message?.toString()}
                        required
                    >
                        <Input {...register("password")} type="password" />
                    </FormItem>
                    <APIError error={errorLogin} />
                    <Link to="/register">I don't have an account</Link>
                    <Horizontal>
                        <Spacer />
                        <Button
                            type="submit"
                            variant="colored"
                            rightIcon={<SignIn />}
                        >
                            Login
                        </Button>
                    </Horizontal>
                </Vertical>
            </form>
        </div>
    );
}
