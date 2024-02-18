import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

const useShowSidebar = (show?: boolean) => {
    const { setShowSidebar } = useContext(PageContext);

    useEffect(() => {
        setShowSidebar?.(show);
        return () => setShowSidebar?.(false);
    }, [show]);
};

export { useShowSidebar };
