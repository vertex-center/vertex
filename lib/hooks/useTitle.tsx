import { useContext, useEffect } from "react";
import { PageContext } from "../contexts/PageContext";

export const useTitle = (title: string) => {
    const { setTitle } = useContext(PageContext);

    useEffect(() => {
        setTitle?.(title);
        return () => setTitle?.(undefined);
    }, [title]);
};
