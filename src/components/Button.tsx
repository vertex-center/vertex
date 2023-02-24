import { PropsWithChildren } from "react";

type Props = PropsWithChildren<{
    rightSymbol: string;
}>;

export default function Button(props: Props) {
    const { children, rightSymbol } = props;

    return (
        <button
            className="flex gap-2 bg-gray-700 hover:bg-gray-900 text-zinc-50 px-3 py-1.5 rounded-md"
            type="button"
        >
            {children}
            <span className="material-symbols-rounded">{rightSymbol}</span>
        </button>
    );
}
