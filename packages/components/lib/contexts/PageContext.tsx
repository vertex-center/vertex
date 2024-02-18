import { createContext, PropsWithChildren, useMemo, useState } from "react";

export const PageContext = createContext<{
    title?: string;
    setTitle?: (title?: string) => void;
    hasSidebar?: boolean;
    setHasSidebar?: (has?: boolean) => void;
    showSidebar?: boolean;
    setShowSidebar?: (show?: boolean) => void;
    toggleShowSidebar?: () => void;
}>({
    title: undefined,
    setTitle: undefined,
    hasSidebar: undefined,
    setHasSidebar: undefined,
    showSidebar: undefined,
    setShowSidebar: undefined,
});

export function PageProvider(props: PropsWithChildren) {
    const { children } = props;

    const [title, setTitle] = useState<string>();
    const [hasSidebar, setHasSidebar] = useState<boolean>();
    const [showSidebar, setShowSidebar] = useState<boolean>();

    const value = useMemo(
        () => ({
            title,
            setTitle,
            hasSidebar,
            setHasSidebar,
            showSidebar,
            setShowSidebar,
            toggleShowSidebar: () => setShowSidebar((show) => !show),
        }),
        [
            title,
            setTitle,
            hasSidebar,
            setHasSidebar,
            showSidebar,
            setShowSidebar,
        ],
    );

    return (
        <PageContext.Provider value={value}>{children}</PageContext.Provider>
    );
}
