import { forwardRef, InputHTMLAttributes, Ref } from "react";
import "./Input.sass";
import cx from "classnames";

type Props = InputHTMLAttributes<HTMLInputElement>;

export const Input = forwardRef(
    (props: Readonly<Props>, ref: Ref<HTMLInputElement>) => {
        const { className, ...others } = props;
        return (
            <input ref={ref} className={cx("input", className)} {...others} />
        );
    },
);
