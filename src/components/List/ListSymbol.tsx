import classNames from "classnames";

import styles from "./List.module.sass";
import { HTMLProps } from "react";

export type ListSymbolProps = HTMLProps<HTMLDivElement>;

export default function ListSymbol(props: Readonly<ListSymbolProps>) {
    const { className, ...others } = props;
    return <div className={classNames(styles.symbol, className)} {...others} />;
}
