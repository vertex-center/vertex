import styles from "./List.module.sass";
import { HTMLProps } from "react";
import classNames from "classnames";

export type ListProps = HTMLProps<HTMLDivElement>;

export default function List(props: ListProps) {
    const { className, ...others } = props;
    return <div className={classNames(styles.list, className)} {...others} />;
}
