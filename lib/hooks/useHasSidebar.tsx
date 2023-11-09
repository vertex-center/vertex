import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

const useHasSidebar = (has?: boolean) => {
    const { setHasSidebar } = useContext(PageContext);

    useEffect(() => {
        setHasSidebar?.(has);
        return () => setHasSidebar?.(false);
    }, [has]);
};

export { useHasSidebar };
