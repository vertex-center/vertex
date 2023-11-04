import { useContext, useEffect } from "react";
import { HeaderContext } from "../components/Header/Header";

export const useAppNavigation = (path: string) => {
    const { setNavigation } = useContext(HeaderContext);

    useEffect(() => {
        setNavigation(path);
        return () => setNavigation(undefined);
    }, [path]);
};
