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
            {props.leftIcon && (
                <span className="button-icon-left">{props.leftIcon}</span>
            )}
            <span className="button-content">{children}</span>
            {props.rightIcon && (
                <span className="button-icon-right">{props.rightIcon}</span>
            )}
        </button>
    );
}
