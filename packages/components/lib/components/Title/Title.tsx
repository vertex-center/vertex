import { HTMLProps } from "react";
import cx from "classnames";
import "./Title.sass";

export type TitleType = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

type TitleProps = HTMLProps<HTMLHeadingElement> & {
    variant?: TitleType;
};

export function Title(props: Readonly<TitleProps>) {
    const { className, variant, ...others } = props;
    const Component = variant ?? "h2";
    return <Component className={cx(variant, className)} {...others} />;
}
