import { HTMLProps } from "react";
import "./Button.sass";
import classNames from "classnames";

export type ButtonType = "colored" | "outlined";

export type ButtonProps = HTMLProps<HTMLButtonElement> & {
    type: ButtonType;
};

function Button(props: Readonly<ButtonProps>) {
    const { type, disabled, className, ...others } = props;
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

Button.defaultProps = {
    type: "outlined",
};

export default Button;
