import React, { ButtonHTMLAttributes } from "react";
import "./Button.sass";
import classNames from "classnames";

export type ButtonType = "colored" | "outlined" | "danger";

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonType;
    leftIcon?: React.JSX.Element;
    rightIcon?: React.JSX.Element;
    borderless?: boolean;
};

export function Button(props: Readonly<ButtonProps>) {
    const {
        variant = "outlined",
        disabled,
        className,
        children,
        leftIcon,
        rightIcon,
        borderless,
        ...others
    } = props;

    return (
        <button
            disabled={disabled}
            className={classNames(className, "button", `button-${variant}`, {
                "button-disabled": disabled,
                "button-borderless": borderless,
            })}
            {...others}
        >
            {leftIcon}
            {children && <span className="button-content">{children}</span>}
            {rightIcon}
        </button>
    );
}
