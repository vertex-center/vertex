import { HTMLProps } from "react";

import styles from "./Input.module.sass";
import classNames from "classnames";
import { Vertical } from "../Layouts/Layouts";

type Props = HTMLProps<HTMLInputElement> & {};

export default function Input(props: Props) {
    const { className, label, ...others } = props;

    return (
        <Vertical gap={6}>
            <label className={styles.label}>{label}</label>
            <input
                {...others}
                className={classNames(styles.input, className)}
            />
        </Vertical>
    );
}
