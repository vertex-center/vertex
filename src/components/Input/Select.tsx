import { HTMLProps } from "react";
import classNames from "classnames";

import styles from "./Input.module.sass";
import { Vertical } from "../Layouts/Layouts";

type OptionProps = HTMLProps<HTMLOptionElement>;

export function Option(props: OptionProps) {
    const { ...others } = props;

    return <option {...others} />;
}

type Props = HTMLProps<HTMLSelectElement> & {};

export default function Select(props: Props) {
    const { className, label, ...others } = props;

    return (
        <Vertical gap={6}>
            <label className={styles.label}>{label}</label>
            <select
                className={classNames(styles.input, className)}
                {...others}
            />
        </Vertical>
    );
}
