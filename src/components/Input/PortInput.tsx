import Input, { InputProps } from "./Input";
import classNames from "classnames";

import styles from "./Input.module.sass";

type Props = InputProps;

export default function PortInput(props: Readonly<Props>) {
    const { className, ...others } = props;

    return (
        <Input
            className={classNames(styles.inputPort, className)}
            {...others}
        />
    );
}
