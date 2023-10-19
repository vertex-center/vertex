import { HTMLProps } from "react";
import "./Button.sass";
import classNames from "classnames";

export type ButtonType = "colored" | "outlined" | "danger";

export type ButtonProps = HTMLProps<HTMLButtonElement> & {
    type?: ButtonType;
};

export function Button(props: Readonly<ButtonProps>) {
    const { type = "outlined", disabled, className, ...others } = props;

    return (
        <button
            disabled={disabled}
            className={classNames(className, "button", `button-${type}`, {
                "button-disabled": disabled,
            })}
            {...others}
        />
    );
}
