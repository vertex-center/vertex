import cx from "classnames";
import { HTMLProps } from "react";
import "./Table.sass";

export type TableProps = HTMLProps<HTMLTableElement>;

export function Table(props: Readonly<TableProps>) {
    const { className, ...others } = props;
    return <table className={cx("table", className)} {...others} />;
}

export type TableHeadProps = HTMLProps<HTMLTableSectionElement>;

export function TableHead(props: Readonly<TableHeadProps>) {
    const { className, ...others } = props;
    return <thead className={cx("table-head", className)} {...others} />;
}

export type TableHeadCellProps = HTMLProps<HTMLTableHeaderCellElement>;

export function TableHeadCell(props: Readonly<TableHeadCellProps>) {
    const { className, ...others } = props;
    return <th className={cx("table-head-cell", className)} {...others} />;
}

export type TableBodyProps = HTMLProps<HTMLTableSectionElement>;

export function TableBody(props: Readonly<TableBodyProps>) {
    const { className, ...others } = props;
    return <tbody className={cx("table-body", className)} {...others} />;
}

export type TableRowProps = HTMLProps<HTMLTableRowElement>;

export function TableRow(props: Readonly<TableRowProps>) {
    const { className, ...others } = props;
    return <tr className={cx("table-row", className)} {...others} />;
}

export type TableCellProps = HTMLProps<HTMLTableCellElement> & {
    right?: boolean;
};

export function TableCell(props: Readonly<TableCellProps>) {
    const { className, right, ...others } = props;
    return (
        <td
            className={cx("table-cell", right && "table-cell-right", className)}
            {...others}
        />
    );
}
