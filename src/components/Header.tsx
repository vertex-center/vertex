import { PropsWithChildren } from "react";
import { Link, NavLink } from "react-router-dom";

function Logo() {
    return (
        <Link
            to="/"
            className="flex flex-row items-center gap-3 px-3 py-2 rounded-md hover:bg-zinc-200"
        >
            <img width={30} src="/images/logo.png" alt="Logo" />
            <h1 className="text-xl font-medium">Vertex</h1>
        </Link>
    );
}

type ItemProps = PropsWithChildren<{
    to: string;
}>;

function Item({ children, to }: ItemProps) {
    return (
        <NavLink to={to}>
            <li className="rounded-md px-3 py-1 cursor-pointer hover:bg-zinc-200">
                {children}
            </li>
        </NavLink>
    );
}

export default function Header() {
    return (
        <div className="flex items-center flex-center block px-3 py-2 border-b text-zinc-900 border-zinc-100">
            <Logo />
            <ul className="flex flex-row ml-4">
                <Item to="marketplace">Marketplace</Item>
                <Item to="installed">Installed</Item>
            </ul>
        </div>
    );
}
