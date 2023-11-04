import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

const useTitle = (title: string) => {
    const { setTitle } = useContext(PageContext);

    useEffect(() => {
        setTitle?.(title);
        return () => setTitle?.(undefined);
    }, [title]);
};

export { useTitle };
