import { InputHTMLAttributes } from "react";
import "./Input.sass";
import cx from "classnames";

type Props = InputHTMLAttributes<HTMLInputElement>;

export default function Input(props: Readonly<Props>) {
    const { className, ...others } = props;
    return <input className={cx("input", className)} {...others} />;
}
