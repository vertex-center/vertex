import { ElementType } from "react";

export type NavLinkProps<T> = T & {
    as?: ElementType;
};

export function NavLink<T>(props: Readonly<NavLinkProps<T>>) {
    const { as, ...others } = props;
    const Component = as ?? "a";
    return <Component {...others} />;
}
