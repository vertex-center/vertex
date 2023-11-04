import { createContext, PropsWithChildren, useMemo, useState } from "react";

export const PageContext = createContext<{
    title?: string;
    setTitle?: (title?: string) => void;
    navigation?: string;
    setNavigation?: (url?: string) => void;
}>({
    title: undefined,
    setTitle: undefined,
    navigation: undefined,
    setNavigation: undefined,
});

export function PageProvider(props: PropsWithChildren) {
    const { children } = props;

    const [title, setTitle] = useState<string>();
    const [navigation, setNavigation] = useState<string>();

    const value = useMemo(
        () => ({ title, setTitle, navigation, setNavigation }),
        [title, setTitle, navigation, setNavigation],
    );

    return (
        <PageContext.Provider value={value}>{children}</PageContext.Provider>
    );
}
