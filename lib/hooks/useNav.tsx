import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

export const useNav = (path: string) => {
    const { setNavigation } = useContext(PageContext);

    useEffect(() => {
        setNavigation?.(path);
        return () => setNavigation?.(undefined);
    }, [path]);
};
