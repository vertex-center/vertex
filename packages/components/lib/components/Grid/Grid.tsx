import cx from "classnames";
import { HTMLProps } from "react";
import "./Grid.sass";

export type GridProps = HTMLProps<HTMLDivElement> & {
    rowSize: number;
};

export function Grid(props: GridProps) {
    const { className, rowSize, ...others } = props;
    return (
        <div
            className={cx("grid", className)}
            style={{
                gridTemplateColumns: `repeat(auto-fill, minmax(${rowSize}px, 1fr))`,
            }}
            {...others}
        />
    );
}
