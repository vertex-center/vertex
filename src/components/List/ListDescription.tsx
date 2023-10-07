import { HTMLProps } from "react";
import classNames from "classnames";

import styles from "./List.module.sass";

export type ListDescriptionProps = HTMLProps<HTMLDivElement>;

export default function ListDescription(props: Readonly<ListDescriptionProps>) {
    const { className, ...others } = props;
    return (
        <div
            className={classNames(styles.description, className)}
            {...others}
        />
    );
}
