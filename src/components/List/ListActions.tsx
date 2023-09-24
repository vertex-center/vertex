import { HTMLProps } from "react";
import classNames from "classnames";

import styles from "./List.module.sass";

type Props = HTMLProps<HTMLDivElement>;

export default function ListActions(props: Props) {
    const { className, ...others } = props;
    return (
        <div className={classNames(styles.actions, className)} {...others} />
    );
}
