import "../Login/Login.sass";
import { useState } from "react";
import { useRegister } from "../../hooks/useRegister";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import {
    Button,
    FormItem,
    Horizontal,
    Input,
    Logo,
    Title,
    Vertical,
} from "@vertex-center/components";
import { APIError } from "../../../../components/Error/APIError";
import Spacer from "../../../../components/Spacer/Spacer";
import { Link, useNavigate } from "react-router-dom";
import { SignIn } from "@phosphor-icons/react";
import * as yup from "yup";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

const schema = yup
    .object({
        username: yup.string().required(),
        password: yup.string().min(8).required(),
        confirmPassword: yup
            .string()
            .oneOf([yup.ref("password"), null], "Passwords must match")
            .required(),
    })
    .required();

export default function Register() {
    const navigate = useNavigate();

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(schema),
    });

    const {
        register: _register,
        isRegistering,
        errorRegister,
    } = useRegister({
        onSuccess: () => navigate("/"),
    });

    const onSubmit = handleSubmit((data) =>
        _register({
            username: data.username,
            password: data.password,
        })
    );

    return (
        <div className="login">
            <form className="login-container" onSubmit={onSubmit}>
                <ProgressOverlay show={isRegistering} />
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Register</Title>
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
                    <FormItem
                        label="Confirm Password"
                        error={errors.confirmPassword?.message?.toString()}
                        required
                    >
                        <Input
                            {...register("confirmPassword")}
                            type="password"
                        />
                    </FormItem>
                    <APIError error={errorRegister} />
                    <Link to="/login">I already have an account</Link>
                    <Horizontal>
                        <Spacer />
                        <Button
                            type="submit"
                            variant="colored"
                            rightIcon={<SignIn />}
                        >
                            Register
                        </Button>
                    </Horizontal>
                </Vertical>
            </form>
        </div>
    );
}
