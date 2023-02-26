import styles from "./Layouts.module.sass";
import { HTMLProps } from "react";
import classNames from "classnames";
import type * as CSS from "csstype";

type Props = HTMLProps<HTMLDivElement> & {
    gap?: number;

    alignItems?: CSS.Property.AlignItems;
    justifyContent?: CSS.Property.JustifyContent;
};

function Layout(props: Props) {
    const { className, gap, alignItems, justifyContent, style, ...others } =
        props;

    return (
        <div
            style={{
                gap,
                alignItems,
                justifyContent,
                ...style,
            }}
            className={classNames({
                [className]: true,
            })}
            {...others}
        />
    );
}

export function Vertical({ className, ...others }: Props) {
    return (
        <Layout
            className={classNames(styles.vertical, className)}
            {...others}
        />
    );
}

export function Horizontal({ className, ...others }: Props) {
    return (
        <Layout
            className={classNames(styles.horizontal, className)}
            {...others}
        />
    );
}
