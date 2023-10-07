import { HTMLProps } from "react";
import classNames from "classnames";

import styles from "./List.module.sass";

type Props = HTMLProps<HTMLDivElement>;

export default function ListActions(props: Readonly<Props>) {
    const { className, ...others } = props;
    return (
        <div className={classNames(styles.actions, className)} {...others} />
    );
}
