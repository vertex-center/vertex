import { forwardRef, InputHTMLAttributes, Ref } from "react";
import "./Input.sass";
import cx from "classnames";

export type InputProps = InputHTMLAttributes<HTMLInputElement>;

export const Input = forwardRef(
    (props: Readonly<InputProps>, ref: Ref<HTMLInputElement>) => {
        const { className, ...others } = props;
        return (
            <input ref={ref} className={cx("input", className)} {...others} />
        );
    },
);
