import cx from "classnames";
import { HTMLProps } from "react";
import "./Table.sass";

export type TableProps = HTMLProps<HTMLTableElement>;

export function Table(props: Readonly<TableProps>) {
    const { className, ...others } = props;
    return <table className={cx("table", className)} {...others} />;
}
