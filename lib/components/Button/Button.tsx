import React, { ButtonHTMLAttributes } from "react";
import "./Button.sass";
import classNames from "classnames";

export type ButtonType = "colored" | "outlined" | "danger";

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonType;
    leftIcon?: React.JSX.Element;
    rightIcon?: React.JSX.Element;
};

export function Button(props: Readonly<ButtonProps>) {
    const {
        variant = "outlined",
        disabled,
        className,
        children,
        leftIcon,
        rightIcon,
        ...others
    } = props;

    return (
        <button
            disabled={disabled}
            className={classNames(className, "button", `button-${variant}`, {
                "button-disabled": disabled,
            })}
            {...others}
        >
            {leftIcon && <span className="button-icon-left">{leftIcon}</span>}
            <span className="button-content">{children}</span>
            {rightIcon && (
                <span className="button-icon-right">{rightIcon}</span>
            )}
        </button>
    );
}
