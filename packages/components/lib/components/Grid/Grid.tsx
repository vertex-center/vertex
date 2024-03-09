import cx from "classnames";
import { HTMLProps } from "react";
import "./Grid.sass";

export type GridProps = HTMLProps<HTMLDivElement> & {
    columnSize: number;
};

export function Grid(props: GridProps) {
    const { className, columnSize, ...others } = props;
    return (
        <div
            className={cx("grid", className)}
            style={{
                gridTemplateColumns: `repeat(auto-fill, minmax(${columnSize}px, 1fr))`,
            }}
            {...others}
        />
    );
}
