import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

const useNav = (path: string) => {
    const { setNavigation } = useContext(PageContext);

    useEffect(() => {
        setNavigation?.(path);
        return () => setNavigation?.(undefined);
    }, [path]);
};

export { useNav };
