import classNames from "classnames";

import styles from "./List.module.sass";
import { HTMLProps } from "react";

export type ListIconProps = HTMLProps<HTMLDivElement>;

export default function ListIcon(props: Readonly<ListIconProps>) {
    const { className, ...others } = props;
    return <div className={classNames(styles.icon, className)} {...others} />;
}
