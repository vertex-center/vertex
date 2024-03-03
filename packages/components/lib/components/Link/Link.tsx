import { ElementType } from "react";

export type LinkProps<T> = T & {
    as?: ElementType;
};

export function Link<T>(props: Readonly<LinkProps<T>>) {
    const { as, ...others } = props;
    const Component = as ?? "a";
    return <Component {...others} />;
}
