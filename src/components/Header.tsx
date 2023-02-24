import { PropsWithChildren } from "react";

function Logo() {
    return (
        <a
            href="/"
            className="flex flex-row items-center gap-3 px-3 py-2 rounded-md hover:bg-zinc-200"
        >
            <img width={30} src="/images/logo.png" alt="Logo" />
            <h1 className="text-xl font-medium">Vertex</h1>
        </a>
    );
}

function Item({ children }: PropsWithChildren) {
    return (
        <li className="rounded-md px-3 py-1 cursor-pointer hover:bg-zinc-200">
            {children}
        </li>
    );
}

export default function Header() {
    return (
        <div className="flex items-center flex-center block px-3 py-2 border-b text-zinc-900 border-zinc-100">
            <Logo />
            <ul className="flex flex-row ml-4">
                <Item>Marketplace</Item>
                <Item>Installed</Item>
            </ul>
        </div>
    );
}
