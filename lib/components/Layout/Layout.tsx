import "./Layout.sass";
import { HTMLProps } from "react";
import cx from "classnames";
import type * as CSS from "csstype";

export type LayoutProps = HTMLProps<HTMLDivElement> & {
    gap?: number;

    alignItems?: CSS.Property.AlignItems;
    justifyContent?: CSS.Property.JustifyContent;
};

function Layout(props: Readonly<LayoutProps>) {
    const { gap, alignItems, justifyContent, style, ...others } = props;

    return (
        <div
            style={{
                gap,
                alignItems,
                justifyContent,
                ...style,
            }}
            {...others}
        />
    );
}

export function Vertical({ className, ...others }: Readonly<LayoutProps>) {
    return <Layout className={cx("vertical", className)} {...others} />;
}

export function Horizontal({ className, ...others }: Readonly<LayoutProps>) {
    return <Layout className={cx("horizontal", className)} {...others} />;
}
